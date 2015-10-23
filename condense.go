package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io"
	"os"
	"path"
	"strings"
)

type Template struct {
	AWSTemplateFormatVersion string                 `json:",omitempty"`
	Description              string                 `json:",omitempty"`
	Metadata                 map[string]interface{} `json:",omitempty"`
	Parameters               map[string]interface{} `json:",omitempty"`
	Mappings                 map[string]interface{} `json:",omitempty"`
	Conditions               map[string]interface{} `json:",omitempty"`
	Resources                map[string]interface{} `json:",omitempty"`
	Outputs                  map[string]interface{} `json:",omitempty"`
}

func getStackOutputs(name string) (map[string]interface{}, bool) {
	svc := cloudformation.New(&aws.Config{Region: aws.String("eu-west-1")})
	description, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: &name,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		return nil, false
	}

	if len(description.Stacks) != 1 {
		panic(fmt.Sprintf(
			"Description of [%s] did not return in exactly one Stack",
			name,
		))
	}

	outputs := map[string]interface{}{}
	stack := description.Stacks[0]
	for _, output := range stack.Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	return outputs, true
}

type ReffableType int

const (
	Parameter ReffableType = iota
	Condition
	Resource
)

func (t *Template) Reffables() map[string]ReffableType {
	reffables := map[string]ReffableType{}

	for k := range t.Parameters {
		reffables[k] = Parameter
	}

	for k := range t.Conditions {
		reffables[k] = Condition
	}

	for k := range t.Resources {
		reffables[k] = Resource
	}

	return reffables
}

func (t *Template) IsReffable(ref string) bool {
	if _, ok := t.Reffables()[ref]; ok {
		return true
	}

	return false
}

func merge(a map[string]interface{}, b map[string]interface{}) {
	for k, spec := range b {
		a[k] = spec
	}
}

func (t *Template) Merge(o *Template) {
	merge(t.Metadata, o.Metadata)
	merge(t.Parameters, o.Parameters)
	merge(t.Mappings, o.Mappings)
	merge(t.Conditions, o.Conditions)
	merge(t.Resources, o.Resources)
	merge(t.Outputs, o.Outputs)
}

type Visitor func(interface{}) interface{}

func walk(p interface{}, v Visitor) interface{} {
	switch p := p.(type) {
	default:
		panic(fmt.Sprintf("unknown type %T\n", p))
	case *Template:
		if len(p.Parameters) > 0 {
			p.Parameters = walk(p.Parameters, v).(map[string]interface{})
		}
		if len(p.Resources) > 0 {
			p.Resources = walk(p.Resources, v).(map[string]interface{})
		}
	case []interface{}:
		for i, value := range p {
			p[i] = walk(value, v)
		}
	case map[string]interface{}:
		for k, value := range p {
			p[k] = walk(value, v)
		}
	case string:
	case bool:
	case int:
	}

	return v(p)
}

func (t *Template) Walk(v Visitor) {
	walk(t, v)
}

func isRef(p map[string]interface{}) bool {
	if len(p) != 1 {
		return false
	}
	if _, ok := p["Ref"]; !ok {
		return false
	}
	if _, ok := p["Ref"].(string); !ok {
		return false
	}
	return true
}

func isGetAtt(p map[string]interface{}) bool {
	if len(p) != 1 {
		return false
	}
	if _, ok := p["Fn::GetAtt"]; !ok {
		return false
	}
	if _, ok := p["Fn::GetAtt"].([]interface{}); !ok {
		return false
	}
	if len(p["Fn::GetAtt"].([]interface{})) != 2 {
		return false
	}
	if _, ok := p["Fn::GetAtt"].([]interface{})[0].(string); !ok {
		return false
	}
	if _, ok := p["Fn::GetAtt"].([]interface{})[1].(string); !ok {
		return false
	}
	return true
}

type RefType int

const (
	Ref RefType = iota
	GetAtt
)

type Refish struct {
	Type RefType
	path []string
}

func IsRefish(p map[string]interface{}) bool {
	return isRef(p) || isGetAtt(p)
}

func NewRefish(p map[string]interface{}) *Refish {
	if isRef(p) {
		return &Refish{
			Type: Ref,
			path: []string{p["Ref"].(string)},
		}
	}

	if isGetAtt(p) {
		var path []string
		args := p["Fn::GetAtt"].([]interface{})
		path = append(path, args[0].(string))
		for _, part := range strings.Split(args[1].(string), ".") {
			path = append(path, part)
		}

		return &Refish{
			Type: GetAtt,
			path: path,
		}
	}

	panic("non-Refish passed to NewRefish")
}

func (r *Refish) Lead(inputs Inputs) string {
	lead := r.path[0]
	aliasKey := fmt.Sprintf("%s.", lead)

	if alias, ok := inputs.Get([]string{aliasKey}); ok {
		return alias.(string)
	}

	return lead
}

func (r *Refish) Path(inputs Inputs) []string {
	var path []string

	path = append(path, r.Lead(inputs))
	for _, part := range r.path[1:] {
		path = append(path, part)
	}

	return path
}

func (r *Refish) Map(inputs Inputs) map[string]interface{} {
	if r.Type == Ref {
		return map[string]interface{}{
			"Ref": r.Lead(inputs),
		}
	}

	return map[string]interface{}{
		"Fn::GetAtt": r.Path(inputs),
	}
}

type Inputs map[string]interface{}

func NewInputs(raw map[string]interface{}) *Inputs {
	inputs := Inputs(raw)
	return &inputs
}

func (inputs *Inputs) Map() map[string]interface{} {
	return map[string]interface{}(*inputs)
}

func (inputs *Inputs) Merge(other *Inputs) {
	merge(map[string]interface{}(*inputs), map[string]interface{}(*other))
}

// Note: intentionally does not handle lookup-by-array-index
func (inputs *Inputs) Get(path []string) (interface{}, bool) {
	var ok bool
	var p map[string]interface{}

	next := interface{}(inputs.Map())
	for _, part := range path {
		p, ok = next.(map[string]interface{})
		if !ok {
			// "next" was not something we could use for mapping
			return nil, false
		}

		next, ok = p[part]
		if !ok {
			// "p" did not contain a [part] key
			return nil, false
		}
	}

	return next, true
}

func getComponent(name string) (*Template, bool) {
	filename := path.Join("components", strings.Join([]string{name, "json"}, "."))
	reader, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, false
	}

	var template Template
	decoder := json.NewDecoder(reader)
	decoder.Decode(&template)
	return &template, true
}

func main() {
	var templateFilename string
	var parametersFilename string

	flag.StringVar(&templateFilename,
		"template", "-",
		"CloudFormation Template to process")

	flag.StringVar(&parametersFilename,
		"parameters", "",
		"CloudFormation Template to process")

	flag.Parse()

	var jsonStream io.Reader
	var inputStream io.Reader
	var err error

	if templateFilename == "-" {
		jsonStream = os.Stdin
	} else if jsonStream, err = os.Open(templateFilename); err != nil {
		panic(err)
	}

	if parametersFilename == "" {
		inputStream = strings.NewReader("{}")
	} else if inputStream, err = os.Open(parametersFilename); err != nil {
		panic(err)
	}

	dec := json.NewDecoder(jsonStream)
	idec := json.NewDecoder(inputStream)

	var t Template
	var inputs Inputs
	if err := idec.Decode(&inputs); err != nil {
		panic(err)
	}

	if err := dec.Decode(&t); err != nil {
		panic(err)
	}

	pending := 1
	for pending > 0 {
		t.Walk(func(p interface{}) interface{} {
			switch p := p.(type) {
			case string:
				if p == "BbbValue" {
					return "ReplacedBbbValue"
				}
			case map[string]interface{}:
				if !IsRefish(p) {
					return p
				}

				r := NewRefish(p)
				if t.IsReffable(r.Lead(inputs)) {
					return p
				}

				value, ok := inputs.Get(r.Path(inputs))
				if ok {
					// existed in inputs: return that
					return value
				}

				// search for a component to Include
				component, ok := getComponent(r.Lead(inputs))
				if ok {
					t.Merge(component)
					pending = pending + 1
					if t.IsReffable(r.Lead(inputs)) {
						// after merging in the Include, everything was okay
						return r.Map(inputs)
					}

					panic(fmt.Sprintf(
						"Including %s component did not result in %s being reffable",
						r.Lead(inputs), r.Lead(inputs),
					))
				}

				outputs, ok := getStackOutputs(r.Lead(inputs))
				if ok {
					outputsMap := map[string]interface{}{}
					outputsMap[r.Lead(inputs)] = map[string]interface{}{
						"Outputs": outputs,
					}
					inputs.Merge(NewInputs(outputsMap))

					if value, ok := inputs.Get(r.Path(inputs)); ok {
						return value
					} else {
						panic(fmt.Sprintf(
							"Unknown value for %v even after Merging in Stack Outputs for %s",
							r.Path(inputs),
							r.Lead(inputs),
						))
					}
				}

				// this is probably a failure, but let the next step in the chain
				// handle the reporting of that failure.
				return p
			}

			return p
		})

		pending = pending - 1
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(t)
}
