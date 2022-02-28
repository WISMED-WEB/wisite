package file

import (
	"net/http"

	fm "github.com/digisan/file-mgr"
	"github.com/digisan/file-mgr/fdb"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/wismed-web/wisite-api/server/api/sign"
)

// *** after implementing, register with path in 'file.go' *** //

// @Title pathcontent
// @Summary get content under specific path.
// @Description
// @Tags    file
// @Accept  json
// @Produce json
// @Param   path query string false "path to some level"
// @Success 200 "OK - upload successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/file/pathcontent [get]
func PathContent(c echo.Context) error {

	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)

	var (
		uname = claims.UName
		path  = c.QueryParam("path")
	)

	// fetch user space for valid login
	us, ok := sign.MapUserSpace.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [pathcontent] @"+uname)
	}

	content := us.(*fm.UserSpace).PathContent(path)
	return c.JSON(http.StatusOK, content)
}

// @Title fileitem
// @Summary get fileitems by given path or id.
// @Description
// @Tags    file
// @Accept  json
// @Produce json
// @Param   path query string false "path to a file"
// @Param   id   query string false "file's id"
// @Success 200 "OK - get fileitems successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/file/fileitem [get]
func FileItem(c echo.Context) error {

	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)

	var (
		uname = claims.UName
		path  = c.QueryParam("path")
		id    = c.QueryParam("id")
	)

	// fetch user space for valid login
	us, ok := sign.MapUserSpace.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [fileitem] @"+uname)
	}

	var fis []*fdb.FileItem
	if fi := us.(*fm.UserSpace).FileItemByPath(path); fi != nil {
		fis = append(fis, fi)
	}
	fis = append(fis, us.(*fm.UserSpace).FileItemByID(id)...)

	// remove duplicated fileitem
	m := map[string]*fdb.FileItem{}
	for _, fi := range fis {
		m[fi.Id+fi.Path] = fi
	}

	fis = []*fdb.FileItem{}
	for _, v := range m {
		fis = append(fis, v)
	}

	return c.JSON(http.StatusOK, fis)
}

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
// @Param   file   formData file   true  "file path for uploading"
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
	if err := us.(*fm.UserSpace).SaveFormFile(file, note, group0, group1, group2); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "uploaded successfully")
}
