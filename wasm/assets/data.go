package assets

import (
	_ "embed"
	"encoding/json"
)

//go:embed members.json
var membersJSON []byte

// RawMember mirrors the members.yaml schema (converted to JSON at build time).
type RawMember struct {
	Name       string  `json:"name"`
	Subdomain  string  `json:"subdomain"`
	Template   string  `json:"template"`
	Generation int     `json:"generation"`
	Parent     *string `json:"parent"`
	Status     string  `json:"status"`
	Icon       string  `json:"icon"`
}

// LoadMembers parses the embedded members.json and returns all records.
func LoadMembers() ([]RawMember, error) {
	var members []RawMember
	if err := json.Unmarshal(membersJSON, &members); err != nil {
		return nil, err
	}
	return members, nil
}
