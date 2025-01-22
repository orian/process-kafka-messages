package message

import (
	"fmt"
	"regexp"
	"strconv"
)

type Field struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`

	Versions         Range `json:"versions,omitempty"`
	TaggedVersions   Range `json:"taggedVersions,omitempty"`
	NullableVersions Range `json:"nullableVersions,omitempty"`
	FlexibleVersions Range `json:"flexibleVersions,omitempty"`

	Ignorable bool   `json:"ignorable,omitempty"`
	About     string `json:"about,omitempty"`

	Tag    int  `json:"tag,omitempty"`
	MapKey bool `json:"mapKey,omitempty"`

	Fields []Field `json:"fields,omitempty"`

	Default    interface{} `json:"default,omitempty"`
	EntityType string      `json:"entityType,omitempty"`
	ZeroCopy   bool        `json:"zeroCopy,omitempty"`
}

type Range struct {
	Value string
	Begin int
	End   int
}

var (
	exactRe     = regexp.MustCompile(`^"?(\d+)"?$`)
	onlyStartRe = regexp.MustCompile(`^"?((\d+)\+)"?$`)
	startStopRe = regexp.MustCompile(`^"?((\d+)-(\d+))"?$`)
)

func (r *Range) UnmarshalJSON(bytes []byte) error {
	if subs := onlyStartRe.FindSubmatch(bytes); len(subs) == 3 {
		r.Value = string(subs[1])
		begin, err := strconv.Atoi(string(subs[2]))
		if err != nil {
			return err
		}
		r.Begin = begin
		return nil
	} else if subs := startStopRe.FindSubmatch(bytes); len(subs) == 4 {
		r.Value = string(subs[1])
		begin, err := strconv.Atoi(string(subs[2]))
		if err != nil {
			return err
		}
		r.Begin = begin

		end, err := strconv.Atoi(string(subs[3]))
		if err != nil {
			return err
		}
		r.End = end
		return nil
	} else if subs := exactRe.FindSubmatch(bytes); len(subs) == 2 {
		r.Value = string(subs[1])
		begin, err := strconv.Atoi(string(subs[1]))
		if err != nil {
			return err
		}
		r.Begin = begin
		r.End = begin
		return nil
	} else if string(bytes) == "none" || string(bytes) == "\"none\"" {
		r.Value = "none"
		return nil
	}
	return fmt.Errorf("invalid range: %s", string(bytes))
}

type Message struct {
	ApiKey                int      `json:"apiKey,omitempty"`
	Type                  string   `json:"type,omitempty"`
	Listeners             []string `json:"listeners,omitempty"`
	Name                  string   `json:"name,omitempty"`
	ValidVersions         Range    `json:"validVersions,omitempty"`
	LatestVersionUnstable bool     `json:"latestVersionUnstable,omitempty"`
	FlexibleVersions      Range    `json:"flexibleVersions,omitempty"`
	Fields                []Field
	CommonStructs         []Field `json:"commonStructs,omitempty"`
}
