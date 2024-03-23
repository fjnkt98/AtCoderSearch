package problem

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Parameter struct {
	Q              string   `json:"q" query:"q"`
	Limit          *int     `json:"limit" query:"limit"`
	Page           int      `json:"page" query:"page"`
	Sort           []string `json:"sort" query:"sort"`
	Facet          string   `json:"facet" query:"facet"`
	Category       []string `json:"category" query:"category"`
	DifficultyFrom int      `json:"difficultyFrom" query:"difficultyFrom"`
	DifficultyTo   int      `json:"difficultyTo" query:"difficultyTo"`
	Color          []string `json:"color" query:"color"`
}

func (p Parameter) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Q, validation.RuneLength(0, 200)),
		validation.Field(&p.Limit, validation.Min(0), validation.Max(200)),
		validation.Field(&p.Page, validation.Min(0)),
	)
}
