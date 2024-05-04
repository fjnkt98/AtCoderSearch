package list

import (
	"encoding/json"
	"fjnkt98/atcodersearch/repository"
	"log/slog"
	"net/http"

	"github.com/coocood/freecache"
	cache "github.com/gitsight/go-echo-cache"
	"github.com/goark/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ListHandler struct {
	pool *pgxpool.Pool
}

func NewListHandler(pool *pgxpool.Pool) *ListHandler {
	return &ListHandler{
		pool: pool,
	}
}

type ListResponse[T any] struct {
	Items   []T    `json:"items"`
	Message string `json:"message,omitempty"`
}

func NewErrorListResponse(message string) ListResponse[any] {
	return ListResponse[any]{
		Items:   []any{},
		Message: message,
	}
}

type ListProblemParameter struct {
	ContestID []string `query:"contestId"`
}

func (p ListProblemParameter) Validate() error {
	return nil
}

func (h *ListHandler) ListProblem(ctx echo.Context) error {
	var p ListProblemParameter
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorListResponse("bad request"))
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorListResponse(err.Error()))
	}

	q := repository.New(h.pool)

	var rows []string
	var err error
	if len(p.ContestID) > 0 {
		rows, err = q.FetchProblemIDsByContestID(ctx.Request().Context(), p.ContestID)
	} else {
		rows, err = q.FetchProblemIDs(ctx.Request().Context())
	}
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusOK, ListResponse[any]{})
		} else {
			slog.Error("failed to fetch problems", slog.Any("error", err))
			return ctx.JSON(http.StatusInternalServerError, NewErrorListResponse("request failed"))
		}
	}
	return ctx.JSON(http.StatusOK, ListResponse[string]{Items: rows})
}

type ListContestParameter struct {
	Category []string `query:"category"`
}

func (p ListContestParameter) Validate() error {
	return nil
}

func (h *ListHandler) ListContest(ctx echo.Context) error {
	var p ListContestParameter
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorListResponse("bad request"))
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorListResponse(err.Error()))
	}

	q := repository.New(h.pool)

	var rows []string
	var err error
	if len(p.Category) > 0 {
		rows, err = q.FetchContestIDsByCategory(ctx.Request().Context(), p.Category)
	} else {
		rows, err = q.FetchContestIDs(ctx.Request().Context())
	}
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusOK, ListResponse[any]{})
		} else {
			slog.Error("failed to fetch contests", slog.Any("error", err))
			return ctx.JSON(http.StatusInternalServerError, NewErrorListResponse("request failed"))
		}
	}
	return ctx.JSON(http.StatusOK, ListResponse[string]{Items: rows})
}

func (h *ListHandler) ListCategory(ctx echo.Context) error {
	q := repository.New(h.pool)

	rows, err := q.FetchCategories(ctx.Request().Context())
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusOK, ListResponse[any]{})
		} else {
			slog.Error("failed to fetch categories", slog.Any("error", err))
			return ctx.JSON(http.StatusInternalServerError, NewErrorListResponse("request failed"))
		}
	}
	return ctx.JSON(http.StatusOK, ListResponse[string]{Items: rows})
}

type ListLanguageParameter struct {
	Group []string `query:"group"`
}

func (p ListLanguageParameter) Validate() error {
	return nil
}

type Language struct {
	Group    string   `json:"group"`
	Language []string `json:"language"`
}

func (h *ListHandler) ListLanguage(ctx echo.Context) error {
	var p ListLanguageParameter
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorListResponse("bad request"))
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorListResponse(err.Error()))
	}

	q := repository.New(h.pool)

	var row []byte
	var err error
	if len(p.Group) > 0 {
		row, err = q.FetchLanguagesByGroup(ctx.Request().Context(), p.Group)
	} else {
		row, err = q.FetchLanguages(ctx.Request().Context())
	}
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusOK, ListResponse[Language]{})
		} else {
			slog.Error("failed to fetch languages", slog.Any("error", err))
			return ctx.JSON(http.StatusInternalServerError, NewErrorListResponse("request failed"))
		}
	}

	items := make([]Language, 0)
	if len(row) != 0 {
		if err := json.Unmarshal(row, &items); err != nil {
			return ctx.JSON(http.StatusInternalServerError, NewErrorListResponse("request failed"))
		}
	}

	return ctx.JSON(http.StatusOK, ListResponse[Language]{Items: items})
}

func (h *ListHandler) Register(e *echo.Echo) {
	g := e.Group("/api/list")

	c := freecache.NewCache(256 * 1024 * 1024)
	g.Use(cache.New(&cache.Config{}, c))
	g.GET("/problem", h.ListProblem)
	g.GET("/contest", h.ListContest)
	g.GET("/category", h.ListCategory)
	g.GET("/language", h.ListLanguage)
}