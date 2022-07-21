package debug

import (
	"net/http"
	"os"

	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
)

// @Title erase all Post data (high risk, only for debugging)
// @Summary erase all Post data collected by wisite service (high risk, only for debugging)
// @Description
// @Tags    Debug
// @Accept  json
// @Produce json
// @Success 200 "OK - delete successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/debug/erase/all-post [delete]
func EraseAllPostData(c echo.Context) error {

	lk.Warn("Deleting Post Folders...")

	paths := []string{
		"./data/id-event",
		"./data/id-flwids",
		"./data/id-ptps",
		"./data/owner-ids",
		"./data/span-ids",
		"./data/user-fdb",
		"./data/user-space",
	}
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, "all data has been deleted")
}
