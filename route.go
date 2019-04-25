package main

import (
	"fmt"

	echo "github.com/labstack/echo/v4"
)

func AddRoute(e *echo.Echo) error {
	if e == nil {
		return fmt.Errorf("[Err] AddRoute empty params")
	}
	return nil
}
