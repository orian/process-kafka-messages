package message

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tailscale/hujson"
)

func TestRange_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		Value string
		Want  Range
		Err   string
	}{
		{"7", Range{"7", 7, 7}, ""},
		{"\"7\"", Range{"7", 7, 7}, ""},
		{"1+", Range{"1+", 1, 0}, ""},
		{"2-4", Range{"2-4", 2, 4}, ""},
		{"none", Range{"none", 0, 0}, ""},
	}

	for _, c := range cases {
		t.Run(c.Value, func(t *testing.T) {
			var got Range
			err := got.UnmarshalJSON([]byte(c.Value))
			if c.Err == "" {
				assert.NoError(t, err)
				assert.Equal(t, c.Want, got)
			} else {
				assert.EqualError(t, err, c.Err)
			}
		})
	}
}

func TestMessageParse(t *testing.T) {
	for _, c := range []struct {
		file string
		want Message
	}{
		{"testdata/ApiVersionsRequest.json", Message{
			ApiKey:    18,
			Type:      "request",
			Listeners: []string{"zkBroker", "broker", "controller"},
			Name:      "ApiVersionsRequest",
			ValidVersions: Range{
				Value: "0-4",
				Begin: 0,
				End:   4,
			},
			FlexibleVersions: Range{
				Value: "3+",
				Begin: 3,
				End:   0,
			},
			Fields: []Field{
				{
					Name: "ClientSoftwareName", Type: "string",
					Versions:  Range{Value: "3+", Begin: 3, End: 0},
					Ignorable: true, About: "The name of the client.",
				},
				{
					Name: "ClientSoftwareVersion", Type: "string",
					Versions:  Range{Value: "3+", Begin: 3, End: 0},
					Ignorable: true, About: "The version of the client.",
				},
			},
		}},
		{
			file: "testdata/ApiVersionsResponse.json",
			want: Message{
				ApiKey:           18,
				Type:             "response",
				Name:             "ApiVersionsResponse",
				ValidVersions:    Range{Value: "0-4", Begin: 0, End: 4},
				FlexibleVersions: Range{Value: "3+", Begin: 3, End: 0},
				Fields: []Field{
					{
						Name: "ErrorCode", Type: "int16",
						Versions:  Range{Value: "0+", Begin: 0, End: 0},
						Ignorable: false,
						About:     "The top-level error code.",
					},
					{
						Name:      "ApiKeys",
						Type:      "[]ApiVersion",
						Versions:  Range{Value: "0+", Begin: 0, End: 0},
						Ignorable: false,
						About:     "The APIs supported by the broker.",
						Fields: []Field{
							{
								Name: "ApiKey", Type: "int16", Versions: Range{Value: "0+", Begin: 0, End: 0}, MapKey: true, Ignorable: false, About: "The API index.",
							},
							{
								Name: "MinVersion", Type: "int16", Versions: Range{Value: "0+", Begin: 0, End: 0}, Ignorable: false, About: "The minimum supported version, inclusive.",
							},
							{
								Name: "MaxVersion", Type: "int16", Versions: Range{Value: "0+", Begin: 0, End: 0}, Ignorable: false, About: "The maximum supported version, inclusive.",
							},
						},
					},
					{
						Name:      "ThrottleTimeMs",
						Type:      "int32",
						Versions:  Range{Value: "1+", Begin: 1, End: 0},
						Ignorable: true,
						About:     "The duration in milliseconds for which the request was throttled due to a quota violation, or zero if the request did not violate any quota.",
					},
					{
						Name:           "SupportedFeatures",
						Type:           "[]SupportedFeatureKey",
						Versions:       Range{Value: "3+", Begin: 3, End: 0},
						Ignorable:      true,
						About:          "Features supported by the broker. Note: in v0-v3, features with MinSupportedVersion = 0 are omitted.",
						TaggedVersions: Range{Value: "3+", Begin: 3, End: 0},
						Fields: []Field{
							{
								Name: "Name", Type: "string", Versions: Range{Value: "3+", Begin: 3, End: 0}, MapKey: true, Ignorable: false, About: "The name of the feature.",
							},
							{
								Name: "MinVersion", Type: "int16", Versions: Range{Value: "3+", Begin: 3, End: 0}, Ignorable: false, About: "The minimum supported version for the feature.",
							},
							{
								Name: "MaxVersion", Type: "int16", Versions: Range{Value: "3+", Begin: 3, End: 0}, Ignorable: false, About: "The maximum supported version for the feature.",
							},
						},
					},
					{
						Name:           "FinalizedFeaturesEpoch",
						Type:           "int64",
						Versions:       Range{Value: "3+", Begin: 3, End: 0},
						Ignorable:      true,
						About:          "The monotonically increasing epoch for the finalized features information. Valid values are >= 0. A value of -1 is special and represents unknown epoch.",
						Tag:            1,
						TaggedVersions: Range{Value: "3+", Begin: 3, End: 0},
						Default:        "-1",
					},
					{
						Name:           "FinalizedFeatures",
						Type:           "[]FinalizedFeatureKey",
						Versions:       Range{Value: "3+", Begin: 3, End: 0},
						Ignorable:      true,
						About:          "List of cluster-wide finalized features. The information is valid only if FinalizedFeaturesEpoch >= 0.",
						Tag:            2,
						TaggedVersions: Range{Value: "3+", Begin: 3, End: 0},
						Fields: []Field{
							{
								Name: "Name", Type: "string", Versions: Range{Value: "3+", Begin: 3, End: 0}, MapKey: true, Ignorable: false, About: "The name of the feature.",
							},
							{
								Name: "MaxVersionLevel", Type: "int16", Versions: Range{Value: "3+", Begin: 3, End: 0}, Ignorable: false, About: "The cluster-wide finalized max version level for the feature.",
							},
							{
								Name: "MinVersionLevel", Type: "int16", Versions: Range{Value: "3+", Begin: 3, End: 0}, Ignorable: false, About: "The cluster-wide finalized min version level for the feature.",
							},
						},
					},
					{
						Name:           "ZkMigrationReady",
						Type:           "bool",
						Versions:       Range{Value: "3+", Begin: 3, End: 0},
						Ignorable:      true,
						About:          "Set by a KRaft controller if the required configurations for ZK migration are present.",
						Tag:            3,
						TaggedVersions: Range{Value: "3+", Begin: 3, End: 0},
						Default:        "false",
					}},
			},
		},
		{
			file: "testdata/CreateTopicsRequest.json",
		},
		{
			file: "testdata/CreateTopicsResponse.json",
		},
		{
			file: "testdata/ProduceRequest.json",
		},
		{
			file: "testdata/ProduceResponse.json",
		},
	} {
		t.Run(c.file, func(t *testing.T) {
			data, err := os.ReadFile(c.file)
			assert.NoError(t, err)

			clean, err := hujson.Standardize(data)
			assert.NoError(t, err)
			var got Message

			dec := json.NewDecoder(bytes.NewBuffer(clean))
			dec.DisallowUnknownFields()
			assert.NoError(t, dec.Decode(&got))
			if c.want.ApiKey > 0 {
				//log.Printf("%#v", got)
				assert.Equal(t, c.want, got)
			}
		})
	}
}
