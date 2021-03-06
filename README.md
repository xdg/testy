[![GoDoc](https://godoc.org/github.com/xdg/testy?status.svg)](https://godoc.org/github.com/xdg/testy)
[![Build Status](https://travis-ci.org/xdg/testy.svg?branch=master)](https://travis-ci.org/xdg/testy)

# Testy – An extensible testing facade

**If Go's standard testing package annoys you, you might like Testy.**

There is a lot to like about Go's [testing](https://golang.org/pkg/testing/)
package.

There are also two extremely annoying things about it:

1. You can't refactor repetitive tests without reporting errors from the
   wrong place in the code.
2. The testing package is locked down tight, preventing trivial solutions
   to the prior problem.

The Go FAQ [suggests](https://golang.org/doc/faq#testing_framework) that
table-driven testing is the way to avoid repetitive test code.  If that
works for you for all situations, fine.

Testy has a different solution.

**Testy implements a facade around the testing package and hijacks its
logging features.**

This means:

* You can report test errors at any level up the call stack.
* You can label all errors in a scope to disambiguate repetitive tests.

The downside is an extra level of log message nesting (which your
editor's quickfix window should ignore, anyway).

Unlike many other [Go testing
libraries](https://github.com/avelino/awesome-go#testing), it doesn't offer
an extensive testing framework.  Unlike some assertion libraries, it
doesn't use stack traces or race-prone print statements.

Testy offers a simple, extensible solution focused on reporting errors
simply and in the right place.  But it does give a few convenient helper
functions for the most common tests you're likely to write.

# Examples

## Using a custom helper function

Consider this [simple example](/_examples/example1_test.go):

```go
package example

import (
	"github.com/xdg/testy"
	"testing"
)

func TestExample1(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }() // Line 10
	is.Error("First failure")            // Line 11
	checkTrue(is, 1+1 == 3)              // Line 12
}

func checkTrue(is *testy.T, cond bool) {
	if !cond {
		is.Uplevel(1).Error("Expression was not true")
	}
}
```

In the `TestExample1` function, the `is` variable wraps the test variable,
`t`.  The `defer` closure schedules test logging output to be delivered to
`t` via `is.Done()` when the test function exits.

When run in Vim, with [vim-go](https://github.com/fatih/vim-go), the
quickfix window looks like this:

```
_examples/example1_test.go|10| TestExample1: 2 tests failed
_examples/example1_test.go|11| First failure
_examples/example1_test.go|12| Expression was not true
```

Note that the `checkTrue` error is reported from the call to `checkTrue` at
line 12, not from inside the `checkTrue` function.  The `Uplevel` method in
`checkTrue` tells Testy to report the error one level up the call stack.

## Using Testy helpers

The `checkTrue` pattern is so common that testing true and false are
built-in to Testy as `True` and `False`.  There are also `Equal` and
`Unequal` helpers that use `reflect.DeepEqual` but provide diagnostic
output on subsequent lines:

```go
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
```

The diagnostic output quotes strings and indicates types where necessary
to disambiguate.  For example:

```
_examples/example2_test.go|15| TestExample2: 8 tests failed
_examples/example2_test.go|17| Expression was not true
_examples/example2_test.go|18| Expression was not false
_examples/example2_test.go|19| Values were not equal:
|| 			   Got: 1 (int)
|| 			Wanted: 2 (int)
_examples/example2_test.go|20| Values were not equal:
|| 			   Got: 1 (float64)
|| 			Wanted: 1 (int)
_examples/example2_test.go|21| Values were not equal:
|| 			   Got: "foo\tbar"
|| 			Wanted: "foo\tbaz"
_examples/example2_test.go|22| Values were not equal:
|| 			   Got: true
|| 			Wanted: false
_examples/example2_test.go|23| Values were not equal:
|| 			   Got: &{1 1} (*example.pair)
|| 			Wanted: &{1.1 1} (*example.pair)
_examples/example2_test.go|24| Values were not unequal:
|| 			   Got: 42 (int)
```

## Using error labels

To prefix error messages with some descriptive text, you can use the
`Label` method like this:

```go
package example

import (
	"github.com/xdg/testy"
	"testing"
)

func TestExample3(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	for i := 1; i <= 5; i++ {
		is.Label("Checking", i).True(i == 3) // Line 13
	}
}
```

The output with labels looks like this, making it clear which tests failed:


```
_examples/example3_test.go|10| TestExample3: 4 tests failed
_examples/example3_test.go|13| Checking 1: Expression was not true
_examples/example3_test.go|13| Checking 2: Expression was not true
_examples/example3_test.go|13| Checking 4: Expression was not true
_examples/example3_test.go|13| Checking 5: Expression was not true
```

## Combining Uplevel and Label in a new facade

Because `Uplevel` and `Label` just return new facades, you can chain them
at the start of a helper function and modify all subsequent logging:

```go
package example

import (
	"github.com/xdg/testy"
	"testing"
)

func TestExample4(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	for i := -1; i <= 2; i++ {
		checkEvenPositive(is, i) // Line 13
	}
}

func checkEvenPositive(is *testy.T, n int) {
	// replace 'is' with a labeled, upleveled equivalent
	is = is.Uplevel(1).Label("Testing", n)

	if n < 1 {
		is.Error("was not positive")
	}
	if n%2 != 0 {
		is.Error("was not even")
	}
}
```

This lets you write test helpers that report errors where they are
called (line 13 in this case), but with detailed errors you can
tie back to the original input data:

```
_examples/example4_test.go|10| TestExample4: 4 tests failed
_examples/example4_test.go|13| Testing -1: was not positive
_examples/example4_test.go|13| Testing -1: was not even
_examples/example4_test.go|13| Testing 0: was not positive
_examples/example4_test.go|13| Testing 1: was not even
```

# Copyright and License

Copyright 2015 by David A. Golden. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
