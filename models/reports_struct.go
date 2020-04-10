package models

import (
	"encoding/json"
	"time"
)

// UserReport is a submission report from a user
type UserReport struct {
	Reason string
	Count  int
}

// UnmarshalJSON helps to get UserReport directly from JSON
func (ur *UserReport) UnmarshalJSON(data []byte) error {
	var v []interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	ur.Reason, _ = v[0].(string)
	ur.Count = int(v[1].(float64))

	return nil
}

// ModReport is a submission report from a mod.
// Unlike UserReport, this includes the name of the mod
type ModReport struct {
	Reason string
	Mod    string
}

// UnmarshalJSON helps to get ModReport directly from JSON
func (mr *ModReport) UnmarshalJSON(data []byte) error {
	var v []interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	mr.Reason, _ = v[0].(string)
	mr.Mod = v[1].(string)

	return nil
}

// AllReports simply combines ModReports & UserReports
type AllReports struct {
	Num  int
	Mod  []ModReport
	User []UserReport
}

// SubModAction is a removal/approval with a timestamp
type SubModAction struct {
	Mod string
	At  time.Time
}
