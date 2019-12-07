package core

import "encoding/json"

type User struct {
	UUID    string   `json:"uuid"`
	Email   string   `json:"email"`
	AlterId uint32   `json:"alter_id"`
	Regions []string `json:"region"`
	data    []byte   `json:"-"`
}

func (u *User) Encode() ([]byte, error) { return json.Marshal(u) }

func (u *User) Decode(data []byte) error {
	u.data = data
	return json.Unmarshal(data, u)
}

func (u *User) Data() []byte { return u.data }
