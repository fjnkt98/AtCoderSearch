package solr

import (
	"time"
)

type SolrResponseHeader struct {
	ZkConnected map[string]any `json:"zkConnected"`
	Status      uint32         `json:"status"`
	QTime       uint32
	Params      map[string]any `json:"params"`
}

type SolrPingResponse struct {
	Header SolrResponseHeader `json:"responseHeader"`
	Status string             `json:"status"`
}

type SolrErrorInfo struct {
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
	Code     uint32   `json:"code"`
}

type LuceneInfo struct {
	SolrSpecVersion   string `json:"solr-spec-version"`
	SolrImplVersion   string `json:"solr-impl-version"`
	LuceneSpecVersion string `json:"lucene-spec-version"`
	LuceneImplVersion string `json:"lucene-impl-version"`
}

type SolrSystemInfo struct {
	Header   SolrResponseHeader `json:"responseHeader"`
	Mode     string             `json:"mode"`
	SolrHome string             `json:"solr_home"`
	CoreRoot string             `json:"core_root"`
	Lucene   LuceneInfo         `json:"lucene"`
	Jvm      map[string]any     `json:"jvm"`
	Security map[string]any     `json:"security"`
	System   map[string]any     `json:"system"`
	Error    SolrErrorInfo      `json:"error"`
}

type SolrIndexInfo struct {
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

type SolrCoreStatus struct {
	Name        string        `json:"name"`
	InstanceDir string        `json:"instanceDir"`
	DataDir     string        `json:"dataDir"`
	Config      string        `json:"config"`
	Schema      string        `json:"schema"`
	StartTime   string        `json:"startTime"`
	UpTime      uint64        `json:"uptime"`
	Index       SolrIndexInfo `json:"index"`
}

type SolrCoreList struct {
	Header       SolrResponseHeader        `json:"responseHeader"`
	InitFailures map[string]any            `json:"initFailures"`
	Status       map[string]SolrCoreStatus `json:"status"`
	Error        SolrErrorInfo             `json:"error"`
}

type SolrSimpleResponse struct {
	Header SolrResponseHeader `json:"responseHeader"`
	Error  SolrErrorInfo      `json:"error"`
}

type SolrSelectResponse[D any, F any] struct {
	Header      SolrResponseHeader `json:"responseHeader"`
	Response    SolrSelectBody[D]  `json:"response"`
	FacetCounts *F                 `json:"facets"`
	Error       *SolrErrorInfo     `json:"error"`
}

type SolrSelectBody[D any] struct {
	NumFound      uint32 `json:"numFound"`
	Start         uint32 `json:"start"`
	NumFoundExact bool   `json:"numFoundExact"`
	Docs          []D    `json:"docs"`
}

type BucketElement interface {
	int | int32 | int64 | uint32 | uint64 | float32 | float64 | time.Time | string
}

type Bucket[T BucketElement] struct {
	Val   T      `json:"val"`
	Count uint32 `json:"count"`
}

type SolrTermFacetCount struct {
	Buckets []Bucket[string] `json:"buckets"`
}

type SolrRangeFacetCount[T BucketElement] struct {
	Buckets []Bucket[T]              `json:"buckets"`
	Before  *SolrRangeFacetCountInfo `json:"before"`
	After   *SolrRangeFacetCountInfo `json:"after"`
	Between *SolrRangeFacetCountInfo `json:"between"`
}

type SolrRangeFacetCountInfo struct {
	Count uint32 `json:"count"`
}

type SolrQueryFacetCount struct {
	Buckets []Bucket[string] `json:"buckets"`
}