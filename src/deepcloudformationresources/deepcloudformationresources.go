package deepcloudformationresources

import (
	"fallbackmap"
	"fmt"
	"regexp"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func NewDeepCloudFormationResources(region string) *DeepCloudFormationResources {
	return &DeepCloudFormationResources{
		Region: region,
		cache: map[string]fallbackmap.Deep{},
	}
}

type DeepCloudFormationResources struct {
	Region string
	cache map[string]fallbackmap.Deep
}

func isValidStackName(candidate string) bool {
	did_match, err := regexp.MatchString("^[a-zA-Z][-a-zA-Z0-9]*$", candidate)
	return err == nil && did_match && len(candidate) < 128
}

func isValidResourceName(candidate string) bool {
	did_match, err := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9]*$", candidate)
	return err == nil && did_match
}

func (catalogue *DeepCloudFormationResources) Get(path []string) (interface{}, bool) {
	// path should always be in the form: [StackName, "Outputs", OutputParameter]
	if len(path) != 3 || !isValidStackName(path[0]) || path[1] != "Resources" || !isValidResourceName(path[2]) {
		return nil, false
	}

	if catalogue.cache != nil {
		cached, ok := catalogue.cache[path[0]]
		if ok {
			return cached.Get(path[1:])
		}
	}

	svc := cloudformation.New(&aws.Config{Region: aws.String(catalogue.Region)})
	response, err := svc.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{
		StackName: &path[0],
	})

	if err != nil {
		// FIXME: need a better way to handle this.
		// Right now we just ignore the error, as we may not have actually wanted a
		// Stack Resource
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		return nil, false
	}

	resources := map[string]interface{}{}
	for _, resource := range response.StackResources {
		resources[*resource.LogicalResourceId] = *resource.PhysicalResourceId
	}

	deep := fallbackmap.DeepMap(map[string]interface{}{
		"Resources": resources,
	})
	catalogue.cache[path[0]] = deep

	return deep.Get(path[1:])
}
