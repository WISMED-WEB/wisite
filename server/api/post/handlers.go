package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	em "github.com/digisan/event-mgr"
	. "github.com/digisan/go-generics/v2"
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
func Template(c echo.Context) error {
	return c.JSON(http.StatusOK, Post{
		Category: "post category",
		Topic:    "post topic",
		Content: []struct {
			Text string "json:\"text\""
			Path string "json:\"path\""
		}{
			{
				Text: "some words for this attach",
				Path: "attached stuff path, which should have been given from 'file upload'",
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
		return c.String(http.StatusBadRequest, fmt.Sprintf("'%s' is invalid storage at server", filepath.Base(epath)))
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

	return c.JSON(http.StatusOK, evt)
}

// @Title get a batch of Post id group
// @Summary get a batch of Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   fetchby query string true "time or count"
// @Param   value   query string true "recent [value] minutes for time OR most recent [value] count"
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - incorrect query param type"
// @Failure 500 "Fail - internal error"
// @Router /api/post/ids [get]
// @Security ApiKeyAuth
func IdBatch(c echo.Context) error {

	// userTkn := c.Get("user").(*jwt.Token)
	// claims := userTkn.Claims.(*usr.UserClaims)

	var (
		// uname   = claims.UName
		fetchby = c.QueryParam("fetchby")
		value   = c.QueryParam("value")
		ids     = []string{}
	)

	if fetchby = strings.ToLower(fetchby); NotIn(fetchby, "time", "count") {
		return c.String(http.StatusBadRequest, "'fetchby' must be one of [time, count]")
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "'value' must be a valid number for time(minutes) or count")
	}

	switch fetchby {
	case "time":
		ids, err = em.FetchEvtIDsByTm(edb, value+"m", "DESC")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	case "count":
		ids, err = em.FetchEvtIDsByCnt(edb, int(n), "", "")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, ids)
}

// @Title get one Post content
// @Summary get one Post content.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id   query string true "Post ID for its content"
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - incorrect query param id"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/one [get]
// @Security ApiKeyAuth
func GetOne(c echo.Context) error {

	// userTkn := c.Get("user").(*jwt.Token)
	// claims := userTkn.Claims.(*usr.UserClaims)

	var (
		// uname = claims.UName
		id = c.QueryParam("id")
	)

	if len(id) == 0 {
		c.String(http.StatusBadRequest, "'id' is invalid (cannot be empty)")
	}

	content, err := edb.GetEvt(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	if content == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Post not found @%s", id))
	}

	return c.JSON(http.StatusOK, content)
}
