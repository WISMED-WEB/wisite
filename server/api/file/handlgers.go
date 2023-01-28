package file

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	fm "github.com/digisan/file-mgr"
	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt/v4"
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
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		ym      = c.QueryParam("ym")
		gpath   = c.QueryParam("gpath")
	)

	// fetch user space for valid login
	us, ok := sign.UserCache.Load(uname)
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
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		id      = c.QueryParam("id")
	)

	// fetch user space for valid login
	us, ok := sign.UserCache.Load(uname)
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

// @Title upload-formfile
// @Summary upload file action via form file input.
// @Description
// @Tags    File
// @Accept  multipart/form-data
// @Produce json
// @Param   note   formData string false "note for uploading file"
// @Param   group0 formData string false "1st category for uploading file"
// @Param   group1 formData string false "2nd category for uploading file"
// @Param   group2 formData string false "3rd category for uploading file"
// @Param   file   formData file   true  "file path for uploading"
// @Success 200 "OK - return storage path"
// @Failure 400 "Fail - file param is incorrect"
// @Failure 500 "Fail - internal error"
// @Router /api/file/upload-formfile [post]
// @Security ApiKeyAuth
func UploadFormFile(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		// Read form fields
		note   = c.FormValue("note")
		group0 = c.FormValue("group0")
		group1 = c.FormValue("group1")
		group2 = c.FormValue("group2")
	)

	// fetch user space for valid login
	us, ok := sign.UserCache.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [upload] @"+uname)
	}

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	fmt.Println("note ---> ", note)

	path, err := us.(*fm.UserSpace).SaveFormFile(file, note, group0, group1, group2)
	if err != nil {
		lk.Warn("UploadFormFile / SaveFormFile ERR: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// * root    path   	 	 "data/user-space/"
	// * storage path   	 	 "data/user-space/cdutwhu/2022-05/g0/g1/g2/document/github key.1652858188.txt"
	// * this    return 	 	 "2022-05/g0/g1/g2/document/github key.1652858188.txt"
	// * future  access url need "[ip:port]/[uname]/2022-05/g0/g1/g2/document/github key.1652858188.txt"

	parts := strings.Split(path, "/")
	path = strings.Join(parts[3:], "/")
	return c.JSON(http.StatusOK, path)
}

// @Title upload-bodydata
// @Summary upload file action via body content.
// @Description
// @Tags    File
// @Accept  application/octet-stream
// @Produce json
// @Param   fname  query string true  "filename for uploading data from body"
// @Param   note   query string false "note for uploading file"
// @Param   group0 query string false "1st category for uploading file"
// @Param   group1 query string false "2nd category for uploading file"
// @Param   group2 query string false "3rd category for uploading file"
// @Param   data   body  string true  "file data for uploading" Format(binary)
// @Success 200 "OK - return storage path"
// @Failure 400 "Fail - file param is incorrect"
// @Failure 500 "Fail - internal error"
// @Router /api/file/upload-bodydata [post]
// @Security ApiKeyAuth
func UploadBodyData(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		fname   = c.QueryParam("fname")
		note    = c.QueryParam("note")
		group0  = c.QueryParam("group0")
		group1  = c.QueryParam("group1")
		group2  = c.QueryParam("group2")
		dataRdr = c.Request().Body
	)

	// fetch user space for valid login
	us, ok := sign.UserCache.Load(uname)
	if !ok || us == nil {
		return c.String(http.StatusInternalServerError, "login error for [upload] @"+uname)
	}

	if len(fname) == 0 {
		return c.String(http.StatusBadRequest, "file name is empty")
	}
	if dataRdr == nil {
		return c.String(http.StatusBadRequest, "body data is empty")
	}

	path, err := us.(*fm.UserSpace).SaveFile(fname, note, dataRdr, group0, group1, group2)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	parts := strings.Split(path, "/")
	path = strings.Join(parts[3:], "/")
	return c.JSON(http.StatusOK, path)
}
