package util

import "testing"

func TestRestoreAbstract(t *testing.T) {
	tests := []struct {
		name string
		idx  map[string][]int
		want string
	}{
		{
			name: "nil index",
			idx:  nil,
			want: "",
		},
		{
			name: "empty index",
			idx:  map[string][]int{},
			want: "",
		},
		{
			name: "single word",
			idx:  map[string][]int{"hello": {0}},
			want: "hello",
		},
		{
			name: "ordered words",
			idx:  map[string][]int{"hello": {0}, "world": {1}},
			want: "hello world",
		},
		{
			name: "unordered positions",
			idx:  map[string][]int{"world": {1}, "hello": {0}},
			want: "hello world",
		},
		{
			name: "repeated word at different positions",
			idx:  map[string][]int{"the": {0, 3}, "cat": {1}, "sat": {2}},
			want: "the cat sat the",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RestoreAbstract(tt.idx)
			if got != tt.want {
				t.Errorf("RestoreAbstract() = %q, want %q", got, tt.want)
			}
		})
	}
}
