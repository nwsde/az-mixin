package az

import (
	"io/ioutil"
	"testing"

	"get.porter.sh/porter/pkg/exec/builder"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

func TestMixin_UnmarshalStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0].TypedCommand.(*UserCommand)
	assert.Equal(t, "Summon Minion", step.Description)
	assert.NotEmpty(t, step.Outputs)
	assert.Equal(t, Output{Name: "VICTORY", JsonPath: "$Id"}, step.Outputs[0])

	require.Len(t, step.Arguments, 1)
	assert.Equal(t, "man-e-faces", step.Arguments[0])

	require.Len(t, step.Flags, 1)
	assert.Equal(t, builder.NewFlag("species", "human"), step.Flags[0])

	assert.Equal(t, false, step.SuppressOutput)
	assert.Equal(t, false, step.SuppressesOutput())
}

func TestStep_GetFlags(t *testing.T) {
	s := UserCommand{}

	f := s.GetFlags()

	require.Len(t, f, 1, "Flags should always have at least 1 entry: --output")
	assert.Equal(t, builder.NewFlag("output", "json"), f[0])
}

func TestStep_SuppressesOutput(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input-suppress-output.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0].TypedCommand.(*UserCommand)
	assert.Equal(t, true, step.SuppressOutput)
	assert.Equal(t, true, step.SuppressesOutput())
}
