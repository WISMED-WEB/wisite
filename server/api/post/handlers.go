package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	em "github.com/digisan/event-mgr"
	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// https://github.com/swaggo/swag
// https://swagger.io/specification/

// *** after implementing, register with path in 'post.go' *** //

// @Title Post template
// @Summary get Post template for dev reference.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Success 200 "OK - upload successfully"
// @Router /api/post/template [get]
// @Security ApiKeyAuth
func GetTemplate(c echo.Context) error {
	return c.JSON(http.StatusOK, Post{
		Category: "post category",
		Topic:    "post topic",
		Content: []struct {
			Text string "json:\"text\""
			Type string "json:\"type\""
			Path string "json:\"path\""
		}{
			{
				Text: "some words for this attach",
				Type: "attachment type, e.g. image/video/audio/pdf/others",
				Path: "attached stuff path, which is given from 'file upload'",
			},
		},
		Summary: "summarize your topic",
	})
}

// @Title upload a Post
// @Summary upload a Post by filling a Post template.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Success 200 "OK - upload successfully"
// @Failure 400 "Fail - incorrect Post format"
// @Failure 500 "Fail - internal error"
// @Router /api/post/upload [post]
// @Security ApiKeyAuth
func Upload(c echo.Context) error {

	userTkn := c.Get("user").(*jwt.Token)
	claims := userTkn.Claims.(*usr.UserClaims)

	var (
		uname = claims.UName
	)

	P := new(Post)
	if err := c.Bind(P); err != nil {
		return c.String(http.StatusBadRequest, "incorrect Post format: "+err.Error())
	}

	lk.Log("%s -- %v", uname, P)

	// validate each path from P
	//
	paths := []string{}
	for _, item := range P.Content {
		if len(item.Path) > 0 {
			path := filepath.Join("data/user-space", uname, item.Path)
			paths = append(paths, path)
		}
	}
	ok, epath := fd.AllExistAsWhole(paths...)
	if !ok {
		return c.String(http.StatusBadRequest, fmt.Sprintf("'%s' is invalid stored at server", epath))
	}

	// save P as JSON for event
	//
	lk.Log("%s", es.CurrIDs())

	data, err := json.Marshal(P)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	evt := em.NewEvent("", uname, "Post", string(data), edb.SaveEvt)
	if err = es.AddEvent(evt); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Post uploaded successfully")
}

// @Title get a batch of Post id group
// @Summary get a batch of Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - "
// @Failure 500 "Fail - internal error"
// @Router /api/post/ids [get]
// @Security ApiKeyAuth
func GetIdBatch(c echo.Context) error {

	panic("TODO:")

	return nil
}
