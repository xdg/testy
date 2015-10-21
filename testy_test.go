// Copyright 2015 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package testy_test

import (
	"github.com/xdg/testy"
	"regexp"
	"testing"
)

func TestHelpers(t *testing.T) {
	mock := &testing.T{}
	test := testy.New(mock)

	// not failures
	test.True(true)
	test.False(false)

	// failures on lines 24+; check later
	test.True(false)
	test.False(true)
	test.Equal(1, 2)
	test.Unequal("foo", "foo")
	test.Unequal(true, true)

	fc := test.FailCount()
	expect := 5
	if fc != expect {
		t.Errorf("Incorrect FailCount. Got %d, but expected %d", fc, expect)
	}

	output := test.Output()
	if ok, _ := regexp.MatchString("testy_test.go:24: Expression was not true", output[0]); !ok {
		t.Errorf("True() had wrong error message: '%s'", output[0])
	}
	if ok, _ := regexp.MatchString("testy_test.go:25: Expression was not false", output[1]); !ok {
		t.Errorf("False() had wrong error message: '%s'", output[1])
	}
	if ok, _ := regexp.MatchString("testy_test.go:26: Values were not equal", output[2]); !ok {
		t.Errorf("Equal() had wrong error message: '%s'", output[2])
	}
	if ok, _ := regexp.MatchString(`(?m)^\s+Got: 1 \(int\)`, output[2]); !ok {
		t.Errorf("Equal() had wrong 'got' message: '%s'", output[2])
	}
	if ok, _ := regexp.MatchString(`(?m)^\s+Wanted: 2 \(int\)`, output[2]); !ok {
		t.Errorf("Equal() had wrong 'got' message: '%s'", output[2])
	}
	if ok, _ := regexp.MatchString("testy_test.go:27: Values were not unequal", output[3]); !ok {
		t.Errorf("Unequal() had wrong error message: '%s'", output[3])
	}
	if ok, _ := regexp.MatchString(`(?m)^\s+Got: "foo"`, output[3]); !ok {
		t.Errorf("Unequal() had wrong 'got' message: '%s'", output[3])
	}
	if ok, _ := regexp.MatchString("testy_test.go:28: Values were not unequal", output[4]); !ok {
		t.Errorf("Unequal() had wrong error message: '%s'", output[4])
	}
	if ok, _ := regexp.MatchString(`(?m)^\s+Got: true`, output[4]); !ok {
		t.Errorf("Unequal() had wrong 'got' message: '%s'", output[4])
	}
}

func TestLabelUplevel(t *testing.T) {
	mock := &testing.T{}
	test := testy.New(mock)

	checkEven(test, 3) // Line 70

	output := test.Output()
	if ok, _ := regexp.MatchString(`testy_test.go:70: Testing 3: Value is not even`, output[0]); !ok {
		t.Errorf("checkEven() had wrong error message: '%s'", output[0])
	}

}

func checkEven(is *testy.T, n int) {
	is = is.Uplevel(1).Label("Testing", n)
	if n%2 != 0 {
		is.Error("Value is not even")
	}
}

func TestFail(t *testing.T) {
	mock := &testing.T{}
	test := testy.New(mock)

	if test.FailCount() != 0 {
		t.Errorf("FailCount didn't start at zero")
	}

	test.Fail()

	if !test.Failed() {
		t.Errorf("Fail() not recorded in facade")
	}

	if !mock.Failed() {
		t.Errorf("Fail() not recorded in test object")
	}
}

func TestError(t *testing.T) {
	mock := &testing.T{}
	test := testy.New(mock)

	test.Error("one", "two")

	if !test.Failed() {
		t.Errorf("Error() not recorded in facade")
	}

	if !mock.Failed() {
		t.Errorf("Error() not recorded in test object")
	}

	output := test.Output()
	if ok, _ := regexp.MatchString(`testy_test.go:\d+: one two`, output[0]); !ok {
		t.Errorf("Error() had wrong error message: '%s'", output[0])
	}
}

func TestErrorf(t *testing.T) {
	mock := &testing.T{}
	test := testy.New(mock)

	test.Errorf("%s %d", "three", 4)

	if !test.Failed() {
		t.Errorf("Errorf() not recorded in facade")
	}

	if !mock.Failed() {
		t.Errorf("Errorf() not recorded in test object")
	}

	output := test.Output()
	if ok, _ := regexp.MatchString(`testy_test.go:\d+: three 4`, output[0]); !ok {
		t.Errorf("Errorf() had wrong error message: '%s'", output[0])
	}
}

func TestLogging(t *testing.T) {
	mock := &testing.T{}
	test := testy.NewCase(mock, "Logging test")

	test.Log("one", "two")
	test.Logf("%s %d", "three", 4)
	log := test.Done()

	// All tests pass case
	if ok, _ := regexp.MatchString(`^Logging test: all tests passed`, log); !ok {
		t.Errorf("Done() had wrong summary: '%s'", log)
	}
	if ok, _ := regexp.MatchString(`testy_test.go:\d+: one two`, log); !ok {
		t.Errorf("Log() message not seen: '%s'", log)
	}
	if ok, _ := regexp.MatchString(`testy_test.go:\d+: three 4`, log); !ok {
		t.Errorf("Logf() message not seen: '%s'", log)
	}

	// 1 tests fails case
	test.Error("inject error")
	log = test.Done()
	if ok, _ := regexp.MatchString(`^Logging test: 1 test failed`, log); !ok {
		t.Errorf("Done() had wrong summary: '%s'", log)
	}

	// 2 tests fail case
	test.Error("inject error")
	log = test.Done()
	if ok, _ := regexp.MatchString(`^Logging test: 2 tests failed`, log); !ok {
		t.Errorf("Done() had wrong summary: '%s'", log)
	}
}
