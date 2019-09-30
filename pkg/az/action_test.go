package az

import (
	"io/ioutil"
	"testing"

	"github.com/deislabs/porter/pkg/exec/builder"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestMixin_UnmarshalStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/step-input.yaml")
	require.NoError(t, err)

	var action Action
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 1)

	step := action.Steps[0]
	assert.Equal(t, "Summon Minion", step.Description)
	assert.NotEmpty(t, step.Outputs)
	assert.Equal(t, Output{Name: "VICTORY", JsonPath: "$Id"}, step.Outputs[0])

	require.Len(t, step.Arguments, 1)
	assert.Equal(t, "man-e-faces", step.Arguments[0])

	require.Len(t, step.Flags, 1)
	assert.Equal(t, builder.NewFlag("species", "human"), step.Flags[0])
}

func TestStep_GetFlags(t *testing.T) {
	s := Step{}

	f := s.GetFlags()

	require.Len(t, f, 1, "Flags should always have at least 1 entry: --output")
	assert.Equal(t, builder.NewFlag("output", "json"), f[0])
}
