package page

import (
	"strconv"
	"time"
)

type Date struct{ time.Time }

func ParseDate(data string) (Date, error) {
	v, err := strconv.Atoi(data)
	if err != nil {
		return Date{}, err
	}
	return Date{time.Unix(int64(v), 0)}, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(d.Unix()))), nil
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	*d, err = ParseDate(string(data))
	return
}
