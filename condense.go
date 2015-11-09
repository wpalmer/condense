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
	"inputs"
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
	if strings.HasPrefix(ref, "AWS::") {
		return true
	}

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
	if o.Metadata != nil && t.Metadata == nil {
		t.Metadata = make(map[string]interface{})
	}
	merge(t.Metadata, o.Metadata)

	if o.Parameters != nil && t.Parameters == nil {
		t.Parameters = make(map[string]interface{})
	}
	merge(t.Parameters, o.Parameters)

	if o.Mappings != nil && t.Mappings == nil {
		t.Mappings = make(map[string]interface{})
	}
	merge(t.Mappings, o.Mappings)

	if o.Conditions != nil && t.Conditions == nil {
		t.Conditions = make(map[string]interface{})
	}
	merge(t.Conditions, o.Conditions)

	if o.Resources != nil && t.Resources == nil {
		t.Resources = make(map[string]interface{})
	}
	merge(t.Resources, o.Resources)

	if o.Outputs != nil && t.Outputs == nil {
		t.Outputs = make(map[string]interface{})
	}
	merge(t.Outputs, o.Outputs)
}

type Visitor func(interface{}) interface{}

func isComment(p interface{}) bool {
	var m map[string]interface{}
	var ok bool

	m, ok = p.(map[string]interface{})
	if !ok {
		return false
	}

	if len(m) != 1 {
		return false
	}

	_, ok = m["$comment"]
	return ok
}

func walk(p interface{}, v Visitor) interface{} {
	switch typed := p.(type) {
	default:
		panic(fmt.Sprintf("unknown type: %T\n", typed))
	case *Template:
		if len(typed.Parameters) > 0 {
			typed.Parameters = walk(typed.Parameters, v).(map[string]interface{})
		}

		if len(typed.Conditions) > 0 {
			typed.Conditions = walk(typed.Conditions, v).(map[string]interface{})
		}

		if len(typed.Resources) > 0 {
			typed.Resources = walk(typed.Resources, v).(map[string]interface{})
		}
	case []interface{}:
		filtered := []interface{}{}
		for _, value := range typed {
			if isComment(value) {
				continue
			}

			filtered = append(filtered, walk(value, v))
		}

		p = interface{}(filtered)
	case map[string]interface{}:
		for k, value := range typed {
			if k == "$comment" {
				delete(typed, k)
				continue
			}

			typed[k] = walk(value, v)
		}
	case string:
	case bool:
	case int:
	case float64:
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

func (r *Refish) Lead(in inputs.Inputs) string {
	lead := r.path[0]
	aliasKey := fmt.Sprintf("%s.", lead)

	if alias, ok := in.Get([]string{aliasKey}); ok {
		return alias.(string)
	}

	return lead
}

func (r *Refish) Path(in inputs.Inputs) []string {
	var path []string

	path = append(path, r.Lead(in))
	for _, part := range r.path[1:] {
		path = append(path, part)
	}

	return path
}

func (r *Refish) Map(in inputs.Inputs) map[string]interface{} {
	if r.Type == Ref {
		return map[string]interface{}{
			"Ref": r.Lead(in),
		}
	}

	return map[string]interface{}{
		"Fn::GetAtt": r.Path(in),
	}
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
	var in inputs.Inputs
	if err := idec.Decode(&in); err != nil {
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
				if t.IsReffable(r.Lead(in)) {
					return p
				}

				value, ok := in.Get(r.Path(in))
				if ok {
					// existed in inputs: return that
					return value
				}

				// search for a component to Include
				component, ok := getComponent(r.Lead(in))
				if ok {
					t.Merge(component)
					pending = pending + 1
					if t.IsReffable(r.Lead(in)) {
						// after merging in the Include, everything was okay
						return r.Map(in)
					}

					panic(fmt.Sprintf(
						"Including %s component did not result in %s being reffable",
						r.Lead(in), r.Lead(in),
					))
				}

				outputs, ok := getStackOutputs(r.Lead(in))
				if ok {
					outputsMap := map[string]interface{}{}
					outputsMap[r.Lead(in)] = map[string]interface{}{
						"Outputs": outputs,
					}
					in.Attach(r.Lead(in), inputs.NewInputs(outputsMap))

					if value, ok := in.Get(r.Path(in)); ok {
						return value
					} else {
						panic(fmt.Sprintf(
							"Unknown value for %v even after Merging in Stack Outputs for %s",
							r.Path(in),
							r.Lead(in),
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
