package example

import (
	"github.com/xdg/testy"
	"testing"
)

type pair struct {
	x float32
	y float32
}

func TestExample2(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	is.True(1+1 == 3)                          // Line 17
	is.False(2 == 2)                           // Line 18
	is.Equal(1, 2)                             // Line 19
	is.Equal(1.0, 1)                           // Line 20
	is.Equal("foo\tbar", "foo\tbaz")           // Line 21
	is.Equal(true, false)                      // Line 22
	is.Equal(&pair{1.0, 1.0}, &pair{1.1, 1.0}) // Line 23
	is.Unequal(42, 42)                         // Line 24
}
