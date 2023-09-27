package sp

import "time"

type SpannerArguments struct {
	Verbose           bool
	Database          string // fullname
	File              string
	PartitionedUpdate bool
	Timeout           time.Duration
}
