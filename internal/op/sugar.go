package op

import "runtime"

// Stack returns the current goroutine stack trace
func Stack() string {
	size := 512
	buf := make([]byte, size)
	length := runtime.Stack(buf, false)
	for length == len(buf) {
		size *= 2
		buf = make([]byte, size)
		length = runtime.Stack(buf, false)
	}
	return string(buf[:length])
}
