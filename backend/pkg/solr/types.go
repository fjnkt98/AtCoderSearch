package solr

import (
	"encoding/json"
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
	Header   ResponseHeader       `json:"responseHeader"`
	Response SelectBody           `json:"response"`
	Facets   RawJSONFacetResponse `json:"facets"`
	Error    ErrorInfo            `json:"error"`
}

type SelectBody struct {
	NumFound      int             `json:"numFound"`
	Start         int             `json:"start"`
	NumFoundExact bool            `json:"numFoundExact"`
	Docs          json.RawMessage `json:"docs"`
}
