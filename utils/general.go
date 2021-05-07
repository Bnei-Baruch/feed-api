package utils

import (
	"database/sql"
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

func Int64InSlice(a int64, list []int64) bool {
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

func IntersectSorted(first []string, second []string) []string {
	firstIndex := 0
	secondIndex := 0
	ret := []string(nil)
	for firstIndex < len(first) && secondIndex < len(second) {
		if cmp := strings.Compare(first[firstIndex], second[secondIndex]); cmp == 0 {
			ret = append(ret, second[secondIndex])
			secondIndex++
			firstIndex++
		} else if cmp == 1 {
			secondIndex++
		} else if cmp == -1 {
			firstIndex++
		}
	}
	return ret
}

func UnionSorted(first []string, second []string) []string {
	firstIndex := 0
	secondIndex := 0
	ret := []string(nil)
	for firstIndex < len(first) || secondIndex < len(second) {
		if firstIndex == len(first) {
			ret = append(ret, second...)
		} else if secondIndex == len(second) {
			ret = append(ret, first...)
		} else {
			cmp := strings.Compare(first[firstIndex], second[secondIndex])
			for cmp == 0 {
				secondIndex++
				firstIndex++
				cmp = strings.Compare(first[firstIndex], second[secondIndex])
			}
			if cmp == 1 {
				ret = append(ret, second[secondIndex])
				secondIndex++
			} else /* cmp == -1 */ {
				ret = append(ret, first[firstIndex])
				firstIndex++
			}
		}
	}
	return ret
}

func Filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func NullStringSliceToStringSlice(in []sql.NullString) []string {
	out := []string(nil)
	for _, nullString := range in {
		if nullString.Valid {
			out = append(out, nullString.String)
		}
	}
	return out
}
