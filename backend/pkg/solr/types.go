package solr

import (
	"encoding/json"
	"time"
)

type ResponseHeader struct {
	ZkConnected map[string]any `json:"zkConnected"`
	Status      int            `json:"status"`
	QTime       int            `json:"QTime"`
	Params      map[string]any `json:"params"`
}

type PingResponse struct {
	Header ResponseHeader `json:"responseHeader"`
	Status string         `json:"status"`
}

type ErrorInfo struct {
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
	Code     int      `json:"code"`
}

type LuceneInfo struct {
	SolrSpecVersion   string `json:"solr-spec-version"`
	SolrImplVersion   string `json:"solr-impl-version"`
	LuceneSpecVersion string `json:"lucene-spec-version"`
	LuceneImplVersion string `json:"lucene-impl-version"`
}

type SystemInfo struct {
	Header   ResponseHeader `json:"responseHeader"`
	Mode     string         `json:"mode"`
	SolrHome string         `json:"solr_home"`
	CoreRoot string         `json:"core_root"`
	Lucene   LuceneInfo     `json:"lucene"`
	Jvm      map[string]any `json:"jvm"`
	Security map[string]any `json:"security"`
	System   map[string]any `json:"system"`
	Error    ErrorInfo      `json:"error"`
}

type IndexInfo struct {
	NumDocs                 int64                  `json:"numDocs"`
	MaxDoc                  int64                  `json:"maxDoc"`
	DeletedDocs             int64                  `json:"deletedDocs"`
	Version                 int64                  `json:"version"`
	SegmentCount            int64                  `json:"segmentCount"`
	Current                 bool                   `json:"current"`
	HasDeletions            bool                   `json:"hasDeletions"`
	Directory               string                 `json:"directory"`
	SegmentsFile            string                 `json:"segmentsFile"`
	SegmentsFileSizeInBytes int64                  `json:"segmentsFileSizeInBytes"`
	UserData                map[string]interface{} `json:"userData"`
	SizeInBytes             int64                  `json:"sizeInBites"`
	Size                    string                 `json:"size"`
}

type CoreStatus struct {
	Name        string    `json:"name"`
	InstanceDir string    `json:"instanceDir"`
	DataDir     string    `json:"dataDir"`
	Config      string    `json:"config"`
	Schema      string    `json:"schema"`
	StartTime   string    `json:"startTime"`
	UpTime      int64     `json:"uptime"`
	Index       IndexInfo `json:"index"`
}

type CoreStatuses struct {
	Header       ResponseHeader        `json:"responseHeader"`
	InitFailures map[string]any        `json:"initFailures"`
	Status       map[string]CoreStatus `json:"status"`
	Error        ErrorInfo             `json:"error"`
}

type SimpleResponse struct {
	Header ResponseHeader `json:"responseHeader"`
	Error  ErrorInfo      `json:"error"`
}

type SelectResponse struct {
	Header   ResponseHeader `json:"responseHeader"`
	Response SelectBody     `json:"response"`
	Facets   Facets         `json:"facets"`
	Error    ErrorInfo      `json:"error"`
}

type SelectBody struct {
	NumFound      int             `json:"numFound"`
	Start         int             `json:"start"`
	NumFoundExact bool            `json:"numFoundExact"`
	Docs          json.RawMessage `json:"docs"`
}

type TermFacetCount struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Val    string                    `json:"val"`
	Count  int                       `json:"count"`
	Nested map[string]TermFacetCount `json:"nested"`
}

func (b *Bucket) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var bucket Bucket
	nested := make(map[string]TermFacetCount)
	for k, v := range raw {
		if k == "val" {
			var val string
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			bucket.Val = val
		} else if k == "count" {
			var count int
			if err := json.Unmarshal(v, &count); err != nil {
				return err
			}
			bucket.Count = count
		} else {
			var c TermFacetCount
			if err := json.Unmarshal(v, &c); err != nil {
				return err
			}
			nested[k] = c
		}
	}
	bucket.Nested = nested
	*b = bucket
	return nil
}

type Facets map[string]TermFacetCount

func (f *Facets) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	res := make(map[string]TermFacetCount)
	for k, v := range raw {
		if k == "count" {
			continue
		}
		var c TermFacetCount
		if err := json.Unmarshal(v, &c); err != nil {
			continue
		}
		res[k] = c
	}

	*f = Facets(res)
	return nil
}

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
