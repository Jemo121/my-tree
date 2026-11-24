//go:build js && wasm

package assets

import (
	_ "embed"
	"encoding/json"
	"syscall/js"
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

// SignConfig holds localised strings for the wooden sign, read from canvas data-* attributes.
// Hugo injects these via i18n so the WASM binary stays language-agnostic.
type SignConfig struct {
	FamilyName  string
	Established string
	LabelActive string
	LabelRoots  string
	LabelEst    string
}

// LoadSignConfig reads sign strings from the #family-tree canvas element's dataset.
func LoadSignConfig() SignConfig {
	ds := js.Global().Get("document").Call("getElementById", "family-tree").Get("dataset")
	return SignConfig{
		FamilyName:  ds.Get("familyName").String(),
		Established: ds.Get("established").String(),
		LabelActive: ds.Get("signActive").String(),
		LabelRoots:  ds.Get("signRoots").String(),
		LabelEst:    ds.Get("signEst").String(),
	}
}

// ActiveCount returns the number of members with status "active".
func ActiveCount(members []RawMember) int {
	n := 0
	for _, m := range members {
		if m.Status == "active" {
			n++
		}
	}
	return n
}

// RootsCount returns the number of members with status "roots".
func RootsCount(members []RawMember) int {
	n := 0
	for _, m := range members {
		if m.Status == "roots" {
			n++
		}
	}
	return n
}
