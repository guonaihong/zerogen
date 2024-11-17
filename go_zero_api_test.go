package zerogen

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word lowercase",
			input:    "test",
			expected: "test",
		},
		{
			name:     "single word uppercase",
			input:    "TEST",
			expected: "test",
		},
		{
			name:     "camel case",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "pascal case",
			input:    "PascalCase",
			expected: "pascal_case",
		},
		{
			name:     "multiple words",
			input:    "ThisIsATest",
			expected: "this_is_a_test",
		},
		{
			name:     "already snake case",
			input:    "already_snake_case",
			expected: "already_snake_case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.expected {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.expected)
			}
		})
	}
}
