package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// panic if err != nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func ToInterfaceSlice(s interface{}) []interface{} {
	slice := reflect.ValueOf(s)
	if slice.Kind() != reflect.Slice {
		panic("Expected slice!")
	}
	c := slice.Len()
	out := make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out
}

func ConvertArgsInt64(args []int64) []interface{} {
	c := make([]interface{}, len(args))
	for i := range args {
		c[i] = args[i]
	}
	return c
}

func ConvertArgsString(args []string) []interface{} {
	c := make([]interface{}, len(args))
	for i := range args {
		c[i] = args[i]
	}
	return c
}

// Like math.Min for int
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) int {
	if x >= y {
		return x
	}
	return y
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func InClause(prefix string, list []string) string {
	inClause := ""
	if len(list) > 0 {
		quoted := []string(nil)
		for _, uid := range list {
			quoted = append(quoted, fmt.Sprintf("'%s'", uid))
		}
		inClause = fmt.Sprintf("%s (%s)", prefix, strings.Join(quoted, ","))
	}
	return inClause
}
