package solr

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshalFromSolrDateTime(t *testing.T) {
	var d FromSolrDateTime = FromSolrDateTime(time.Date(2023, time.Month(7), 18, 18, 31, 21, 0, time.Local))

	b, err := json.Marshal(d)
	if err != nil {
		t.Errorf("failed to marshal FromSolrDateTime `%+v`: %s", d, err.Error())
	}

	expected := `"2023-07-18T18:31:21+09:00"`
	if string(b) != expected {
		t.Errorf("mismatch result of marshalling FromSolrDateTime `%+v`: expected %s but got %s", d, expected, string(b))
	}
}

func TestUnmarshalFromSolrDateTime(t *testing.T) {
	type MyStruct struct {
		Time FromSolrDateTime `json:"time"`
	}

	source := `{"time":"2023-07-18T09:31:21Z"}`
	var actual MyStruct
	if err := json.Unmarshal([]byte(source), &actual); err != nil {
		t.Errorf("failed to unmarshal SolrDateTime from `%s`: %s", source, err.Error())
		return
	}
	expected := MyStruct{Time: FromSolrDateTime(time.Date(2023, time.Month(7), 18, 18, 31, 21, 0, time.Local))}
	if actual != expected {
		t.Errorf("mismatch result of un-marshalling FromSolrDateTime `%s`: expected %+v but got %+v", source, expected, actual)
	}
}

func TestMarshalIntoSolrDateTime(t *testing.T) {
	var d IntoSolrDateTime = IntoSolrDateTime(time.Date(2023, time.Month(7), 18, 18, 31, 21, 0, time.Local))

	b, err := json.Marshal(d)
	if err != nil {
		t.Errorf("failed to marshal IntoSolrDateTime `%+v`: %s", d, err.Error())
		return
	}

	expected := `"2023-07-18T09:31:21Z"`
	if string(b) != expected {
		t.Errorf("mismatch result of marshalling IntoSolrDateTime `%+v`: expected %s but got %s", d, expected, string(b))
	}
}

func TestUnmarshalIntoSolrDateTime(t *testing.T) {
	type MyStruct struct {
		Time IntoSolrDateTime `json:"time"`
	}

	source := `{"time":"2023-07-18T18:31:21+09:00"}`
	var actual MyStruct
	if err := json.Unmarshal([]byte(source), &actual); err != nil {
		t.Errorf("failed to unmarshal SolrDateTime from `%s`: %s", source, err.Error())
		return
	}
	expected := MyStruct{Time: IntoSolrDateTime(time.Date(2023, time.Month(7), 18, 18, 31, 21, 0, time.Local))}
	if actual != expected {
		t.Errorf("mismatch result of un-marshalling IntoSolrDateTime `%s`: expected %+v but got %+v", source, expected, actual)
	}
}
