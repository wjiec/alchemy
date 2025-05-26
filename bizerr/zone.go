package bizerr

import "sync/atomic"

// Zone manages a range of error codes, ensuring that new errors
// within the zone are unique and sequential.
type Zone struct {
	cursor uint32
	stop   uint32
}

// New creates a new Error instance within the Zone.
//
// It increments the cursor to obtain a new error code, ensuring
// it does not exceed the stop limit if set.
func (z *Zone) New(status uint32, template ...string) *Error {
	nextCode := atomic.AddUint32(&z.cursor, 1)
	if z.stop != 0 && nextCode >= z.stop {
		panic("too many errors in zone")
	}

	return New(nextCode, status, template...)
}

// NewZone initializes and returns a new Zone instance,
func NewZone(start uint32, stop ...uint32) *Zone {
	zone := &Zone{cursor: start - 1}
	if len(stop) > 0 {
		zone.stop = stop[0]
	}

	return zone
}

// Step defines an interface for managing steps in error code zones.
type Step interface {
	// Reset resets the current step's starting point.
	Reset(uint32)

	// Next retrieves the next Zone based on the current step.
	Next() *Zone
}

var (
	// ThousandStep step increments by 1000
	ThousandStep Step = &realStep{curr: 1000, step: 1000}
	// TenThousandStep step increments by 10000
	TenThousandStep Step = &realStep{curr: 10000, step: 10000}
)

// realStep implements the Step interface
type realStep struct {
	curr, step uint32
}

// Reset sets the current code to a specified value
func (r *realStep) Reset(curr uint32) {
	atomic.StoreUint32(&r.curr, curr)
}

// Next computes the next Zone based on the current step,
func (r *realStep) Next() *Zone {
	stop := atomic.AddUint32(&r.curr, r.step)
	return NewZone(stop-r.step, stop)
}
