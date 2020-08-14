package internal

import (
	"reflect"
	"testing"

	"github.com/gjbae1212/go-counter-badge/badge"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFlatBadge(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		leftText     string
		leftBgColor  string
		rightText    string
		rightBgColor string
		edgeFlat     bool
		output       badge.Badge
	}{
		"not-edge": {
			leftText:     "allan",
			leftBgColor:  "#555",
			rightText:    " 0 / 10 ",
			rightBgColor: "#79c83d",
			edgeFlat:     false,
			output: badge.Badge{
				FontType:             badge.VeraSans,
				LeftText:             "allan",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            " 0 / 10 ",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#79c83d",
				XRadius:              "3",
				YRadius:              "3",
			},
		},
		"edge": {
			leftText:     "allan",
			leftBgColor:  "#555",
			rightText:    " 0 / 10 ",
			rightBgColor: "#79c83d",
			edgeFlat:     true,
			output: badge.Badge{
				FontType:             badge.VeraSans,
				LeftText:             "allan",
				LeftTextColor:        "#fff",
				LeftBackgroundColor:  "#555",
				RightText:            " 0 / 10 ",
				RightTextColor:       "#fff",
				RightBackgroundColor: "#79c83d",
				XRadius:              "0",
				YRadius:              "0",
			},
		},
	}

	for _, t := range tests {
		bg := GenerateFlatBadge(t.leftText, t.leftBgColor, t.rightText, t.rightBgColor, t.edgeFlat)
		assert.True(reflect.DeepEqual(t.output, bg))
	}
}
