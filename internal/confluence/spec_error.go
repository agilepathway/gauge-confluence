package confluence

import "fmt"

type invalidSpecError struct {
	spec Spec
}

func (e *invalidSpecError) Error() string {
	return fmt.Sprintf("could not find a spec heading in spec %s", e.spec.path)
}

func (e *invalidSpecError) Nonfatal() bool {
	return true
}
