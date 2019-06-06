package podio

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Marshals and Unmarshals time in the common podio format.
// Podio time is always UTC.
type Time struct {
	time.Time
}

const podioDatetimeLayout = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(buf []byte) error {
	// apparently we need to trim "
	raw := strings.Trim(string(buf), "\"")

	if raw == "null" {
		// on null value we set the time to the time.Time zero value
		t.Time = time.Time{}
		return nil
	}

	tm, err := time.ParseInLocation(podioDatetimeLayout, raw, time.UTC)
	if err == nil {
		t.Time = tm
	}
	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	s := t.Format(podioDatetimeLayout)
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

type Date struct {
	time.Time
}

const podioDateLayout = "2006-01-02"

func (d *Date) UnmarshalJSON(buf []byte) error {
	// apparently we need to trim "
	raw := strings.Trim(string(buf), "\"")

	if raw == "null" {
		// on null value we set the time to the time.Time zero value
		d.Time = time.Time{}
		return nil
	}

	tm, err := time.ParseInLocation(podioDateLayout, raw, time.UTC)
	if err == nil {
		d.Time = tm
	}
	return err
}

func (d *Date) MarshalJSON() ([]byte, error) {
	s := d.Format(podioDateLayout)
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	t.Time = time.Unix(int64(ts), 0)

	return nil
}
