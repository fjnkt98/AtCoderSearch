package solr

import (
	"encoding/json"
	"time"
)

type FromSolrDateTime time.Time

func (t FromSolrDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

func (t *FromSolrDateTime) UnmarshalJSON(data []byte) error {
	dataString := string(data)

	if dataString == "null" {
		return nil
	}

	parsed, err := time.ParseInLocation(`"2006-01-02T15:04:05Z"`, dataString, time.UTC)
	if err != nil {
		return err
	}

	*t = FromSolrDateTime(parsed.Local())
	return nil
}

type IntoSolrDateTime time.Time

func (t IntoSolrDateTime) String() string {
	return time.Time(t).UTC().Format(`2006-01-02T15:04:05Z`)
}

func (t IntoSolrDateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).UTC().Format(`"2006-01-02T15:04:05Z"`)), nil
}

func (t *IntoSolrDateTime) UnmarshalJSON(data []byte) error {
	var d time.Time
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	*t = IntoSolrDateTime(d.Local())
	return nil
}
