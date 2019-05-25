package api_handler

import (
	"fmt"

	"github.com/gjbae1212/hit-counter/handler"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GithubRankOfTotal(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)

	group := "github.com"
	scores, err := h.Counter.GetRankTotalByLimit(group, 10)
	if err != nil {
		return err
	}
	var result []string
	for i, score := range scores {
		result = append(result, fmt.Sprintf("(%d) %s%s <%d>", i+1, group, score.Name, score.Value))
	}
	return hctx.JSON(200, result)
}
