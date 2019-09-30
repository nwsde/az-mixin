package az

import (
	"github.com/deislabs/porter/pkg/exec/builder"
)

var _ builder.ExecutableAction = Action{}

type Action struct {
	Steps []Steps // using UnmarshalYAML so that we don't need a custom type per action
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - az: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var steps []Steps
	results, err := builder.UnmarshalAction(unmarshal, &steps)
	if err != nil {
		return err
	}

	// Go doesn't have generics, nothing to see here...
	for _, result := range results {
		step := result.(*[]Steps)
		a.Steps = append(a.Steps, *step...)
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	// Go doesn't have generics, nothing to see here...
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

type Steps struct {
	Step `yaml:"az"`
}

var _ builder.ExecutableStep = Step{}
var _ builder.StepWithOutputs = Step{}

type Step struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Arguments   []string      `yaml:"arguments,omitempty"`
	Flags       builder.Flags `yaml:"flags,omitempty"`
	Outputs     []Output      `yaml:"outputs,omitempty"`
}

func (s Step) GetCommand() string {
	return "az"
}

func (s Step) GetArguments() []string {
	return s.Arguments
}

func (s Step) GetFlags() builder.Flags {
	return append(s.Flags, builder.NewFlag("output", "json"))
}

func (s Step) GetOutputs() []builder.Output {
	// Go doesn't have generics, nothing to see here...
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}

var _ builder.OutputJsonPath = Output{}
var _ builder.OutputFile = Output{}

type Output struct {
	Name string `yaml:"name"`

	// See https://porter.sh/mixins/exec/#outputs
	JsonPath string `yaml:"jsonPath,omitempty"`
	FilePath string `yaml:"path,omitempty"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}

func (o Output) GetFilePath() string {
	return o.FilePath
}
