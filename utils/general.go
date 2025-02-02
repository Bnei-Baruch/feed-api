package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
)

type ProfileData struct {
	Name     string
	Count    int
	Duration time.Duration
}

var PROFILE_MUTEX = sync.Mutex{}
var PROFILE = make(map[string]*ProfileData)

func StringKeys(s interface{}) []string {
	m := reflect.ValueOf(s)
	if m.Kind() != reflect.Map {
		panic(fmt.Sprintf("Expected map! got: %+v", m.Kind()))
	}
	out := []string(nil)
	for _, key := range m.MapKeys() {
		out = append(out, key.String())
	}
	return out
}

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

type loggingFFunc func(format string, args ...interface{})

func PrintProfile(reset bool) {
	PROFILE_MUTEX.Lock()
	defer PROFILE_MUTEX.Unlock()
	logf := log.Debugf
	for _, pd := range PROFILE {
		if pd.Duration > 5*time.Second {
			logf = log.Infof
			break
		}
	}
	keys := StringKeys(PROFILE)
  if len(keys) > 0 {
    logf("===== Profile =====")
  }
	sort.Strings(keys)
	for _, k := range keys {
		logf("%s: %+v", k, PROFILE[k])
	}
	if reset {
		PROFILE = make(map[string]*ProfileData)
	}
  if len(keys) > 0 {
    logf("===== End Profile =====")
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
		panic(fmt.Sprintf("Expected slice! got: %+v", slice.Kind()))
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

func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MaxInt64(x, y int64) int64 {
	if x >= y {
		return x
	}
	return y
}

func MinDuration(x, y time.Duration) time.Duration {
	if x < y {
		return x
	}
	return y
}

func MaxDuration(x, y time.Duration) time.Duration {
	if x >= y {
		return x
	}
	return y
}

func StringInSlice(a string, list []string) (ret bool) {
	/*start := time.Now()
	defer func() {
		Profile("StringInSlice", time.Now().Sub(start))
		//if len(list) > 3 {
		//	log.Infof("StringInSlice: %d -> %t", len(list), ret)
		//}
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
	start := time.Now()
	defer func() { Profile("IntersectSorted", time.Now().Sub(start)) }()
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

// Changes the frist map.
func UnionMaps(a map[string]bool, b map[string]bool) {
	start := time.Now()
	defer func() {
		Profile("UnionMaps", time.Now().Sub(start))
	}()
	for k, _ := range b {
		a[k] = true
	}
}

func UnionSorted(first []string, second []string) (ret []string) {
	start := time.Now()
	defer func() {
		Profile("UnionSorted", time.Now().Sub(start))
		defer log.Debugf("UnionSorted: %d, %d -> %d", len(first), len(second), len(ret))
	}()
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
	start := time.Now()
	defer func() {
		Profile("Filter", time.Now().Sub(start))
		defer log.Debugf("Filter: %d -> %d", len(ss), len(ret))
	}()
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func FilterMap(m map[string]bool, test func(string) bool) {
	startLen := len(m)
	start := time.Now()
	defer func() {
		Profile("FilterMap", time.Now().Sub(start))
		defer log.Debugf("FilterMap: %d -> %d", startLen, len(m))
	}()
	for k, _ := range m {
		if !test(k) {
			delete(m, k)
		}
	}
}

func IntersectMaps(a map[string]bool, b map[string]bool) {
	start := time.Now()
	lenA := len(a)
	defer func() {
		Profile("IntersectMaps", time.Now().Sub(start))
		defer log.Debugf("IntersectMaps: %d %d -> %d", lenA, len(b), len(a))
	}()
	for k, _ := range a {
		if _, ok := b[k]; !ok {
			delete(a, k)
		}
	}
}

func MapFromSlice(ss []string) map[string]bool {
	ret := make(map[string]bool, len(ss))
	for _, s := range ss {
		ret[s] = true
	}
	return ret
}

func SliceFromMap(m map[string]bool) []string {
	ret := make([]string, 0, len(m))
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
}

func CopyStringMap(m map[string]bool) map[string]bool {
	start := time.Now()
	ret := make(map[string]bool, len(m))
	defer func() {
		Profile("CopyStringMap", time.Now().Sub(start))
		defer log.Debugf("CopyStringMap: %d", len(m))
	}()
	for k, v := range m {
		ret[k] = v
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

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func NullStringToValue(s null.String) string {
	if s.Valid {
		return fmt.Sprintf("'%s'", strings.ReplaceAll(s.String, "'", "''"))
	}
	return "NULL"
}

func NullJSONToValue(j null.JSON) string {
	if j.Valid {
		return fmt.Sprintf("'%s'", strings.ReplaceAll(string(j.JSON), "'", "''"))
	}
	return "NULL"
}
