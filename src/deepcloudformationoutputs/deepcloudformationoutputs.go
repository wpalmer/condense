package deepcloudformationoutputs

import (
	"fallbackmap"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"os"
	"regexp"
)

func NewDeepCloudFormationOutputs(region string) *DeepCloudFormationOutputs {
	return &DeepCloudFormationOutputs{
		Region: region,
		cache:  map[string]fallbackmap.Deep{},
	}
}

type DeepCloudFormationOutputs struct {
	Region string
	cache  map[string]fallbackmap.Deep
}

func isValidStackName(candidate string) bool {
	did_match, err := regexp.MatchString("^[a-zA-Z][-a-zA-Z0-9]*$", candidate)
	return err == nil && did_match && len(candidate) < 128
}

func isValidOutputName(candidate string) bool {
	did_match, err := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9]*$", candidate)
	return err == nil && did_match
}

func (catalogue *DeepCloudFormationOutputs) Get(path []string) (interface{}, bool) {
	// path should always be in the form: [StackName, "Outputs", OutputParameter]
	if len(path) != 3 || !isValidStackName(path[0]) || path[1] != "Outputs" || !isValidOutputName(path[2]) {
		return nil, false
	}

	if catalogue.cache != nil {
		cached, ok := catalogue.cache[path[0]]
		if ok {
			return cached.Get(path[1:])
		}
	}

	svc := cloudformation.New(&aws.Config{Region: aws.String(catalogue.Region)})
	description, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: &path[0],
	})

	if err != nil {
		// FIXME: need a better way to handle this.
		// Right now we just ignore the error, as we may not have actually wanted a
		// Stack Output
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		return nil, false
	}

	if len(description.Stacks) != 1 {
		panic(fmt.Sprintf(
			"Description of [%s] did not return in exactly one Stack",
			path[0],
		))
	}

	outputs := map[string]interface{}{}
	stack := description.Stacks[0]
	for _, output := range stack.Outputs {
		outputs[*output.OutputKey] = *output.OutputValue
	}

	deep := fallbackmap.DeepMap(map[string]interface{}{
		"Outputs": outputs,
	})
	catalogue.cache[path[0]] = deep

	return deep.Get(path[1:])
}
