package example_test

import (
	"testing"

	"github.com/csgriffis/go-library-template/pkg/example"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	var tt = []struct {
		a, b int
		want int
	}{
		{1, 2, 3},
		{3, 4, 7},
	}

	for _, tc := range tt {
		got := example.Add(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Add(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}
