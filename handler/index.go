package handler

import (
	"net/http"

	"bytes"
	"fmt"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Index(c echo.Context) error {
	group := "github.com"
	scores, err := h.Counter.GetRankTotalByLimit(group, 10)
	if err != nil {
		return err
	}
	var ranks []string
	for i, score := range scores {
		ranks = append(ranks, fmt.Sprintf("(%d) %s%s : (%d count)", i+1, group, score.Name, score.Value))
	}

	buf := new(bytes.Buffer)
	h.IndexTemplate.Execute(buf, struct {
		Ranks []string
	}{Ranks: ranks})
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}
