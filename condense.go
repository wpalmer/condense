package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"deepcloudformationoutputs"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io"
	"os"
	"path"
	"strings"
	"fallbackmap"
	"deepalias"
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
	sources fallbackmap.Deep
}

func IsRefish(p map[string]interface{}) bool {
	return isRef(p) || isGetAtt(p)
}

func NewRefish(p map[string]interface{}, sources fallbackmap.Deep) *Refish {
	if isRef(p) {
		return &Refish{
			Type: Ref,
			path: strings.Split(p["Ref"].(string), "."),
			sources: sources,
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
			sources: sources,
		}
	}

	panic("non-Refish passed to NewRefish")
}

func (r *Refish) Lead() string {
	return r.Path()[0]
}

func (r *Refish) Path() []string {
	translated, _ := deepalias.DeAlias(r.path, r.sources)
	return translated
}

func (r *Refish) Map(in *fallbackmap.FallbackMap) map[string]interface{} {
	if r.Type == Ref {
		return map[string]interface{}{
			"Ref": strings.Join(r.Path(), "."),
		}
	}

	var outpath []interface{}
	path := r.Path()

	if len(path) == 1 {
		outpath = []interface{}{path[0]}
	} else if len(path) > 1 {
		outpath = []interface{}{path[0], strings.Join(path[1:], ".")}
	}

	return map[string]interface{}{
		"Fn::GetAtt": outpath,
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

type InputsFlag struct {
	inputs *fallbackmap.FallbackMap
	sources []string
}

func (f *InputsFlag) Get() *fallbackmap.FallbackMap {
	if f.inputs == nil {
		f.inputs = &fallbackmap.FallbackMap{}
	}

	return f.inputs
}

func (f InputsFlag) String() string {
	return fmt.Sprintf("%v", f.sources)
}

func (f *InputsFlag) Set(parametersFilename string) (err error) {
	var inputStream io.Reader
	var in map[string]interface{}

	if inputStream, err = os.Open(parametersFilename); err != nil {
		return err
	}

	inputDecoder := json.NewDecoder(inputStream)
	if err := inputDecoder.Decode(&in); err != nil {
		return err
	}

	if f.inputs == nil {
		f.inputs = &fallbackmap.FallbackMap{}
	}

	f.inputs.Override(fallbackmap.DeepMap(in))
	return nil
}

type OutputWhat int
const (
	OutputTemplate = iota
	OutputParameters
)

type OutputWhatFlag struct {
	what OutputWhat
}

func (f OutputWhatFlag) Get() OutputWhat {
	return f.what
}

func (f OutputWhatFlag) String() string {
	switch f.what {
	case OutputTemplate:
		return "template"
	case OutputParameters:
		return "parameters"
	default:
		return "[unknown]"
	}
}

func (f *OutputWhatFlag) Set(input string) error {
	switch input {
	case "template":
		f.what = OutputTemplate
	case "parameters":
		f.what = OutputParameters
	default:
		return fmt.Errorf("Unknown -output `%s' requested", input)
	}

	return nil
}

func main() {
	var templateFilename string
	var outputWhat OutputWhatFlag
	var inputParameters InputsFlag

	flag.StringVar(&templateFilename,
		"template", "-",
		"CloudFormation Template to process")

	flag.Var(&inputParameters,
		"parameters",
		"File to use of input parameters (can be specified multiple times)")

	flag.Var(&outputWhat,
		"output",
		"What to output after processing the Template")

	flag.Parse()

	var jsonStream io.Reader
	var err error

	if templateFilename == "-" {
		jsonStream = os.Stdin
	} else if jsonStream, err = os.Open(templateFilename); err != nil {
		panic(err)
	}

	dec := json.NewDecoder(jsonStream)

	var t Template
	in := inputParameters.Get()
	if err := dec.Decode(&t); err != nil {
		panic(err)
	}

	sources := fallbackmap.FallbackMap{}
	sources.Attach(in)
	sources.Attach(deepalias.DeepAlias{&sources})
	sources.Attach(deepcloudformationoutputs.NewDeepCloudFormationOutputs("eu-west-1"))

	pending := 1
	for pending > 0 {
		t.Walk(func(p interface{}) interface{} {
			switch p := p.(type) {
			case map[string]interface{}:
				if !IsRefish(p) {
					return p
				}

				r := NewRefish(p, fallbackmap.Deep(&sources))
				if t.IsReffable(r.Lead()) {
					return p
				}

				value, ok := sources.Get(r.Path())
				if ok {
					// existed in inputs: return that
					return value
				}

				// search for a component to Include
				component, ok := getComponent(r.Lead())
				if ok {
					t.Merge(component)
					pending = pending + 1
					if t.IsReffable(r.Lead()) {
						// after merging in the Include, everything was okay
						return r.Map(in)
					}

					panic(fmt.Sprintf(
						"Including %s component did not result in %s being reffable",
						r.Lead(), r.Lead(),
					))
				}

				// this is probably a failure, but let the next step in the chain
				// handle the reporting of that failure.
				return p
			}

			return p
		})

		pending = pending - 1
	}

	switch outputWhat.Get() {
	case OutputTemplate:
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(t)
	case OutputParameters:
		parameters := []cloudformation.Parameter{}
		for name, _ := range t.Parameters {
			input_value, ok := in.Get([]string{name})
			if !ok {
				continue
			}

			parameters = append(parameters, func(name string, value interface{}) cloudformation.Parameter {
				var stringval string
				var boolval bool
				var p map[string]interface{}
				var ok bool

				if p, ok = value.(map[string]interface{}); ok {
					if !IsRefish(p) {
						panic(fmt.Sprintf("non-string Parameter: %s", name))
					}

					var looked_up interface{}
					r := NewRefish(p, fallbackmap.Deep(&sources))
					looked_up, ok = sources.Get(r.Path())
					if ok {
						// existed in inputs: return that
						stringval = fmt.Sprintf("%s", looked_up)
					} else {
						panic(fmt.Sprintf("Invalid reference %v for Parameter: %s", r.Path(), name))
					}
				} else {
					stringval = fmt.Sprintf("%s", value)
				}

				boolval = false
				return cloudformation.Parameter{
					ParameterKey: &name,
					ParameterValue: &stringval,
					UsePreviousValue: &boolval,
				}
			}(name, input_value))
		}

		enc := json.NewEncoder(os.Stdout)
		enc.Encode(parameters)
	}
}
