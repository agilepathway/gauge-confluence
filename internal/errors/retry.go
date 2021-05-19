package errors

import (
	"fmt"
	"time"

	"github.com/agilepathway/gauge-confluence/internal/logger"
)

// Retry allows flaky functionality to be retried.
// Based on: https://github.com/abourget/blog/blob/master/content/post/my-favorite-golang-retry-function.md
// Also see: https://stackoverflow.com/a/47606858
func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		logger.Debug(true, fmt.Sprintln("retrying after error:", err))
	}

	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
