// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package validate

import "github.com/go-openapi/spec"

const (
	extensionNullable   = "x-nullable"
	extensionIsNullable = "x-isnullable"
)

// isNullable reports whether a schema accepts JSON null values.
//
// Swagger 2.0 uses the x-nullable and x-isnullable vendor extensions for this.
// OpenAPI 3 uses the nullable keyword, mapped to spec.Schema.Nullable.
func isNullable(nullable bool, extensions spec.Extensions) bool {
	if b, ok := nullableExtension(extensions); ok {
		return b
	}

	return nullable
}

func nullableExtension(extensions spec.Extensions) (bool, bool) {
	if extensions == nil {
		return false, false
	}

	if b, ok := extensions[extensionNullable].(bool); ok {
		return b, true
	}

	if b, ok := extensions[extensionIsNullable].(bool); ok {
		return b, true
	}

	return false, false
}
