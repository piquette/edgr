package model

import (
	"bytes"
	"errors"
	"math"
	"time"
)

// Time wraps time.Time with serialization support.
type Time struct{ time.Time }

// UnmarshalJSON implements json.Unmarshaler
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Handle empty case.
	if bytes.Equal(data, []byte(`null`)) || bytes.Equal(data, []byte(`""`)) {
		*t = Time{}
		return nil
	}

	t.Time, err = time.Parse(`"`+time.RFC3339Nano+`"`, string(data))
	return
}

// MarshalJSON implements json.Marshaler
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(t.Format(`null`)), nil
	}

	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	return []byte(t.Format(`"` + time.RFC3339Nano + `"`)), nil
}

// TimeFromFloat converts timestamp floats to a Time object.
func TimeFromFloat(ts float64) Time {
	if ts == 0 {
		return Time{}
	}

	i, f := math.Modf(ts)
	return Time{time.Unix(int64(i), int64(f*1000000000.0))}
}

// Float converts a Time object to a timestamp float.
func (t Time) Float() float64 {
	if t.IsZero() {
		return 0
	}

	return float64(t.UnixNano()) / float64(1000000000.0)
}
