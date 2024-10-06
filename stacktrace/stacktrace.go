// Package stacktrace provides functionality for capturing, filtering, and formatting stack traces in Go.
// It allows developers to generate structured stack traces with frames.
package stacktrace

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// Config holds the configuration for stack trace generation.
type Config struct {
	// BufferSize is the stack trace buffer size.
	BufferSize int
	// SkipFrames is the number of frames to skip.
	SkipFrames int
}

// DefaultConfig provides default configuration values.
var DefaultConfig = Config{
	BufferSize: 2048,
	SkipFrames: 2,
}

const (
	// validSuffix is the valid suffix for Go source files.
	validSuffix = ".go"
)

var (
	// functionNameRegexp is the regular expression used to extract function names.
	functionNameRegexp = regexp.MustCompile(`\/([^\/]+)$`)
)

// StackTrace represents a stack trace with frames and text representation.
type StackTrace struct {
	frames Frames // Filtered frames of the stack trace.
	raw    string // Raw text representation of the stack trace.
}

// Frames represents a collection of Frame objects.
type Frames []Frame

// filter removes invalid frames and normalizes function names for readability.
func (frames Frames) filter() Frames {
	filtered := make(Frames, 0, len(frames))
	for _, frame := range frames {
		if frame.File == "" || frame.Function == "" || frame.Line < 1 || !strings.HasSuffix(frame.File, validSuffix) {
			continue
		}
		// Simplify function name extraction
		functionName := frame.Function
		if match := functionNameRegexp.FindStringSubmatch(frame.Function); len(match) == 2 {
			functionName = match[1]
		}
		// Append the structured frame
		filtered = append(filtered, Frame{
			Function: functionName,
			File:     frame.File,
			Line:     frame.Line,
		})
	}
	return filtered
}

// String returns the string representation of the frames.
func (frames Frames) String() string {
	var builder strings.Builder
	for i, frame := range frames {
		if i > 0 {
			builder.WriteString("\n")
		}
		fmt.Fprintf(&builder, "%s:%d %s", frame.File, frame.Line, frame.Function)
	}
	return builder.String()
}

// Frame represents a single function call in the stack trace.
type Frame struct {
	Function string // Name of the function.
	File     string // File where the function is located.
	Line     int    // Line number in the file.
}

// NewStackTrace creates a new stack trace starting from the given skip level.
func NewStackTrace(config *Config) *StackTrace {
	stackTrace := &StackTrace{
		frames: make(Frames, 0),
		raw:    "",
	}

	// Use default config if not provided
	if config == nil {
		config = &DefaultConfig
	}

	// Get the stack trace
	uIntPtr := make([]uintptr, config.BufferSize)
	n := runtime.Callers(config.SkipFrames+2, uIntPtr) // +2 to skip runtime.Callers and NewStackTrace
	if n > 0 {
		uIntPtr = uIntPtr[:n]
		// Extract the structured frames
		frames := runtime.CallersFrames(uIntPtr)
		var structuredFrames Frames
		for {
			frame, more := frames.Next()
			// Append the structured frame
			structuredFrames = append(structuredFrames, Frame{
				Function: frame.Function,
				File:     frame.File,
				Line:     frame.Line,
			})
			// Break if no more frames
			if !more {
				break
			}
		}
		stackTrace.frames = structuredFrames.filter()
	}

	// Get the raw stack trace text
	buf := make([]byte, config.BufferSize)
	nBytes := runtime.Stack(buf, false)
	if nBytes > 0 {
		stackTrace.raw = strings.TrimSpace(string(buf[:nBytes]))
	}

	return stackTrace
}

// String returns the raw text stack trace.
func (stackTrace *StackTrace) String() string {
	return stackTrace.raw
}

// Frames returns the filtered frames of the stack trace.
func (stackTrace *StackTrace) Frames() Frames {
	return stackTrace.frames
}

// Limit returns a new StackTrace with at most n frames.
func (stackTrace *StackTrace) Limit(n int) *StackTrace {
	if n >= len(stackTrace.frames) {
		return stackTrace
	}
	return &StackTrace{
		frames: stackTrace.frames[:n],
		raw:    stackTrace.raw, // Note: raw string is not limited
	}
}
