package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type ProfileData struct {
	Name     string
	Count    int
	Duration time.Duration
}

var PROFILE_MUTEX = sync.Mutex{}
var PROFILE = make(map[string]*ProfileData)

func Profile(name string, duration time.Duration) {
	PROFILE_MUTEX.Lock()
	defer PROFILE_MUTEX.Unlock()
	if _, ok := PROFILE[name]; !ok {
		PROFILE[name] = &ProfileData{name, 0, 0}
	}
	p := PROFILE[name]
	p.Count++
	p.Duration += duration
}

func PrintProfile(reset bool) {
	for k, v := range PROFILE {
		log.Infof("%s: %+v", k, v)
	}
	if reset {
		PROFILE = make(map[string]*ProfileData)
	}
}

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

func StringInSlice(a string, list []string) (ret bool) {
	/*start := time.Now()
	defer func() {
		Profile("StringInSlice", time.Now().Sub(start))
		if len(list) > 3 {
			log.Infof("StringInSlice: %d -> %t", len(list), ret)
		}
	}()*/
	for _, b := range list {
		if b == a {
			ret = true
			return
		}
	}
	ret = false
	return
}

func Int64InSlice(a int64, list []int64) (ret bool) {
	/*start := time.Now()
	defer func() {
		Profile("Int64InSlice", time.Now().Sub(start))
		if len(list) > 3 {
			defer log.Infof("Int64InSlice: %d -> %t", len(list), ret)
		}
	}()*/
	for _, b := range list {
		if b == a {
			ret = true
			return
		}
	}
	ret = false
	return
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

func IntersectSorted(first []string, second []string) (ret []string) {
	/*start := time.Now()
	defer func() {
		Profile("IntersectSorted", time.Now().Sub(start))
		log.Infof("IntersectSorted: %d %d -> %d", len(first), len(second), len(ret))
	}()*/
	firstIndex := 0
	secondIndex := 0
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
	return
}

func UnionSorted(first []string, second []string) (ret []string) {
	/*start := time.Now()
	defer func() {
		Profile("UnionSorted", time.Now().Sub(start))
		defer log.Infof("UnionSorted: %d, %d -> %d", len(first), len(second), len(ret))
	}()*/
	firstIndex := 0
	secondIndex := 0
	for firstIndex < len(first) || secondIndex < len(second) {
		if firstIndex == len(first) {
			ret = append(ret, second[secondIndex:]...)
			return
		} else if secondIndex == len(second) {
			ret = append(ret, first[firstIndex:]...)
			return
		} else {
			if cmp := strings.Compare(first[firstIndex], second[secondIndex]); cmp == 0 {
				firstIndex++
			} else if cmp == 1 {
				ret = append(ret, second[secondIndex])
				secondIndex++
			} else if cmp == -1 {
				ret = append(ret, first[firstIndex])
				firstIndex++
			}
		}
	}
	return
}

func Filter(ss []string, test func(string) bool) []string {
	ret := make([]string, 0, len(ss))
	/*start := time.Now()
	defer func() {
		Profile("Filter", time.Now().Sub(start))
		defer log.Infof("Filter: %d -> %d", len(ss), len(ret))
	}()*/
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
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
