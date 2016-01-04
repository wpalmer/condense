package main

import (
	"deepalias"
	"deepcloudformationoutputs"
	"encoding/json"
	"fallbackmap"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"io"
	"os"
	"path"
	"strings"

	"condense/template"
	"condense/template/rules"
)

type InputsFlag struct {
	inputs  *fallbackmap.FallbackMap
	sources map[string]map[string]interface{}
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

	if f.sources == nil {
		f.sources = make(map[string]map[string]interface{})
	}
	f.sources[parametersFilename] = in
	return nil
}

func (f *InputsFlag) Sources() (map[string]map[string]interface{}) {
	return f.sources
}

type OutputWhat int

const (
	OutputTemplate = iota
	OutputParameters
	OutputCredentials
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
	case OutputCredentials:
		return "credentials"
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
	case "credentials":
			f.what = OutputCredentials
	default:
		return fmt.Errorf("Unknown -output `%s' requested", input)
	}

	return nil
}

func getComponent(name string, templateRules *template.Rules) (*map[string]interface{}, bool) {
	filename := path.Join("components", strings.Join([]string{name, "json"}, "."))
	reader, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, false
	}

	var t map[string]interface{}
	decoder := json.NewDecoder(reader)
	decoder.Decode(&t)

	t = template.Process(interface{}(t), templateRules).(map[string]interface{})
	return &t, true
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
	t := make(map[string]interface{})
	if err := dec.Decode(&t); err != nil {
		panic(err)
	}

	components := make(map[string]*map[string]interface{})

	sources := fallbackmap.FallbackMap{}
	templateRules := template.Rules{}

	sources.Attach(inputParameters.Get())
	sources.Attach(deepalias.DeepAlias{&sources})
	sources.Attach(deepcloudformationoutputs.NewDeepCloudFormationOutputs("eu-west-1"))
	sources.Attach(fallbackmap.DeepFunc(func(path []string) (interface{}, bool) {
		var ok bool
		var component *map[string]interface{}

		if len(path) < 1 {
			return nil, false
		}

		// if all else fails, try a component (always return nothing)
		if _, ok = components[path[0]]; ok {
			return nil, false // already loaded
		}

		if component, ok = getComponent(path[0], &templateRules); ok {
			components[path[0]] = component
		}

		return nil, false
	}))

	templateRules.AttachEarly(rules.ExcludeComments)
	templateRules.Attach(rules.FnJoin)
	templateRules.Attach(rules.MakeFnGetAtt(&sources, &templateRules))
	templateRules.Attach(rules.MakeRef(&sources, &templateRules))

	processed := template.Process(t, &templateRules)

	var ok bool
	for _, component := range components {
		for subcomponentType, subcomponents := range *component { // eg: "Resource", "Parameter", etc
			for subcomponentKey, subcomponent := range subcomponents.(map[string]interface{}) {
				if _, ok = processed.(map[string]interface{})[subcomponentType]; !ok {
					processed.(map[string]interface{})[subcomponentType] = interface{}(subcomponents)
				} else {
					processed.(map[string]interface{})[subcomponentType].(map[string]interface{})[subcomponentKey] = subcomponent
				}
			}
		}
	}

	switch outputWhat.Get() {
	case OutputTemplate:
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(processed)
	case OutputCredentials:
		var credentials []interface{}
		credentialMap := make(map[string]interface{})
		for filename, input := range inputParameters.Sources() {
			credentialMap = template.Process(input, &templateRules).(map[string]interface{})
			credentialMap["$comment"] = map[string]interface{}{"filename": filename}
			credentials = append(credentials, credentialMap)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.Encode(credentials)
	case OutputParameters:
		var templateParameters interface{}
		var templateParameterMap map[string]interface{}
		parameters := []cloudformation.Parameter{}
		if templateParameters, ok = processed.(map[string]interface{})["Parameters"]; ok {
			if templateParameterMap, ok = templateParameters.(map[string]interface{}); ok {
				for name, _ := range templateParameterMap {
					value, ok := sources.Get([]string{name})
					if !ok {
						continue
					}

					value = template.Process(value, &templateRules)

					parameters = append(parameters, func(name string, value interface{}) cloudformation.Parameter {
						stringval := fmt.Sprintf("%s", value)
						boolval := false
						return cloudformation.Parameter{
							ParameterKey:     &name,
							ParameterValue:   &stringval,
							UsePreviousValue: &boolval,
						}
					}(name, value))
				}
			}
		}

		enc := json.NewEncoder(os.Stdout)
		enc.Encode(parameters)
	}
}
