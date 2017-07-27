package common

import (
	"time"
	"strings"
)

type JsonTimed time.Time

func (j JsonTimed) Time() time.Time {
	return time.Time(j)
}

func (j JsonTimed) String() string {
	return time.Time(j).String()
}

func (j *JsonTimed) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*j).Format("2006-01-02 15:04:05") + `"` ), nil
}

func (j *JsonTimed) UnmarshalJSON(p []byte) error {
	t, parseErr := time.ParseInLocation("2006-01-02 15:04:05", strings.Replace(string(p), `"`, "", -1), time.Local)
	if parseErr != nil {
		return parseErr
	}
	*j = JsonTimed(t)
	return nil
}

