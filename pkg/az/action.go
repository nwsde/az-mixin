package az

import (
	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []TypedStep // using UnmarshalYAML so that we don't need a custom type per action
}

// MakeSteps builds a slice of Steps for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]TypedStep{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - az: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			steps := result.(*[]TypedStep)
			for _, step := range *steps {
				step.SetAction(actionName)
				a.Steps = append(a.Steps, step)
			}
		}
		break // There is only 1 action
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
	TypedStep `yaml:"az"`
}

type TypedStep struct {
	Description string
	TypedCommand
}

// UnmarshalYAML takes any yaml in this form
// az:
//   description: something
//   COMMAND: # custom... -> make the CustomCommand for us
func (s *TypedStep) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Turn the yaml into a raw map so we can iterate over the values and
	// look for which command was used
	stepMap := map[string]map[string]interface{}{}
	err := unmarshal(&stepMap)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal yaml into a raw az command")
	}

	// get at the values defined under "az"
	step := stepMap["az"]

	// Turn each command into its typed data structure
	for key, value := range step {
		var cmd TypedCommand

		switch key {
		case "description":
			s.Description = value.(string)
			continue
		case "group":
			cmd = &GroupCommand{}
		default: // It's a custom user command
			customCmd := &UserCommand{}
			b, err := yaml.Marshal(step)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(b, customCmd)
			if err != nil {
				return err
			}
			s.TypedCommand = customCmd
			return nil
		}

		// We have a typed command, unmarshal it
		b, err := yaml.Marshal(value)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(b, cmd)
		if err != nil {
			return err
		}

		s.TypedCommand = cmd
		return nil
	}

	return nil
}

type TypedCommand interface {
	SetAction(action string)
	builder.ExecutableStep
	builder.SuppressesOutput
}

var _ TypedCommand = UserCommand{}

type UserCommand struct {
	Name           string        `yaml:"name"`
	Description    string        `yaml:"description"`
	Arguments      []string      `yaml:"arguments,omitempty"`
	Flags          builder.Flags `yaml:"flags,omitempty"`
	Outputs        []Output      `yaml:"outputs,omitempty"`
	SuppressOutput bool          `yaml:"suppress-output,omitempty"`

	// Support custom error handling
	builder.IgnoreErrorHandler `yaml:"ignoreErrors,omitempty"`
}

func (s UserCommand) GetWorkingDir() string {
	return ""
}

var _ builder.ExecutableStep = UserCommand{}
var _ builder.StepWithOutputs = UserCommand{}
var _ builder.SuppressesOutput = UserCommand{}

func (s UserCommand) GetCommand() string {
	return "az"
}

func (s UserCommand) GetArguments() []string {
	return s.Arguments
}

func (s UserCommand) GetFlags() builder.Flags {
	return append(s.Flags, builder.NewFlag("output", "json"))
}

func (s UserCommand) GetOutputs() []builder.Output {
	// Go doesn't have generics, nothing to see here...
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}

func (s UserCommand) SuppressesOutput() bool {
	return s.SuppressOutput
}

func (s UserCommand) SetAction(_ string) {}

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
