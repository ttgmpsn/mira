package models

import "encoding/json"

type UserReport struct {
	Reason string
	Count  int
}

func (ur *UserReport) UnmarshalJSON(data []byte) error {
	var v []interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	ur.Reason, _ = v[0].(string)
	ur.Count = int(v[1].(float64))

	return nil
}

type ModReport struct {
	Reason string
	Mod    string
}

func (mr *ModReport) UnmarshalJSON(data []byte) error {
	var v []interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	mr.Reason, _ = v[0].(string)
	mr.Mod = v[1].(string)

	return nil
}
