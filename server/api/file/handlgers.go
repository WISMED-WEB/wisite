package file

import (
	"net/http"
	"path/filepath"

	fm "github.com/digisan/file-mgr"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/wismed-web/wisite-api/server/api/sign"
)

// *** after implementing, register with path in 'file.go' *** //

// @Title pathcontent
// @Summary get content under specific path.
// @Description
// @Tags    File
// @Accept  json
// @Produce json
// @Param   ym    query string true "year-month, e.g. 2022-05"
// @Param   gpath query string true "group path, e.g. group1/group2/group3"
// @Success 200 "OK - get content successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/file/pathcontent [get]
// @Security ApiKeyAuth
func PathContent(c echo.Context) error {

	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)

	var (
		uname = claims.UName
		ym    = c.QueryParam("ym")
		gpath = c.QueryParam("gpath")
	)

	// fetch user space for valid login
	us, ok := sign.MapUserSpace.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [pathcontent] @"+uname)
	}

	content := us.(*fm.UserSpace).PathContent(filepath.Join(ym, gpath))
	return c.JSON(http.StatusOK, content)
}

// @Title fileitem
// @Summary get fileitems by given path or id.
// @Description
// @Tags    File
// @Accept  json
// @Produce json
// @Param   id   query string true "file ID (md5)"
// @Success 200 "OK - get fileitems successfully"
// @Failure 400 "Fail - incorrect query param id"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/file/fileitems [get]
// @Security ApiKeyAuth
func FileItems(c echo.Context) error {

	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)

	var (
		uname = claims.UName
		id    = c.QueryParam("id")
	)

	// fetch user space for valid login
	us, ok := sign.MapUserSpace.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [fileitem] @"+uname)
	}

	fis, err := us.(*fm.UserSpace).FileItems(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if len(fis) == 0 {
		return c.JSON(http.StatusNotFound, fis)
	}

	return c.JSON(http.StatusOK, fis)
}

// @Title upload
// @Summary upload file action.
// @Description
// @Tags    File
// @Accept  multipart/form-data
// @Produce json
// @Param   note   formData string false "note for uploading file"
// @Param   group0 formData string false "1st category for uploading file"
// @Param   group1 formData string false "2nd category for uploading file"
// @Param   group2 formData string false "3rd category for uploading file"
// @Param   file   formData file   true  "file path for uploading"
// @Success 200 "OK - upload successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/file/upload [post]
// @Security ApiKeyAuth
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

	// fetch user space for valid login
	us, ok := sign.MapUserSpace.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [upload] @"+uname)
	}

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if _, err := us.(*fm.UserSpace).SaveFormFile(file, note, group0, group1, group2); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "uploaded successfully")
}
