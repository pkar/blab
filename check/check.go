package check

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}, msg ...string) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		if len(msg) > 0 {
			fmt.Printf(
				"\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n%v\n\n",
				filepath.Base(file), line, exp, act, msg)
		} else {
			fmt.Printf(
				"\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
				filepath.Base(file), line, exp, act)
		}
		tb.FailNow()
	}
}

// EqualJSON checks if two JSON are equal.
func EqualJSON(tb testing.TB, s1, s2 string) {
	var o1 interface{}
	var o2 interface{}

	if err := json.Unmarshal([]byte(s1), &o1); err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\tunexpectedly equal: %#v\n\n\t\033[39m\n\n", filepath.Base(file), line, err)
		tb.FailNow()
	}
	if err := json.Unmarshal([]byte(s2), &o2); err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\tunexpectedly equal: %#v\n\n\t\033[39m\n\n", filepath.Base(file), line, err)
		tb.FailNow()
	}

	if !reflect.DeepEqual(o1, o2) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\033[31m%s:%d:\n\n\texp: %s\n\n\tgot: %s\033[39m\n\n",
			filepath.Base(file), line, s1, s2)
		tb.FailNow()
	}
}
