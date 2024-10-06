package asserts

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func mockTestingEnable() {
	mockTesting = true
	mockTestMessage = ""
}

func mockTestMessageCheck(t *testing.T, expected string) {
	if !strings.Contains(mockTestMessage, expected) {
		t.Errorf("Expected message %q, got %q", expected, mockTestMessage)
	}
	mockTestMessage = ""
}

func TestIsNil(t *testing.T) {
	if !isNil(nil) {
		t.Error("Expected true, got false")
	}
	if isNil(5) {
		t.Error("Expected false, got true")
	}
	var a *int
	if !isNil(a) {
		t.Error("Expected true, got false")
	}
}

func TestCompareNumeric(t *testing.T) {
	t.Run("Test Equal", func(t *testing.T) {
		num, noErr := compareNumeric(5, 5)
		if noErr != nil {
			t.Errorf("Unexpected error: %v", noErr)
		}
		if num != 0 {
			t.Errorf("Expected 0, got %d", num)
		}
	})
	t.Run("Test GreaterThan", func(t *testing.T) {
		num, noErr := compareNumeric(8, 5)
		if noErr != nil {
			t.Errorf("Unexpected error: %v", noErr)
		}
		if num != 1 {
			t.Errorf("Expected 1, got %d", num)
		}
	})
	t.Run("Test LowerThan", func(t *testing.T) {
		num, noErr := compareNumeric(5, 8)
		if noErr != nil {
			t.Errorf("Unexpected error: %v", noErr)
		}
		if num != -1 {
			t.Errorf("Expected -1, got %d", num)
		}
	})
	t.Run("Test Error", func(t *testing.T) {
		_, err := compareNumeric(5, "1")
		if err == nil {
			t.Error("Expected error, got none")
		}
	})
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		input    any
		expected float64
		hasError bool
	}{
		{input: int(5), expected: 5.0, hasError: false},
		{input: int8(5), expected: 5.0, hasError: false},
		{input: int16(5), expected: 5.0, hasError: false},
		{input: int32(5), expected: 5.0, hasError: false},
		{input: int64(5), expected: 5.0, hasError: false},
		{input: uint(5), expected: 5.0, hasError: false},
		{input: uint8(5), expected: 5.0, hasError: false},
		{input: uint16(5), expected: 5.0, hasError: false},
		{input: uint32(5), expected: 5.0, hasError: false},
		{input: uint64(5), expected: 5.0, hasError: false},
		{input: float32(5.5), expected: 5.5, hasError: false},
		{input: float64(5.5), expected: 5.5, hasError: false},
		// Unsupported type should return an error
		{input: "5", expected: 0, hasError: true},
		{input: struct{}{}, expected: 0, hasError: true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.input), func(t *testing.T) {
			result, err := toFloat64(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for input: %v, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input: %v, error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("expected: %v, but got: %v", tt.expected, result)
				}
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
	}{
		{input: nil, expected: true},
		{input: 0, expected: true},
		{input: "", expected: true},
		{input: []int{}, expected: true},
		{input: map[string]int{}, expected: true},
		{input: false, expected: true},
		{input: 0.0, expected: true},
		{input: uint(0), expected: true},
		{input: struct{}{}, expected: true},
		{input: (*int)(nil), expected: true},
		{input: 1, expected: false},
		{input: "hello", expected: false},
		{input: []int{1, 2}, expected: false},
		{input: map[string]int{"a": 1}, expected: false},
		{input: true, expected: false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.input), func(t *testing.T) {
			result := isEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("expected: %v, but got: %v", tt.expected, result)
			}
		})
	}
	b := 1
	if isEmpty(&b) {
		t.Error("Expected false, got true")
	}
}

func TestHaveSameElements(t *testing.T) {
	t.Run("Same elements but in different order", func(t *testing.T) {
		if !haveSameElements([]int{1, 2, 3}, []int{3, 2, 1}) {
			t.Error("Expected true, got false for same elements in different order")
		}
	})

	t.Run("Different elements", func(t *testing.T) {
		if haveSameElements([]int{1, 2, 3}, []int{1, 2, 4}) {
			t.Error("Expected false, got true for different elements")
		}
	})

	t.Run("Same elements with different duplicates", func(t *testing.T) {
		if haveSameElements([]int{1, 2, 2}, []int{1, 2, 3}) {
			t.Error("Expected false, got true for same elements but different duplicates")
		}
	})

	t.Run("Different length lists", func(t *testing.T) {
		if haveSameElements([]int{1, 2, 3}, []int{1, 2}) {
			t.Error("Expected false, got true for lists with different lengths")
		}
	})

	t.Run("Empty lists", func(t *testing.T) {
		if !haveSameElements([]int{}, []int{}) {
			t.Error("Expected true, got false for two empty lists")
		}
	})

	t.Run("Same struct elements", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		listA := []person{{"Alice", 30}, {"Bob", 25}}
		listB := []person{{"Bob", 25}, {"Alice", 30}}

		if !haveSameElements(listA, listB) {
			t.Error("Expected true, got false for same struct elements in different order")
		}
	})
}

func TestEqual(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		Equal(t, 5, 5)
		Equal(t, "hello", "hello")
		Equal(t, []int{1, 2}, []int{1, 2})
		Equal(t, map[string]int{"a": 1}, map[string]int{"a": 1})
		Equal(t, nil, nil)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		Equal(t, 5, "5")
		mockTestMessageCheck(t, "values not equal: expected: 5 actual: 5")
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		NotEqual(t, 5, 6)
		NotEqual(t, "hello", "world")
		NotEqual(t, []int{1, 2}, []int{2, 3})
		NotEqual(t, map[string]int{"a": 1}, map[string]int{"b": 2})
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		NotEqual(t, 5, 5)
		mockTestMessageCheck(t, "values unexpectedly equal: not expected: 5 actual: 5")
	})
}

func TestNil(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		var ptr *int = nil
		Nil(t, ptr)
		Nil(t, nil)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		var i int = 5
		Nil(t, &i)
		mockTestMessageCheck(t, "expected nil, but got:")
	})
}

func TestNotNil(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		var i int = 5
		NotNil(t, &i)
		NotNil(t, i)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		var ptr *int = nil
		NotNil(t, ptr)
		mockTestMessageCheck(t, "expected non-nil value, but got nil")
	})
}

func TestNotEmpty(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		NotEmpty(t, 1)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		var a int = 0
		NotEmpty(t, a)
		mockTestMessageCheck(t, "expected non-empty value, but got empty")
	})
}

func TestEmpty(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		var a int
		Empty(t, a)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		var a int = 1
		Empty(t, a)
		mockTestMessageCheck(t, "expected empty value, but got: 1")
	})
}

func TestNoError(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		var err error = nil
		NoError(t, err)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		var err error = fmt.Errorf("an error occurred")
		NoError(t, err)
		mockTestMessageCheck(t, "unexpected error: an error occurred")
	})
}

func TestError(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		var err error = fmt.Errorf("an error occurred")
		Error(t, err)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		var err error = nil
		Error(t, err)
		mockTestMessageCheck(t, "expected an error, but got nil")
	})
}

func TestTrue(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		True(t, true)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		True(t, false)
		mockTestMessageCheck(t, "expected true, but got false")
	})
}

func TestFalse(t *testing.T) {
	t.Run("Truthy", func(t *testing.T) {
		False(t, false)
	})

	t.Run("Falsy", func(t *testing.T) {
		mockTestingEnable()
		False(t, true)
		mockTestMessageCheck(t, "expected false, but got true")
	})
}

func TestContains(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		Contains(t, "hello world", "world")
	})

	t.Run("Slice", func(t *testing.T) {
		Contains(t, []int{1, 2, 3}, 2)
	})

	t.Run("Map", func(t *testing.T) {
		Contains(t, map[string]int{"a": 1, "b": 2}, "a")
	})

	t.Run("Not Contains different type", func(t *testing.T) {
		mockTestingEnable()
		Contains(t, "hello world", 1)
		mockTestMessageCheck(t, "item must be a string when container is a string")
	})

	t.Run("Not Contains different type", func(t *testing.T) {
		mockTestingEnable()
		Contains(t, struct{}{}, 1)
		mockTestMessageCheck(t, "unsupported container type: struct {}")
	})

	t.Run("Not Contains String", func(t *testing.T) {
		mockTestingEnable()
		Contains(t, "hello world", "universe")
		mockTestMessageCheck(t, "expected hello world to contain universe, but it did not")
	})

	t.Run("Not Contains Slice", func(t *testing.T) {
		mockTestingEnable()
		Contains(t, []int{1, 2, 3}, 4)
		mockTestMessageCheck(t, "expected [1 2 3] to contain 4, but it did not")
	})

	t.Run("Not Contains Map", func(t *testing.T) {
		mockTestingEnable()
		Contains(t, map[string]int{"a": 1, "b": 2}, "c")
		mockTestMessageCheck(t, "expected map[a:1 b:2] to contain c, but it did not")
	})
}

func TestNotContains(t *testing.T) {
	t.Run("Not Contains String", func(t *testing.T) {
		NotContains(t, "hello world", "universe")
	})

	t.Run("Not Contains Slice", func(t *testing.T) {
		NotContains(t, []int{1, 2, 3}, 4)
	})

	t.Run("Not Contains Map", func(t *testing.T) {
		NotContains(t, map[string]int{"a": 1, "b": 2}, "c")
	})

	t.Run("Not Contains different type", func(t *testing.T) {
		mockTestingEnable()
		NotContains(t, "hello world", 1)
		mockTestMessageCheck(t, "item must be a string when container is a string")
	})

	t.Run("Not Contains different type", func(t *testing.T) {
		mockTestingEnable()
		NotContains(t, struct{}{}, 1)
		mockTestMessageCheck(t, "unsupported container type: struct {}")
	})

	t.Run("Contains String", func(t *testing.T) {
		mockTestingEnable()
		NotContains(t, "hello world", "world")
		mockTestMessageCheck(t, "expected hello world to not contain world, but it did")
	})

	t.Run("Contains Slice", func(t *testing.T) {
		mockTestingEnable()
		NotContains(t, []int{1, 2, 3}, 2)
		mockTestMessageCheck(t, "expected [1 2 3] to not contain 2, but it did")
	})

	t.Run("Contains Map", func(t *testing.T) {
		mockTestingEnable()
		NotContains(t, map[string]int{"a": 1, "b": 2}, "a")
		mockTestMessageCheck(t, "expected map[a:1 b:2] to not contain a, but it did")
	})
}

func TestLen(t *testing.T) {
	t.Run("Correct", func(t *testing.T) {
		Len(t, "hello", 5)
		Len(t, []int{1, 2, 3}, 3)
		Len(t, map[string]int{"a": 1, "b": 2}, 2)
	})

	t.Run("Incorrect", func(t *testing.T) {
		mockTestingEnable()
		Len(t, "hello", 4)
		mockTestMessageCheck(t, "expected length 4, but got 5")
	})

	t.Run("Different kind", func(t *testing.T) {
		mockTestingEnable()
		Len(t, struct{}{}, 4)
		mockTestMessageCheck(t, "unsupported type for length check: struct {}")
	})
}

func TestPanics(t *testing.T) {
	t.Run("Panics", func(t *testing.T) {
		Panics(t, func() { panic("panic") })
	})

	t.Run("Does Not Panic", func(t *testing.T) {
		mockTestingEnable()
		Panics(t, func() {})
		mockTestMessageCheck(t, "expected panic, but none occurred")
	})
}

func TestNotPanics(t *testing.T) {
	t.Run("Not Panics", func(t *testing.T) {
		NotPanics(t, func() {})
	})

	t.Run("Panics", func(t *testing.T) {
		mockTestingEnable()
		NotPanics(t, func() { panic("panic") })
		mockTestMessageCheck(t, "unexpected panic: panic")
	})
}

func TestSame(t *testing.T) {
	t.Run("Same Pointers", func(t *testing.T) {
		var a = &struct{}{}
		Same(t, a, a)
	})

	t.Run("Different Pointers", func(t *testing.T) {
		mockTestingEnable()
		var a = &struct{ x int }{x: 1}
		var b = &struct{ x int }{x: 2}
		Same(t, a, b)
		mockTestMessageCheck(t, "expected same address, but got different:")
	})

	t.Run("Non-Pointers", func(t *testing.T) {
		mockTestingEnable()
		Same(t, 5, 5)
		mockTestMessageCheck(t, "expected and actual must both be pointers, but got: int vs int")
	})

	t.Run("Different kinds", func(t *testing.T) {
		mockTestingEnable()
		Same(t, struct{}{}, 4)
		mockTestMessageCheck(t, "expected and actual must both be pointers, but got: struct {} vs int")
	})
}

func TestGreater(t *testing.T) {
	t.Run("Greater", func(t *testing.T) {
		Greater(t, 5, 3)
		Greater(t, 5.1, 5.0)
		Greater(t, uint(5), uint(3))
	})

	t.Run("Not Greater", func(t *testing.T) {
		mockTestingEnable()
		Greater(t, 3, 5)
		mockTestMessageCheck(t, "expected 3 to be greater than 5")
	})

	t.Run("Unsupported Types", func(t *testing.T) {
		mockTestingEnable()
		Greater(t, "5", "3")
		mockTestMessageCheck(t, "failed to compare values: unsupported numeric types: string vs string")
	})
}

func TestLess(t *testing.T) {
	t.Run("Less", func(t *testing.T) {
		Less(t, 3, 5)
		Less(t, 5.0, 5.1)
		Less(t, uint(3), uint(5))
	})

	t.Run("Not Less", func(t *testing.T) {
		mockTestingEnable()
		Less(t, 5, 3)
		mockTestMessageCheck(t, "expected 5 to be less than 3")
	})

	t.Run("Unsupported Types", func(t *testing.T) {
		mockTestingEnable()
		Less(t, "3", "5")
		mockTestMessageCheck(t, "failed to compare values: unsupported numeric types: string vs string")
	})
}

func TestIsOfType(t *testing.T) {
	t.Run("IsType", func(t *testing.T) {
		IsOfType(t, 5, 10)
		IsOfType(t, "hello", "world")
		IsOfType(t, []int{}, []int{1, 2, 3})
	})

	t.Run("IsType Fail", func(t *testing.T) {
		mockTestingEnable()
		IsOfType(t, 5, "5")
		mockTestMessageCheck(t, "expected type int, but got string")
	})
}

func TestLessOrEqual(t *testing.T) {
	t.Run("LessOrEqual", func(t *testing.T) {
		LessOrEqual(t, 5, 5)
		LessOrEqual(t, 5, 6)
		LessOrEqual(t, 5.0, 5.1)
		LessOrEqual(t, uint(5), uint(5))
	})

	t.Run("Not LessOrEqual", func(t *testing.T) {
		mockTestingEnable()
		LessOrEqual(t, 5, 3)
		mockTestMessageCheck(t, "expected 5 to be less than or equal to 3")
	})

	t.Run("Unsupported Types", func(t *testing.T) {
		mockTestingEnable()
		LessOrEqual(t, "5", "3")
		mockTestMessageCheck(t, "failed to compare values: unsupported numeric types: string vs string")
	})
}

func TestGreaterOrEqual(t *testing.T) {
	t.Run("GreaterOrEqual", func(t *testing.T) {
		GreaterOrEqual(t, 5, 5)
		GreaterOrEqual(t, 6, 5)
		GreaterOrEqual(t, 5.1, 5.0)
		GreaterOrEqual(t, uint(5), uint(5))
	})

	t.Run("Not GreaterOrEqual", func(t *testing.T) {
		mockTestingEnable()
		GreaterOrEqual(t, 3, 5)
		mockTestMessageCheck(t, "expected 3 to be greater than or equal to 5")
	})

	t.Run("Unsupported Types", func(t *testing.T) {
		mockTestingEnable()
		GreaterOrEqual(t, "5", "3")
		mockTestMessageCheck(t, "failed to compare values: unsupported numeric types: string vs string")
	})
}

func TestIsZero(t *testing.T) {
	t.Run("IsZero", func(t *testing.T) {
		IsZero(t, 0)
		IsZero(t, 0.0)
		IsZero(t, uint(0))
		IsZero(t, false)
	})

	t.Run("Not IsZero", func(t *testing.T) {
		mockTestingEnable()
		IsZero(t, 5)
		mockTestMessageCheck(t, "expected zero value, but got: 5")
	})
}

func TestSubset(t *testing.T) {
	t.Run("Subset", func(t *testing.T) {
		Subset(t, []int{1, 2, 3}, []int{1, 2})
		Subset(t, map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1})
	})

	t.Run("Subset Map", func(t *testing.T) {
		Subset(t, map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2})
	})

	t.Run("Subset Map Falsy", func(t *testing.T) {
		mockTestingEnable()
		Subset(t, map[string]int{"a": 1, "b": 2}, map[string]int{"c": 3})
		mockTestMessageCheck(t, "expected map[c:3] to be a subset of map[a:1 b:2], but it's not")
	})

	t.Run("Not Subset", func(t *testing.T) {
		mockTestingEnable()
		Subset(t, []int{1, 2, 3}, []int{4, 5})
		mockTestMessageCheck(t, "expected [4 5] to be a subset of [1 2 3], but it's not")
	})

	t.Run("Unsupported Types", func(t *testing.T) {
		mockTestingEnable()
		Subset(t, "hello", "world")
		mockTestMessageCheck(t, "unsupported type for Subset: string")
	})
}

func TestErrorContains(t *testing.T) {
	t.Run("ErrorContains", func(t *testing.T) {
		ErrorContains(t, fmt.Errorf("an error occurred"), "error")
	})

	t.Run("ErrorContains nil", func(t *testing.T) {
		mockTestingEnable()
		ErrorContains(t, nil, "world")
		mockTestMessageCheck(t, "expected an error, but got nil")
	})

	t.Run("ErrorContains Fail", func(t *testing.T) {
		mockTestingEnable()
		ErrorContains(t, fmt.Errorf("an error occurred"), "world")
		mockTestMessageCheck(t, "expected error message to contain \"world\", but got \"an error occurred\"")
	})
}

func TestImplements(t *testing.T) {
	t.Run("Implements", func(t *testing.T) {
		Implements(t, (*testing.TB)(nil), t)
	})

	t.Run("Not Implements", func(t *testing.T) {
		mockTestingEnable()
		Implements(t, (*testing.TB)(nil), 5)
		mockTestMessageCheck(t, "expected int to implement *testing.TB, but it does not")
	})
}

func TestSameElements(t *testing.T) {
	t.Run("SameElements", func(t *testing.T) {
		SameElements(t, []int{1, 2, 3}, []int{3, 2, 1})
	})

	t.Run("Not SameElements 1", func(t *testing.T) {
		mockTestingEnable()
		SameElements(t, []int{1, 2, 3}, []int{4, 5})
		mockTestMessageCheck(t, "expected slices of the same length, but got 3 and 2")
	})

	t.Run("Not SameElements 2", func(t *testing.T) {
		mockTestingEnable()
		SameElements(t, []int{1, 2, 3}, []int{4, 5, 6})
		mockTestMessageCheck(t, "expected same elements in both slices, but")
	})

	t.Run("Not SameElements with non-hashable slices", func(t *testing.T) {
		mockTestingEnable()
		SameElements(t, [][]int{{1, 2}, {3, 4}}, [][]int{{3, 4}, {1, 2}})
		mockTestMessageCheck(t, "unsupported element type for comparison")
	})

	t.Run("Different Types 1", func(t *testing.T) {
		mockTestingEnable()
		SameElements(t, 1, []string{"1", "2", "3"})
		mockTestMessageCheck(t, "first argument must be a slice or array")
	})

	t.Run("Different Types 2", func(t *testing.T) {
		mockTestingEnable()
		SameElements(t, []string{"1", "2", "3"}, 2)
		mockTestMessageCheck(t, "second argument must be a slice or array")
	})
}

func TestMatchesRegex(t *testing.T) {
	t.Run("MatchesRegex", func(t *testing.T) {
		MatchesRegex(t, "hello", "hello")
	})

	t.Run("Invalid MatchesRegex", func(t *testing.T) {
		mockTestingEnable()
		MatchesRegex(t, "hello", "(")
		mockTestMessageCheck(t, "invalid regex pattern: error parsing regexp: missing closing ): `(`")
	})

	t.Run("Not MatchesRegex", func(t *testing.T) {
		mockTestingEnable()
		MatchesRegex(t, "hello", "world")
		mockTestMessageCheck(t, "expected string \"hello\" to match regex \"world\", but it did not")
	})
}

func TestHasSuffix(t *testing.T) {
	t.Run("HasSuffix", func(t *testing.T) {
		HasSuffix(t, "hello", "lo")
	})

	t.Run("Not HasSuffix", func(t *testing.T) {
		mockTestingEnable()
		HasSuffix(t, "hello", "world")
		mockTestMessageCheck(t, "expected string \"hello\" to have suffix \"world\", but it did not")
	})
}

func TestHasPrefix(t *testing.T) {
	t.Run("HasPrefix", func(t *testing.T) {
		HasPrefix(t, "hello", "he")
	})

	t.Run("Not HasPrefix", func(t *testing.T) {
		mockTestingEnable()
		HasPrefix(t, "hello", "world")
		mockTestMessageCheck(t, "expected string \"hello\" to have prefix \"world\", but it did not")
	})
}

func TestWithinDuration(t *testing.T) {
	t.Run("WithinDuration", func(t *testing.T) {
		time := time.Now()
		WithinDuration(t, time, time, 0)
	})

	t.Run("Not WithinDuration", func(t *testing.T) {
		mockTestingEnable()
		time1 := time.Now()
		time2 := time1.Add(30 * time.Second)
		WithinDuration(t, time1, time2, 0)
		mockTestMessageCheck(t, ", but difference was -30s")
	})
}

func TestJSONEq(t *testing.T) {
	t.Run("JSONEq", func(t *testing.T) {
		JSONEq(t, `{"hello": "world"}`, `{"hello": "world"}`)
	})

	t.Run("Not JSONEq NotJson1", func(t *testing.T) {
		mockTestingEnable()
		JSONEq(t, `{"hello" / }`, `{"hello": "universe"}`)
		mockTestMessageCheck(t, "failed to unmarshal expected JSON: invalid character")
	})

	t.Run("Not JSONEq NotJson2", func(t *testing.T) {
		mockTestingEnable()
		JSONEq(t, `{"hello": "universe"}`, `{"hello" / }`)
		mockTestMessageCheck(t, "failed to unmarshal actual JSON: invalid character")
	})

	t.Run("Not JSONEq", func(t *testing.T) {
		mockTestingEnable()
		JSONEq(t, `{"hello": "world"}`, `{"hello": "universe"}`)
		mockTestMessageCheck(t, "JSON not equal: expected: map[hello:world] actual: map[hello:universe]")
	})
}

func TestPanicsWithValue(t *testing.T) {
	t.Run("PanicsWithValue", func(t *testing.T) {
		PanicsWithValue(t, "panic", func() { panic("panic") })
	})

	t.Run("Not PanicsWithValue", func(t *testing.T) {
		mockTestingEnable()
		PanicsWithValue(t, "panic", func() {})
		mockTestMessageCheck(t, "expected panic, but none occurred")
	})

	t.Run("Not PanicsWithValue Value", func(t *testing.T) {
		mockTestingEnable()
		PanicsWithValue(t, "special", func() { panic("normal") })
		mockTestMessageCheck(t, "expected panic value special, but got normal")
	})
}

func TestInDelta(t *testing.T) {
	t.Run("InDelta", func(t *testing.T) {
		InDelta(t, 5.0, 5.1, 0.2)
	})

	t.Run("InDelta DataType 1", func(t *testing.T) {
		mockTestingEnable()
		InDelta(t, "a", 5.1, 0.05)
		mockTestMessageCheck(t, "expected value is not numeric: unsupported type for numeric comparison: string")
	})

	t.Run("InDelta DataType 2", func(t *testing.T) {
		mockTestingEnable()
		InDelta(t, 5.1, "a", 0.05)
		mockTestMessageCheck(t, "actual value is not numeric: unsupported type for numeric comparison: string")
	})

	t.Run("Not InDelta", func(t *testing.T) {
		mockTestingEnable()
		InDelta(t, 5.0, 5.1, 0.05)
		mockTestMessageCheck(t, "expected 5.1 to be within 0.05 of 5, but difference was 0.099")
	})
}

func TestInEpsilon(t *testing.T) {
	t.Run("InEpsilon", func(t *testing.T) {
		InEpsilon(t, 5.0, 5.0, 0.2)
		InEpsilon(t, 5.0, 5.1, 0.2)
	})

	t.Run("InEpsilon DataType 1", func(t *testing.T) {
		mockTestingEnable()
		InEpsilon(t, "a", 5.1, 0.05)
		mockTestMessageCheck(t, "expected value is not numeric: unsupported type for numeric comparison: string")
	})

	t.Run("InEpsilon DataType 2", func(t *testing.T) {
		mockTestingEnable()
		InEpsilon(t, 5.1, "a", 0.05)
		mockTestMessageCheck(t, "actual value is not numeric: unsupported type for numeric comparison: string")
	})

	t.Run("Not InEpsilon", func(t *testing.T) {
		mockTestingEnable()
		InEpsilon(t, 5.0, 5.1, 0.0006)
		mockTestMessageCheck(t, "expected 5.1 to be within 0.06% of 5, but difference was 1.98")
	})
}

func TestElementsMatch(t *testing.T) {
	t.Run("ElementsMatch", func(t *testing.T) {
		ElementsMatch(t, []int{1, 2, 3}, []int{3, 2, 1})
	})

	t.Run("ElementsMatch Fail", func(t *testing.T) {
		mockTestingEnable()
		ElementsMatch(t, []int{1, 2, 3}, []int{4, 5})
		mockTestMessageCheck(t, "element lists are not equal: expected: [1 2 3] actual: [4 5]")
	})
}
