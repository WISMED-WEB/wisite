package module2

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	lk "github.com/digisan/logkit"
)

// after implementing, register with path in 'api_reg.go'

// @Title Test
// @Description POST test
// @Tags sign-in
// -- @Accept json                         # using form, comment '@Accept json'
// @Param name formData string true "Name"
// @Param age formData int true "Age"
// @Success 200 "OK"
// @Failure 400 "Fail"
// @Router /api/module2/test [post]
func TestPost(c echo.Context) error {
	name := c.FormValue("name")
	age := c.FormValue("age")
	lk.Log("%s-%s", name, age)
	return c.String(http.StatusOK, fmt.Sprintf("%s-%s", name, age))
}
