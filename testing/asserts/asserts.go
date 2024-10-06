// Package asserts provides a set of assertion functions for use in Go tests.
// It allows developers to write more expressive and readable tests with detailed failure messages.
package asserts

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/turtak/go-kit/stacktrace"
)

var (
	// mockTesting is used internally to mock test failures without calling t.Error.
	mockTesting bool

	// mockTestMessage stores the message when mockTesting is true.
	mockTestMessage string = ""

	// stacktraceConfig holds the configuration for stack trace generation
	stacktraceConfig = &stacktrace.Config{
		BufferSize: 2048,
		SkipFrames: 2,
	}
)

// failTest reports a test failure and prints the stack trace.
// If mockTesting is true, it stores the error message without stopping the test.
func failTest(t *testing.T, msg string) {
	if mockTesting {
		mockTestMessage = msg
		return
	}
	stackTrace := stacktrace.NewStackTrace(stacktraceConfig)
	fmt.Printf("--- Stack trace ---\n%s\n-------------------\n", stackTrace.Frames().String())
	t.Error(msg)
}

// isNil checks whether the given value is nil.
func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// isEmpty determines whether the specified object is considered empty.
func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

// compareNumeric compares two numeric values with a small epsilon for float comparisons.
func compareNumeric(a, b any) (int, error) {
	aFloat, errA := toFloat64(a)
	bFloat, errB := toFloat64(b)
	if errA != nil || errB != nil {
		return 0, fmt.Errorf("unsupported numeric types: %T vs %T", a, b)
	}
	diff := aFloat - bFloat
	epsilon := 1e-9
	switch {
	case math.Abs(diff) < epsilon:
		return 0, nil
	case diff > 0:
		return 1, nil
	default:
		return -1, nil
	}
}

// toFloat64 converts a numeric value to float64.
func toFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(val).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(val).Uint()), nil
	case float32, float64:
		return reflect.ValueOf(val).Float(), nil
	default:
		return 0, fmt.Errorf("unsupported type for numeric comparison: %T", v)
	}
}

// Equal asserts that two values are equal using reflect.DeepEqual.
// It fails the test if the values are not equal.
func Equal(t *testing.T, expected, actual any) {
	if !reflect.DeepEqual(expected, actual) {
		failTest(t, fmt.Sprintf("values not equal: expected: %v actual: %v", expected, actual))
	}
}

// NotEqual asserts that two values are not equal using reflect.DeepEqual.
// It fails the test if the values are equal.
func NotEqual(t *testing.T, notExpected, actual any) {
	if reflect.DeepEqual(notExpected, actual) {
		failTest(t, fmt.Sprintf("values unexpectedly equal: not expected: %v actual: %v", notExpected, actual))
	}
}

// Nil asserts that a value is nil.
// It fails the test if the value is not nil.
func Nil(t *testing.T, actual any) {
	if !isNil(actual) {
		failTest(t, fmt.Sprintf("expected nil, but got: %v", actual))
	}
}

// NotNil asserts that a value is not nil.
// It fails the test if the value is nil.
func NotNil(t *testing.T, value any) {
	if isNil(value) {
		failTest(t, "expected non-nil value, but got nil")
	}
}

// NotEmpty asserts that a value is not empty.
// It fails the test if the value is empty.
func NotEmpty(t *testing.T, value any) {
	if isEmpty(value) {
		failTest(t, fmt.Sprintf("expected non-empty value, but got empty: %v", value))
	}
}

// Empty asserts that a value is empty.
// It fails the test if the value is not empty.
func Empty(t *testing.T, value any) {
	if !isEmpty(value) {
		failTest(t, fmt.Sprintf("expected empty value, but got: %v", value))
	}
}

// NoError asserts that an error is nil.
// It fails the test if the error is not nil.
func NoError(t *testing.T, err error) {
	if err != nil {
		failTest(t, fmt.Sprintf("unexpected error: %v", err))
	}
}

// Error asserts that an error is not nil.
// It fails the test if the error is nil.
func Error(t *testing.T, err error) {
	if err == nil {
		failTest(t, "expected an error, but got nil")
	}
}

// True asserts that a condition is true.
// It fails the test if the condition is false.
func True(t *testing.T, condition bool) {
	if !condition {
		failTest(t, "expected true, but got false")
	}
}

// False asserts that a condition is false.
// It fails the test if the condition is true.
func False(t *testing.T, condition bool) {
	if condition {
		failTest(t, "expected false, but got true")
	}
}

// Contains asserts that a container includes a specific element.
// Supported container types are strings, slices, arrays, and maps.
func Contains(t *testing.T, container, item any) {
	var exists bool
	switch c := container.(type) {
	case string:
		s, ok := item.(string)
		if !ok {
			failTest(t, fmt.Sprintf("item must be a string when container is a string, got %T", item))
			return
		}
		exists = strings.Contains(c, s)
	default:
		v := reflect.ValueOf(container)
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				if reflect.DeepEqual(v.Index(i).Interface(), item) {
					exists = true
					break
				}
			}
		case reflect.Map:
			exists = v.MapIndex(reflect.ValueOf(item)).IsValid()
		default:
			failTest(t, fmt.Sprintf("unsupported container type: %T", container))
			return
		}
	}

	if !exists {
		failTest(t, fmt.Sprintf("expected %v to contain %v, but it did not", container, item))
	}
}

// NotContains asserts that a container does not include a specific element.
// Supported container types are strings, slices, arrays, and maps.
func NotContains(t *testing.T, container, item any) {
	var exists bool
	switch c := container.(type) {
	case string:
		s, ok := item.(string)
		if !ok {
			failTest(t, fmt.Sprintf("item must be a string when container is a string, got %T", item))
			return
		}
		exists = strings.Contains(c, s)
	default:
		v := reflect.ValueOf(container)
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				if reflect.DeepEqual(v.Index(i).Interface(), item) {
					exists = true
					break
				}
			}
		case reflect.Map:
			exists = v.MapIndex(reflect.ValueOf(item)).IsValid()
		default:
			failTest(t, fmt.Sprintf("unsupported container type: %T", container))
			return
		}
	}

	if exists {
		failTest(t, fmt.Sprintf("expected %v to not contain %v, but it did", container, item))
	}
}

// Len asserts that an object has a specific length.
// Supported types are arrays, slices, maps, and strings.
func Len(t *testing.T, object any, length int) {
	objectValue := reflect.ValueOf(object)
	switch objectValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		if objectValue.Len() != length {
			failTest(t, fmt.Sprintf("expected length %d, but got %d", length, objectValue.Len()))
		}
	default:
		failTest(t, fmt.Sprintf("unsupported type for length check: %T", object))
	}
}

// Panics asserts that a function panics when called.
func Panics(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r == nil {
			failTest(t, "expected panic, but none occurred")
		}
	}()
	fn()
}

// NotPanics asserts that a function does not panic when called.
func NotPanics(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			failTest(t, fmt.Sprintf("unexpected panic: %v", r))
		}
	}()
	fn()
}

// Same asserts that two pointers reference the same object.
func Same(t *testing.T, expected, actual any) {
	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)

	// Check if both expected and actual are pointers
	if expectedVal.Kind() != reflect.Ptr || actualVal.Kind() != reflect.Ptr {
		failTest(t, fmt.Sprintf("expected and actual must both be pointers, but got: %T vs %T", expected, actual))
		return
	}

	// Compare the pointers' addresses
	if expectedVal.Pointer() != actualVal.Pointer() {
		failTest(t, fmt.Sprintf("expected same address, but got different: %p vs %p", expected, actual))
	}
}

// Greater asserts that the first value is greater than the second.
func Greater(t *testing.T, a, b any) {
	cmp, err := compareNumeric(a, b)
	if err != nil {
		failTest(t, fmt.Sprintf("failed to compare values: %v", err))
		return
	}
	if cmp <= 0 {
		failTest(t, fmt.Sprintf("expected %v to be greater than %v", a, b))
	}
}

// Less asserts that the first value is less than the second.
func Less(t *testing.T, a, b any) {
	cmp, err := compareNumeric(a, b)
	if err != nil {
		failTest(t, fmt.Sprintf("failed to compare values: %v", err))
		return
	}
	if cmp >= 0 {
		failTest(t, fmt.Sprintf("expected %v to be less than %v", a, b))
	}
}

// IsOfType asserts that an object is of a specific type.
func IsOfType(t *testing.T, expectedType, obj any) {
	if reflect.TypeOf(obj) != reflect.TypeOf(expectedType) {
		failTest(t, fmt.Sprintf("expected type %T, but got %T", expectedType, obj))
	}
}

// LessOrEqual asserts that the first value is less than or equal to the second.
func LessOrEqual(t *testing.T, a, b any) {
	cmp, err := compareNumeric(a, b)
	if err != nil {
		failTest(t, fmt.Sprintf("failed to compare values: %v", err))
		return
	}
	if cmp > 0 {
		failTest(t, fmt.Sprintf("expected %v to be less than or equal to %v", a, b))
	}
}

// GreaterOrEqual asserts that the first value is greater than or equal to the second.
func GreaterOrEqual(t *testing.T, a, b any) {
	cmp, err := compareNumeric(a, b)
	if err != nil {
		failTest(t, fmt.Sprintf("failed to compare values: %v", err))
		return
	}
	if cmp < 0 {
		failTest(t, fmt.Sprintf("expected %v to be greater than or equal to %v", a, b))
	}
}

// IsZero asserts that the value is the zero value for its type.
func IsZero(t *testing.T, value any) {
	if !reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface()) {
		failTest(t, fmt.Sprintf("expected zero value, but got: %v", value))
	}
}

// Subset asserts that a slice, array, or map contains all elements of another.
func Subset(t *testing.T, list, subset any) {
	listVal := reflect.ValueOf(list)
	subsetVal := reflect.ValueOf(subset)

	switch listVal.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < subsetVal.Len(); i++ {
			if !contains(listVal, subsetVal.Index(i).Interface()) {
				failTest(t, fmt.Sprintf("expected %v to be a subset of %v, but it's not", subset, list))
				return
			}
		}
	case reflect.Map:
		for _, key := range subsetVal.MapKeys() {
			if !listVal.MapIndex(key).IsValid() || !reflect.DeepEqual(listVal.MapIndex(key).Interface(), subsetVal.MapIndex(key).Interface()) {
				failTest(t, fmt.Sprintf("expected %v to be a subset of %v, but it's not", subset, list))
				return
			}
		}
	default:
		failTest(t, fmt.Sprintf("unsupported type for Subset: %T", list))
	}
}

// contains is a helper function to check if a value is in a slice or array.
func contains(listVal reflect.Value, item interface{}) bool {
	for i := 0; i < listVal.Len(); i++ {
		if reflect.DeepEqual(listVal.Index(i).Interface(), item) {
			return true
		}
	}
	return false
}

// ErrorContains asserts that the error message contains a specific substring.
func ErrorContains(t *testing.T, err error, substr string) {
	if err == nil {
		failTest(t, "expected an error, but got nil")
		return
	}
	if !strings.Contains(err.Error(), substr) {
		failTest(t, fmt.Sprintf("expected error message to contain %q, but got %q", substr, err.Error()))
	}
}

// Implements asserts that an object implements a specific interface type.
// The interfaceType argument must be a pointer to an interface.
func Implements(t *testing.T, interfaceType, obj any) {
	objType := reflect.TypeOf(obj)
	if !objType.Implements(reflect.TypeOf(interfaceType).Elem()) {
		failTest(t, fmt.Sprintf("expected %T to implement %T, but it does not", obj, interfaceType))
	}
}

// SameElements asserts that two slices or arrays contain the same elements, regardless of order.
func SameElements(t *testing.T, a, b any) {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != reflect.Slice && aVal.Kind() != reflect.Array {
		failTest(t, "first argument must be a slice or array")
		return
	}
	if bVal.Kind() != reflect.Slice && bVal.Kind() != reflect.Array {
		failTest(t, "second argument must be a slice or array")
		return
	}

	if aVal.Len() != bVal.Len() {
		failTest(t, fmt.Sprintf("expected slices of the same length, but got %d and %d", aVal.Len(), bVal.Len()))
		return
	}

	aMap := make(map[interface{}]int)
	bMap := make(map[interface{}]int)

	// Ensure only hashable (comparable) types are used as keys
	for i := 0; i < aVal.Len(); i++ {
		aElem := aVal.Index(i).Interface()
		bElem := bVal.Index(i).Interface()

		if !isHashable(reflect.ValueOf(aElem).Kind()) || !isHashable(reflect.ValueOf(bElem).Kind()) {
			failTest(t, "unsupported element type for comparison")
			return
		}

		aMap[aElem]++
		bMap[bElem]++
	}

	for key, countA := range aMap {
		if countB, ok := bMap[key]; !ok || countA != countB {
			failTest(t, fmt.Sprintf("expected same elements in both slices, but %v differs", key))
			return
		}
	}
}

// Helper function to check if a type is hashable
func isHashable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128, reflect.String, reflect.Chan, reflect.Func, reflect.Ptr:
		return true
	default:
		return false
	}
}

// MatchesRegex asserts that a string matches a regular expression.
func MatchesRegex(t *testing.T, str, pattern string) {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		failTest(t, fmt.Sprintf("invalid regex pattern: %v", err))
		return
	}
	if !matched {
		failTest(t, fmt.Sprintf("expected string %q to match regex %q, but it did not", str, pattern))
	}
}

// HasPrefix asserts that a string has a specific prefix.
func HasPrefix(t *testing.T, str, prefix string) {
	if !strings.HasPrefix(str, prefix) {
		failTest(t, fmt.Sprintf("expected string %q to have prefix %q, but it did not", str, prefix))
	}
}

// HasSuffix asserts that a string has a specific suffix.
func HasSuffix(t *testing.T, str, suffix string) {
	if !strings.HasSuffix(str, suffix) {
		failTest(t, fmt.Sprintf("expected string %q to have suffix %q, but it did not", str, suffix))
	}
}

// WithinDuration asserts that two time.Time values are within a certain duration of each other.
func WithinDuration(t *testing.T, expected, actual time.Time, delta time.Duration) {
	diff := expected.Sub(actual)
	if diff < -delta || diff > delta {
		failTest(t, fmt.Sprintf("expected time %v to be within %v of %v, but difference was %v", actual, delta, expected, diff))
	}
}

// JSONEq asserts that two JSON strings are equivalent, ignoring differences in whitespace or key ordering.
func JSONEq(t *testing.T, expected, actual string) {
	var expectedJSON, actualJSON interface{}
	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		failTest(t, fmt.Sprintf("failed to unmarshal expected JSON: %v", err))
		return
	}
	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		failTest(t, fmt.Sprintf("failed to unmarshal actual JSON: %v", err))
		return
	}
	if !reflect.DeepEqual(expectedJSON, actualJSON) {
		failTest(t, fmt.Sprintf("JSON not equal: expected: %v actual: %v", expectedJSON, actualJSON))
	}
}

// PanicsWithValue asserts that a function panics with a specific value.
func PanicsWithValue(t *testing.T, expected any, fn func()) {
	defer func() {
		if r := recover(); r == nil {
			failTest(t, "expected panic, but none occurred")
		} else if !reflect.DeepEqual(r, expected) {
			failTest(t, fmt.Sprintf("expected panic value %v, but got %v", expected, r))
		}
	}()
	fn()
}

// InDelta asserts that two numeric values are within delta of each other.
func InDelta(t *testing.T, expected, actual any, delta float64) {
	a, err := toFloat64(expected)
	if err != nil {
		failTest(t, fmt.Sprintf("expected value is not numeric: %v", err))
		return
	}
	b, err := toFloat64(actual)
	if err != nil {
		failTest(t, fmt.Sprintf("actual value is not numeric: %v", err))
		return
	}
	if diff := math.Abs(a - b); diff > delta {
		failTest(t, fmt.Sprintf("expected %v to be within %v of %v, but difference was %v", actual, delta, expected, diff))
	}
}

// InEpsilon asserts that two numeric values are within epsilon percent of each other.
func InEpsilon(t *testing.T, expected, actual any, epsilon float64) {
	a, err := toFloat64(expected)
	if err != nil {
		failTest(t, fmt.Sprintf("expected value is not numeric: %v", err))
		return
	}
	b, err := toFloat64(actual)
	if err != nil {
		failTest(t, fmt.Sprintf("actual value is not numeric: %v", err))
		return
	}
	if a == b {
		return
	}
	diff := math.Abs(a - b)
	mean := math.Abs(a+b) / 2
	if diff/mean > epsilon {
		failTest(t, fmt.Sprintf("expected %v to be within %v%% of %v, but difference was %v%%", actual, epsilon*100, expected, diff/mean*100))
	}
}

// ElementsMatch asserts that two slices or arrays have the same elements in any order.
// Duplicate elements are checked for and must appear the same number of times in both slices.
func ElementsMatch(t *testing.T, listA, listB any) {
	if !haveSameElements(listA, listB) {
		failTest(t, fmt.Sprintf("element lists are not equal: expected: %v actual: %v", listA, listB))
	}
}

// haveSameElements is a helper function for ElementsMatch.
func haveSameElements(listA, listB any) bool {
	valA := reflect.ValueOf(listA)
	valB := reflect.ValueOf(listB)

	aLen := valA.Len()
	bLen := valB.Len()

	if aLen != bLen {
		return false
	}

	// Use maps to count element occurrences
	countA := make(map[interface{}]int)
	countB := make(map[interface{}]int)

	for i := 0; i < aLen; i++ {
		countA[valA.Index(i).Interface()]++
	}
	for i := 0; i < bLen; i++ {
		countB[valB.Index(i).Interface()]++
	}

	// Compare element counts in both maps
	return reflect.DeepEqual(countA, countB)
}
