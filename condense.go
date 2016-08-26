package main

import (
	"deepalias"
	"deepcloudformationoutputs"
	"deepcloudformationresources"
	"deepstack"
	"encoding/json"
	"fallbackmap"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"golang.org/x/tools/godoc/vfs"
	"io"
	"lazymap"
	"os"
	"strings"

	"condense/template"
	"condense/template/rules"
)

type inputSource struct {
	filename string
	data     map[string]interface{}
}

type InputsFlag struct {
	inputs  *fallbackmap.FallbackMap
	sources []inputSource
	rules   *template.Rules
}

func NewInputsFlag(rules *template.Rules) InputsFlag {
	return InputsFlag{
		&fallbackmap.FallbackMap{},
		[]inputSource{},
		rules,
	}
}

func (f *InputsFlag) Get() *fallbackmap.FallbackMap {
	return f.inputs
}

func (f InputsFlag) String() string {
	return fmt.Sprintf("%v", f.sources)
}

func (f *InputsFlag) Set(parametersFilename string) (err error) {
	var inputStream io.Reader
	var raw interface{}
	var ok bool
	var gotRaw bool
	var ins []interface{}
	var in map[string]interface{}

	rawJson := strings.NewReader(parametersFilename)
	rawDecoder := json.NewDecoder(rawJson)
	if err := rawDecoder.Decode(&raw); err == nil {
		if ins, ok = raw.([]interface{}); ok {
			gotRaw = true
		} else if in, ok = raw.(map[string]interface{}); ok {
			ins = append(ins, interface{}(in))
			gotRaw = true
		}
	}

	if gotRaw {
		parametersFilename = "[inline]"
	} else {
		// reset, as we cannot rely on the state of raw after decoding into it
		raw = interface{}(nil)

		if inputStream, err = os.Open(parametersFilename); err != nil {
			return err
		}

		inputDecoder := json.NewDecoder(inputStream)
		if err := inputDecoder.Decode(&raw); err != nil {
			return err
		}

		if ins, ok = raw.([]interface{}); !ok {
			if in, ok = raw.(map[string]interface{}); !ok {
				return fmt.Errorf("JSON data does not decode into an array or map")
			}

			ins = append(ins, interface{}(in))
		}
	}

	var i int
	for i, raw = range ins {
		if in, ok = raw.(map[string]interface{}); !ok {
			return fmt.Errorf("JSON data does not decode into a map or array of maps")
		}

		var parametersFilespec string
		if len(ins) == 1 {
			parametersFilespec = parametersFilename
		} else {
			parametersFilespec = fmt.Sprintf("%s[%d]", parametersFilename, i)
		}

		f.inputs.Override(lazymap.NewLazyMap(fallbackmap.DeepMap(in), f.rules))
		f.sources = append(f.sources, inputSource{filename: parametersFilespec, data: in})
	}

	return nil
}

func (f *InputsFlag) Sources() []inputSource {
	return f.sources
}

type OutputWhat int

const (
	OutputTemplate = iota
	OutputParameters
	OutputCredentials
)

type OutputWhatFlag struct {
	what   OutputWhat
	hasKey bool
	key    string
}

func (f OutputWhatFlag) Get() OutputWhatFlag {
	return f
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
	if strings.HasPrefix(input, "credentials:") {
		f.what = OutputCredentials
		f.hasKey = true
		f.key = input[12:]

		return nil
	}

	f.hasKey = false
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

func main() {
	templateRules := template.Rules{}
	inputParameters := NewInputsFlag(&templateRules)
	var templateFilename string
	var outputWhat OutputWhatFlag

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

	sources := fallbackmap.FallbackMap{}
	stack := deepstack.DeepStack{}

	sources.Attach(inputParameters.Get())
	sources.Attach(deepalias.DeepAlias{&stack})
	sources.Attach(deepcloudformationoutputs.NewDeepCloudFormationOutputs("eu-west-1"))
	sources.Attach(deepcloudformationresources.NewDeepCloudFormationResources("eu-west-1"))

	stack.Push(&sources)

	templateRules.AttachEarly(rules.ExcludeComments)
	templateRules.AttachEarly(rules.MakeFnFor(&stack, &templateRules))
	templateRules.AttachEarly(rules.MakeFnWith(&stack, &templateRules))
	templateRules.Attach(rules.FnAdd)
	templateRules.Attach(rules.FnIf)
	templateRules.Attach(rules.FnAnd)
	templateRules.Attach(rules.FnOr)
	templateRules.Attach(rules.FnNot)
	templateRules.Attach(rules.FnEquals)
	templateRules.Attach(rules.FnConcat)
	templateRules.Attach(rules.FnFromEntries)
	templateRules.Attach(rules.FnHasKey)
	templateRules.Attach(rules.FnJoin)
	templateRules.Attach(rules.FnKeys)
	templateRules.Attach(rules.FnLength)
	templateRules.Attach(rules.FnMerge)
	templateRules.Attach(rules.FnMergeDeep)
	templateRules.Attach(rules.FnMod)
	templateRules.Attach(rules.FnSplit)
	templateRules.Attach(rules.FnToEntries)
	templateRules.Attach(rules.FnUnique)
	templateRules.Attach(rules.MakeFnGetAtt(&stack, &templateRules))
	templateRules.Attach(rules.MakeRef(&stack, &templateRules))
	templateRules.Attach(rules.MakeFnIncludeFile(vfs.OS("/"), &templateRules))
	templateRules.Attach(rules.MakeFnIncludeFileRaw(vfs.OS("/")))
	templateRules.Attach(rules.ReduceConditions)

	// First Pass (to collect Parameter names)
	processed := template.Process(t, &templateRules)

	parameterRefs := map[string]interface{}{}
	if processedMap, ok := processed.(map[string]interface{}); ok {
		if processedParameters, ok := processedMap["Parameters"]; ok {
			if processedParametersMap, ok := processedParameters.(map[string]interface{}); ok {
				for parameterName, _ := range processedParametersMap {
					parameterRefs[parameterName] = map[string]interface{}{
						"ParamRef": parameterName,
					}
				}
			}
		}
	}

	stack.Push(fallbackmap.DeepMap(parameterRefs))
	templateRules.Attach(func(path []interface{}, node interface{}) (interface{}, interface{}) {
		key := interface{}(nil)
		if len(path) > 0 {
			key = path[len(path)-1]
		}

		if nodeMap, ok := node.(map[string]interface{}); !ok || len(nodeMap) != 1 {
			return key, node //passthru
		}

		if refName, ok := node.(map[string]interface{})["ParamRef"]; ok {
			return key, interface{}(map[string]interface{}{"Ref": interface{}(refName)})
		}

		return key, node
	})
	processed = template.Process(t, &templateRules)
	stack.PopDiscard()

	switch outputWhat.Get().what {
	case OutputTemplate:
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(processed)
	case OutputCredentials:
		credentials := []interface{}{}
		credentialMap := make(map[string]interface{})
		for _, input := range inputParameters.Sources() {
			if !outputWhat.Get().hasKey || outputWhat.Get().key == input.filename {
				credentialMap = template.Process(input.data, &templateRules).(map[string]interface{})
				credentialMap["$comment"] = map[string]interface{}{"filename": input.filename}
				credentials = append(credentials, credentialMap)
			}
		}

		if len(credentials) == 0 && outputWhat.Get().hasKey {
			panic(fmt.Errorf("No parameters file '%s' was input", outputWhat.Get().key))
		}

		enc := json.NewEncoder(os.Stdout)
		if len(credentials) == 1 {
			enc.Encode(credentials[0])
		} else {
			enc.Encode(credentials)
		}
	case OutputParameters:
		parameters := []cloudformation.Parameter{}

		for name, _ := range parameterRefs {
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

		enc := json.NewEncoder(os.Stdout)
		enc.Encode(parameters)
	}
}
