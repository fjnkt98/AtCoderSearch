package generate

import (
	"reflect"
	"testing"
)

func TestGenerateDocumentOptions(t *testing.T) {
	cases := []struct {
		name     string
		options  []option
		expected GenerateDocumentOptions
	}{
		{name: "default", options: nil, expected: GenerateDocumentOptions{ChunkSize: 0, Concurrent: 0}},
		{name: "with chunk size", options: []option{WithChunkSize(6)}, expected: GenerateDocumentOptions{ChunkSize: 6, Concurrent: 0}},
		{name: "with concurrent", options: []option{WithConcurrent(10)}, expected: GenerateDocumentOptions{ChunkSize: 0, Concurrent: 10}},
		{name: "with chunk size and concurrent", options: []option{WithChunkSize(4), WithConcurrent(4)}, expected: GenerateDocumentOptions{ChunkSize: 4, Concurrent: 4}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := &GenerateDocumentOptions{}
			for _, opt := range tt.options {
				opt(actual)
			}

			if !reflect.DeepEqual(tt.expected, *actual) {
				t.Errorf("expected \n%+v\n, but got \n%+v\n", tt.expected, *actual)
			}
		})
	}
}
