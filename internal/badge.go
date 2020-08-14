package internal

import "github.com/gjbae1212/go-counter-badge/badge"

// GenerateFlatBadge makes Flat-Badge struct which is used go-counter-badge/badge.
func GenerateFlatBadge(leftText, leftBgColor, rightText, rightBgColor string, edgeRound bool) badge.Badge {
	flatBadge := badge.Badge{
		FontType:             badge.VeraSans,
		LeftText:             leftText,
		LeftTextColor:        "#fff",
		LeftBackgroundColor:  leftBgColor,
		RightText:            rightText,
		RightTextColor:       "#fff",
		RightBackgroundColor: rightBgColor,
	}
	if edgeRound {
		flatBadge.XRadius = "3"
		flatBadge.YRadius = "3"
	} else {
		flatBadge.XRadius = "0"
		flatBadge.YRadius = "0"
	}
	return flatBadge
}
