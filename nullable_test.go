// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package validate

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestIsNullable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		nullable  bool
		extension spec.Extensions
		want      bool
	}{
		{
			name:     "nullable keyword",
			nullable: true,
			want:     true,
		},
		{
			name: "x-nullable extension",
			extension: spec.Extensions{
				extensionNullable: true,
			},
			want: true,
		},
		{
			name: "x-isnullable extension",
			extension: spec.Extensions{
				extensionIsNullable: true,
			},
			want: true,
		},
		{
			name: "x-nullable overrides nullable false",
			extension: spec.Extensions{
				extensionNullable: true,
			},
			want: true,
		},
		{
			name:     "x-nullable false overrides nullable true",
			nullable: true,
			extension: spec.Extensions{
				extensionNullable: false,
			},
			want: false,
		},
		{
			name: "x-nullable takes precedence over x-isnullable",
			extension: spec.Extensions{
				extensionNullable:   false,
				extensionIsNullable: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.EqualT(t, tt.want, isNullable(tt.nullable, tt.extension))
		})
	}
}

func TestSchemaValidator_XNullable(t *testing.T) {
	t.Parallel()

	schemaJSON := `{
		"type": "object",
		"required": ["age"],
		"properties": {
			"age": {
				"type": "integer",
				"x-nullable": true
			}
		}
	}`

	schema := new(spec.Schema)
	require.NoError(t, json.Unmarshal([]byte(schemaJSON), schema))

	t.Run("accepts null property value", func(t *testing.T) {
		t.Parallel()
		input := map[string]any{"age": nil}
		require.NoError(t, AgainstSchema(schema, input, strfmt.Default))
	})

	t.Run("accepts integer property value", func(t *testing.T) {
		t.Parallel()
		input := map[string]any{"age": json.Number("42")}
		require.NoError(t, AgainstSchema(schema, input, strfmt.Default))
	})
}

func TestParamValidator_XNullable(t *testing.T) {
	t.Parallel()

	paramJSON := `{
		"name": "count",
		"in": "query",
		"type": "integer",
		"x-nullable": true
	}`

	param := new(spec.Parameter)
	require.NoError(t, json.Unmarshal([]byte(paramJSON), param))

	validator := NewParamValidator(param, strfmt.Default)

	t.Run("accepts null", func(t *testing.T) {
		t.Parallel()
		require.TrueT(t, validator.Validate(nil).IsValid())
	})

	t.Run("accepts integer", func(t *testing.T) {
		t.Parallel()
		require.TrueT(t, validator.Validate(int64(7)).IsValid())
	})
}
