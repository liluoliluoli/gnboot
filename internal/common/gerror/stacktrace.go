package gerror

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type stacktrace struct {
	header []byte // goroutine header (e.g., "goroutine 3 [running]:")
	frames []byte // the stacktrace details
}

func dumpStacktrace(skip int) *stacktrace {
	skip += 2 // skip debug.Stack and this frame
	skip *= 2 // 2-line each frame
	stack := debug.Stack()
	var header []byte
	for i, b := range stack { // assumes no unicode in stack, iterate on bytes
		if b == '\n' {
			if header == nil {
				// consume first line as goroutine header
				header = stack[:i]
			} else {
				skip--
				if skip == 0 {
					stack = stack[i:]
					break
				}
			}
		}
	}
	if skip > 0 {
		panic("skip overflow")
	}
	return &stacktrace{header, stack}
}

//goland:noinspection GoUnhandledErrorResult
func (this *stacktrace) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		f.Write(this.header)
		f.Write(this.frames)
	}
}

func (this *stacktrace) String() string {
	sb := strings.Builder{}
	sb.Write(this.header)
	sb.Write(this.frames)
	return sb.String()
}
