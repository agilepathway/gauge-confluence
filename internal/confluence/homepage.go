package confluence

import "github.com/agilepathway/gauge-confluence/internal/confluence/time"

type homepage struct {
	id        string
	created   time.Time
	childless bool
}
