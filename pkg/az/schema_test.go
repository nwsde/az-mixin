package az

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/ghodss/yaml" // We are not using go-yaml because of serialization problems with jsonschema, don't use this library elsewhere
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestMixin_PrintSchema(t *testing.T) {
	m := NewTestMixin(t)

	err := m.PrintSchema()
	require.NoError(t, err)

	gotSchema := m.TestContext.GetOutput()

	wantSchema, err := os.ReadFile("schema/schema.json")
	require.NoError(t, err)

	assert.Equal(t, string(wantSchema), gotSchema)
}

func TestMixin_ValidateSchema(t *testing.T) {
	// Load the mixin schema
	schemaLoader := gojsonschema.NewStringLoader(schema)

	testcases := []struct {
		name      string
		file      string
		wantError string
	}{
		{"install", "testdata/install-input.yaml", ""},
		{"upgrade", "testdata/upgrade-input.yaml", ""},
		{"invoke", "testdata/invoke-input.yaml", ""},
		{"uninstall", "testdata/uninstall-input.yaml", ""},
		{"invalid property", "testdata/invalid-input.yaml", "Additional property args is not allowed"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Read the mixin input as a go dump
			mixinInputB, err := os.ReadFile(tc.file)
			require.NoError(t, err)
			mixinInputMap := make(map[string]interface{})
			err = yaml.Unmarshal(mixinInputB, &mixinInputMap)
			require.NoError(t, err)
			mixinInputLoader := gojsonschema.NewGoLoader(mixinInputMap)

			// Validate the manifest against the schema
			result, err := gojsonschema.Validate(schemaLoader, mixinInputLoader)
			require.NoError(t, err)

			if tc.wantError == "" {
				assert.True(t, result.Valid())
				assert.Empty(t, result.Errors())
			} else {
				assert.False(t, result.Valid())
				assert.Contains(t, fmt.Sprintf("%v", result.Errors()), tc.wantError)
			}
		})
	}
}

func TestMixin_CheckSchema(t *testing.T) {
	// Long term it would be great to have a helper function in Porter that a mixin can use to check that it meets certain interfaces
	// check that certain characteristics of the schema that Porter expects are present
	// Once we have a mixin library, that would be a good place to package up this type of helper function
	var schemaMap map[string]interface{}
	err := json.Unmarshal([]byte(schema), &schemaMap)
	require.NoError(t, err, "could not unmarshal the schema into a map")

	t.Run("mixin configuration", func(t *testing.T) {
		// Check that mixin config is defined, and has all the supported fields
		configSchema, err := jsonpath.Get("$.definitions.config", schemaMap)
		require.NoError(t, err, "could not find the mixin config schema declaration")
		_, err = jsonpath.Get("$.properties.az.properties.clientVersion", configSchema)
		require.NoError(t, err, "clientVersion was not included in the mixin config schema")
	})

	// Check that schema are defined for each action
	actions := []string{"install", "upgrade", "invoke", "uninstall"}
	for _, action := range actions {
		t.Run("supports "+action, func(t *testing.T) {
			actionPath := fmt.Sprintf("$.definitions.%sStep", action)
			_, err := jsonpath.Get(actionPath, schemaMap)
			require.NoErrorf(t, err, "could not find the %sStep declaration", action)
		})
	}

	// Check that the invoke action is registered
	additionalSchema, err := jsonpath.Get("$.additionalProperties.items", schemaMap)
	require.NoError(t, err, "the invoke action was not registered in the schema")
	require.Contains(t, additionalSchema, "$ref")
	invokeRef := additionalSchema.(map[string]interface{})["$ref"]
	require.Equal(t, "#/definitions/invokeStep", invokeRef, "the invoke action was not registered correctly")
}
