package solr

import (
	"encoding/json"
	"time"
)

type ResponseHeader struct {
	ZkConnected map[string]any `json:"zkConnected"`
	Status      uint           `json:"status"`
	QTime       uint
	Params      map[string]any `json:"params"`
}

type PingResponse struct {
	Header ResponseHeader `json:"responseHeader"`
	Status string         `json:"status"`
}

type ErrorInfo struct {
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
	Code     uint     `json:"code"`
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
	NumDocs                 uint64                 `json:"numDocs"`
	MaxDoc                  uint64                 `json:"maxDoc"`
	DeletedDocs             uint64                 `json:"deletedDocs"`
	Version                 uint64                 `json:"version"`
	SegmentCount            uint64                 `json:"segmentCount"`
	Current                 bool                   `json:"current"`
	HasDeletions            bool                   `json:"hasDeletions"`
	Directory               string                 `json:"directory"`
	SegmentsFile            string                 `json:"segmentsFile"`
	SegmentsFileSizeInBytes uint64                 `json:"segmentsFileSizeInBytes"`
	UserData                map[string]interface{} `json:"userData"`
	SizeInBytes             uint64                 `json:"sizeInBites"`
	Size                    string                 `json:"size"`
}

type CoreStatus struct {
	Name        string    `json:"name"`
	InstanceDir string    `json:"instanceDir"`
	DataDir     string    `json:"dataDir"`
	Config      string    `json:"config"`
	Schema      string    `json:"schema"`
	StartTime   string    `json:"startTime"`
	UpTime      uint64    `json:"uptime"`
	Index       IndexInfo `json:"index"`
}

type CoreList struct {
	Header       ResponseHeader        `json:"responseHeader"`
	InitFailures map[string]any        `json:"initFailures"`
	Status       map[string]CoreStatus `json:"status"`
	Error        ErrorInfo             `json:"error"`
}

type SimpleResponse struct {
	Header ResponseHeader `json:"responseHeader"`
	Error  ErrorInfo      `json:"error"`
}

type SelectResponse[D any, F any] struct {
	Header      ResponseHeader `json:"responseHeader"`
	Response    SelectBody[D]  `json:"response"`
	FacetCounts F              `json:"facets"`
	Error       ErrorInfo      `json:"error"`
}

type SelectBody[D any] struct {
	NumFound      uint `json:"numFound"`
	Start         uint `json:"start"`
	NumFoundExact bool `json:"numFoundExact"`
	Docs          []D  `json:"docs"`
}

type BucketElement interface {
	int | uint | float64 | time.Time | string
}

type Bucket[T BucketElement] struct {
	Val   T    `json:"val"`
	Count uint `json:"count"`
}

type TermFacetCount struct {
	Buckets []Bucket[string] `json:"buckets"`
}

type RangeFacetCount[T BucketElement] struct {
	Buckets []Bucket[T]              `json:"buckets"`
	Before  *RangeFacetCountInfo `json:"before"`
	After   *RangeFacetCountInfo `json:"after"`
	Between *RangeFacetCountInfo `json:"between"`
}

type RangeFacetCountInfo struct {
	Count uint `json:"count"`
}

type QueryFacetCount struct {
	Buckets []Bucket[string] `json:"buckets"`
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
	return time.Time(t).UTC().Format(`"2006-01-02T15:04:05Z"`)
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
