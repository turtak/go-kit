package stacktrace

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewStackTrace(t *testing.T) {
	config := Config{
		BufferSize: 2048,
		SkipFrames: 0,
	}

	st := NewStackTrace(&config)

	if st == nil {
		t.Fatal("NewStackTrace returned nil")
	}

	if len(st.frames) == 0 {
		t.Error("StackTrace has no frames")
	}

	if st.raw == "" {
		t.Error("StackTrace has no raw representation")
	}

	limited := st.Limit(1)
	if len(limited.frames) != 1 {
		t.Error("StackTrace.Limit(1) did not return 1 frame")
	}
}

func TestStackTraceString(t *testing.T) {
	config := Config{
		BufferSize: 2048,
		SkipFrames: 0,
	}

	st := NewStackTrace(&config)
	str := st.String()

	if str == "" {
		t.Error("StackTrace.String() returned empty string")
	}

	if !strings.Contains(str, "testing.tRunner") {
		t.Error("StackTrace.String() does not contain expected function name")
	}
}

func TestFramesString(t *testing.T) {
	frames := Frames{
		{Function: "main.main", File: "/path/to/main.go", Line: 10},
		{Function: "main.helper", File: "/path/to/helper.go", Line: 20},
	}

	str := frames.String()
	expected := "/path/to/main.go:10 main.main\n/path/to/helper.go:20 main.helper"

	if str != expected {
		t.Errorf("Frames.String() returned %q, want %q", str, expected)
	}
}

func TestFramesFilter(t *testing.T) {
	frames := Frames{
		{Function: "main.main", File: "/path/to/main.go", Line: 10},
		{Function: "invalid", File: "invalid.txt", Line: 20},
		{Function: "github.com/user/project/package.Function", File: "/path/to/helper.go", Line: 30},
	}

	filtered := frames.filter()

	if len(filtered) != 2 {
		t.Errorf("frames.filter() returned %d frames, want 2", len(filtered))
	}

	if filtered[0].Function != "main.main" || filtered[1].Function != "package.Function" {
		t.Errorf("frames.filter() did not correctly simplify function names: %v", filtered)
	}
}

func TestStackTraceFrames(t *testing.T) {
	config := Config{
		BufferSize: 2048,
		SkipFrames: 0,
	}

	st := NewStackTrace(&config)
	frames := st.Frames()

	if len(frames) == 0 {
		t.Error("StackTrace.Frames() returned no frames")
	}

	for _, frame := range frames {
		if frame.Function == "" || frame.File == "" || frame.Line == 0 {
			t.Error("StackTrace.Frames() returned an invalid frame")
		}
	}
}

func TestStackTraceLimit(t *testing.T) {
	config := Config{
		BufferSize: 2048,
		SkipFrames: 0,
	}

	st := NewStackTrace(&config)
	limitedST := st.Limit(2)

	if len(limitedST.frames) != 2 {
		t.Errorf("StackTrace.Limit(2) returned %d frames, want 2", len(limitedST.frames))
	}

	if limitedST.raw != st.raw {
		t.Error("StackTrace.Limit() modified the raw stack trace")
	}
}

func TestStackTraceWithConfig(t *testing.T) {
	config := Config{
		BufferSize: 1024,
		SkipFrames: 1,
	}

	st := NewStackTrace(&config)

	if len(st.raw) > 1024 {
		t.Error("StackTrace raw representation exceeds specified buffer size")
	}

	if strings.Contains(st.Frames().String(), "TestStackTraceWithConfig") {
		t.Error("StackTrace does not respect SkipFrames configuration")
	}
}

func TestFunctionNameRegexp(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"github.com/user/project/package.Function", "package.Function"},
		{"main.main", "main.main"},
		{"runtime.goexit", "runtime.goexit"},
	}

	for _, tc := range testCases {
		match := functionNameRegexp.FindStringSubmatch(tc.input)
		var result string
		if len(match) == 2 {
			result = match[1]
		} else {
			result = tc.input
		}
		if result != tc.expected {
			t.Errorf("functionNameRegexp.FindStringSubmatch(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestStackTraceEdgeCases(t *testing.T) {
	t.Run("ZeroBufferSize", func(t *testing.T) {
		st := NewStackTrace(&Config{BufferSize: 0, SkipFrames: 0})
		if st == nil {
			t.Error("NewStackTrace returned nil for zero buffer size")
		}
		if st != nil && len(st.frames) != 0 {
			t.Error("Expected no frames for zero buffer size")
		}
		if st != nil && st.raw != "" {
			t.Error("Expected empty raw stack trace for zero buffer size")
		}
	})

	t.Run("LargeSkipFrames", func(t *testing.T) {
		st := NewStackTrace(nil)
		if st == nil {
			t.Error("NewStackTrace returned nil for large SkipFrames")
		}
		if st != nil && len(st.frames) != 0 {
			t.Error("Expected no frames for large SkipFrames")
		}
	})
}

func BenchmarkNewStackTrace(b *testing.B) {
	config := Config{BufferSize: 2048, SkipFrames: 0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewStackTrace(&config)
	}
}

func ExampleNewStackTrace() {
	config := Config{BufferSize: 2048, SkipFrames: 0}
	st := NewStackTrace(&config)
	fmt.Println(st.Frames()[0].Function) // Print the first function name
}
