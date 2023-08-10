package user

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
)

func grade(c uint) string {
	if c < 10 {
		return "    ~  10"
	} else if c < 100 {
		return fmt.Sprintf("%d0  ~  %d9", c/10, c/10)
	} else {
		return fmt.Sprintf("%d00 ~ %d99", c/100, c/100)
	}
}

func (u *User) ToDocument() (Document, error) {
	color := acs.RateToColor(u.Rating)
	highestColor := acs.RateToColor(u.HighestRating)
	var period string
	if u.BirthYear != nil {
		period = fmt.Sprintf("%d0's", *u.BirthYear/10)
	}

	return Document{
		UserName:       u.UserName,
		Rating:         u.Rating,
		HighestRating:  u.HighestRating,
		Affiliation:    u.Affiliation,
		BirthYear:      u.BirthYear,
		Country:        u.Country,
		Crown:          u.Crown,
		JoinCount:      u.JoinCount,
		Rank:           u.Rank,
		ActiveRank:     u.ActiveRank,
		Wins:           u.Wins,
		Color:          color,
		HighestColor:   highestColor,
		Period:         period,
		JoinCountGrade: grade(u.JoinCount),
	}, nil
}

type Document struct {
	UserName       string  `json:"user_name"`
	Rating         int     `json:"rating"`
	HighestRating  int     `json:"highest_rating"`
	Affiliation    *string `json:"affiliation"`
	BirthYear      *uint   `json:"birth_year"`
	Country        *string `json:"country"`
	Crown          *string `json:"crown"`
	JoinCount      uint    `json:"join_count"`
	Rank           uint    `json:"rank"`
	ActiveRank     *uint   `json:"active_rank"`
	Wins           uint    `json:"wins" `
	Color          string  `json:"color"`
	HighestColor   string  `json:"highest_color"`
	Period         string  `json:"period,omitempty"`
	JoinCountGrade string  `json:"join_count_grade"`
}

type RowReader[R acs.ToDocument[D], D any] struct {
	db *sqlx.DB
}

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- *User) error {
	sql := `
	SELECT
		"user_name",
		"rating",
		"highest_rating",
		"affiliation",
		"birth_year",
		"country",
		"crown",
		"join_count",
		"rank",
		"active_rank",
		"wins"
	FROM
		"users"
	`
	rows, err := r.db.Queryx(sql)
	if err != nil {
		return failure.Translate(err, DBError, failure.Context{"sql": sql}, failure.Message("failed to read rows"))
	}
	defer rows.Close()
	defer close(tx)

	for rows.Next() {
		select {
		case <-ctx.Done():
			log.Println("ReadRows canceled.")
			return nil
		default:
			var row User
			err := rows.StructScan(&row)
			if err != nil {
				return failure.Translate(err, DBError, failure.Message("failed to scan row"))
			}
			tx <- &row
		}
	}

	return nil
}

type DocumentGenerator struct {
	saveDir string
	reader  *RowReader[*User, Document]
}

func NewDocumentGenerator(db *sqlx.DB, saveDir string) DocumentGenerator {
	return DocumentGenerator{
		saveDir: saveDir,
		reader:  &RowReader[*User, Document]{db: db},
	}
}

func (g *DocumentGenerator) Clean() error {
	if err := acs.CleanDocument(g.saveDir); err != nil {
		return failure.Translate(err, FileOperationError, failure.Context{"directory": g.saveDir}, failure.Message("failed to delete problem document files"))
	}
	return nil
}

func (g *DocumentGenerator) Generate(chunkSize int, concurrent int) error {
	if err := acs.GenerateDocument[*User, Document](g.reader, g.saveDir, chunkSize, concurrent); err != nil {
		return failure.Wrap(err)
	}
	return nil
}

func (g *DocumentGenerator) Run(chunkSize int, concurrent int) error {
	if err := g.Clean(); err != nil {
		return failure.Wrap(err)
	}

	if err := g.Generate(chunkSize, concurrent); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
