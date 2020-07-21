package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Index is API for main page.
func (h *Handler) Index(c echo.Context) error {
	group := "github.com"
	scores, err := h.Counter.GetRankTotalByLimit(group, 20)
	if err != nil {
		return err
	}

	var ranks []string
	ranksMap := make(map[string]bool, 20)
	for _, score := range scores {
		if len(ranks) == 10 {
			break
		}

		path := strings.TrimSpace(score.Name)
		if strings.HasSuffix(path, "/") {
			path = path[:len(path)-1]
		}

		// add projects if score-name is /profile/project format
		seps := strings.Split(path, "/")
		if len(seps) == 3 && !ranksMap[path] {
			ranksMap[path] = true
			ranks = append(ranks, fmt.Sprintf("[%d] %s%s", len(ranks)+1, group, path))
		}
	}

	buf := new(bytes.Buffer)
	h.IndexTemplate.Execute(buf, struct {
		Ranks []string
	}{Ranks: ranks})
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}
