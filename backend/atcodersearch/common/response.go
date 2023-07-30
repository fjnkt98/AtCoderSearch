package common

type SearchResultResponse[D any] struct {
	Stats   SearchResultStats `json:"status,omitempty"`
	Items   []D               `json:"items"`
	Message string            `json:"message,omitempty"`
}

type SearchResultStats struct {
	Time   uint `json:"time"`
	Total  uint `json:"total"`
	Index  uint `json:"index"`
	Pages  uint `json:"pages"`
	Count  uint `json:"count"`
	Params any  `json:"params,omitempty"`
	Facet  any  `json:"facet,omitempty"`
}

func NewErrorResponse[D any](msg string, params any) SearchResultResponse[D] {
	return SearchResultResponse[D]{
		Message: msg,
	}
}

type QueryLog struct {
	Domain string `json:"domain"`
	Time   uint   `json:"time"`
	Hits   uint   `json:"hits"`
	Params any    `json:"params"`
}
