package file

import (
	"net/http"

	fm "github.com/digisan/file-mgr"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'file.go' *** //

// @Title upload
// @Summary upload file action.
// @Description
// @Tags    file
// @Accept  multipart/form-data
// @Produce json
// @Param   note   formData string false "note for uploading file"
// @Param   group0 formData string false "1st category for uploading file"
// @Param   group1 formData string false "2nd category for uploading file"
// @Param   group2 formData string false "3rd category for uploading file"
// @Success 200 "OK - upload successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/file/upload [post]
func Upload(c echo.Context) error {

	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)

	// Read form fields
	var (
		uname  = claims.UName
		note   = c.FormValue("note")
		group0 = c.FormValue("group0")
		group1 = c.FormValue("group1")
		group2 = c.FormValue("group2")
	)

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	us, err := fm.UseUser(uname) // *** should refer from login ***
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err := us.SaveFile(file.Filename, note, src, group0, group1, group2); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "uploaded successfully")
}
