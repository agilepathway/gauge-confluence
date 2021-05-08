// Package errors implements functions to manipulate errors
package errors

// nonfatal allows errors to indicate if they are fatal or not
// this is the "assert errors for behaviour, not type" principle, see:
// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
type nonfatal interface {
	Nonfatal() bool
}

// IsNonfatal returns true if err is nonfatal
func IsNonfatal(err error) bool {
	nf, ok := err.(nonfatal)
	return ok && nf.Nonfatal()
}
