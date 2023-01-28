package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	em "github.com/digisan/event-mgr"
	fm "github.com/digisan/file-mgr"
	"github.com/digisan/file-mgr/fdb"
	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt/v4"
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
		Category: "Post category",
		Topic:    "Post topic",
		Keywords: "keywords for this Post",
		Content: []Paragraph{
			{
				Text:     "some words for this paragraph",
				RichText: "html format for text",
				Atch: Attachment{
					Path: "attachment path, which should have been given from 'file upload'",
					Type: "attachment file type, e.g. image, video, audio, etc",
					Size: "if attachment file type is image or video, return its size as format 'width,height'",
				},
			},
		},
	})
}

// @Title upload a Post
// @Summary upload a Post by filling a Post template.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   data body string true "filled Post template json file"
// @Param   followee  query string false "followee Post ID (empty when doing a new post)"
// @Success 200 "OK - upload successfully"
// @Failure 400 "Fail - incorrect Post format"
// @Failure 500 "Fail - internal error"
// @Router /api/post/upload [post]
// @Security ApiKeyAuth
func Upload(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		flwee   = c.QueryParam("followee")
	)

	P := new(Post)
	if err := c.Bind(P); err != nil {
		lk.Warn("incorrect Uploaded Post format: " + err.Error())
		return c.String(http.StatusBadRequest, "incorrect Post format: "+err.Error())
	}
	lk.Log("Uploading ---> [%s] --- %v", uname, P)

	// get rid of empty paragraph
	//
	FilterFast(&P.Content, func(i int, e Paragraph) bool {
		return len(e.Text) > 0 || len(e.Atch.Path) > 0
	})

	// validate each path from P
	//
	paths := []string{}
	for _, item := range P.Content {
		if len(item.Atch.Path) > 0 {
			path := filepath.Join("data/user-space", uname, item.Atch.Path)
			paths = append(paths, path)
		}
	}
	if ok, epath := fd.AllExistAsWhole(paths...); !ok {
		return c.String(http.StatusBadRequest, fmt.Sprintf("'%s' is invalid storage at server", filepath.Base(epath)))
	}

	// set P Category
	//
	switch {
	case len(flwee) > 0:
		P.Category = "comment"
	default:
		P.Category = "post"
	}

	// save P as JSON for event
	//
	data, err := json.Marshal(P)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	evt := em.NewEvent("", uname, "Post", string(data), flwee)
	if len(evt.ID) > 0 {
		if err = em.AddEvent(evt); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// DEBUG
		// gio.MustAppendFile("./debug.txt", []byte(evt.ID), true)

		// FOLLOWING...
		if len(flwee) > 0 {
			ef, err := em.FetchFollow(flwee)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			if ef == nil {
				if ef, err = em.NewEventFollow(flwee, true); err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
			}
			if err := ef.AddFollower(evt.ID); err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
	}

	// lk.Log("---> %s", em.CurrIDs())

	return c.JSON(http.StatusOK, evt.ID)
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
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/ids [get]
// @Security ApiKeyAuth
func IdBatch(c echo.Context) error {
	var (
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
		ids, err = em.FetchEvtIDsByTm(value + "m")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	case "count":
		ids, err = em.FetchEvtIDsByCnt(int(n), "")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	// lk.Log("IdBatch ---> %d : %v", len(ids), ids)

	// if len(ids) == 0 {
	// 	return c.JSON(http.StatusNotFound, ids)
	// }
	return c.JSON(http.StatusOK, ids)
}

// @Title get all Post id group
// @Summary get all Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Success 200 "OK - get successfully"
// @Failure 404 "Fail - empty event ids"
// @Failure 500 "Fail - internal error"
// @Router /api/post/ids-all [get]
// @Security ApiKeyAuth
func IdAll(c echo.Context) error {

	ids, err := em.FetchEvtIDs(nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// lk.Log("IdAll ---> %d : %v", len(ids), ids)

	// if len(ids) == 0 {
	// 	return c.JSON(http.StatusNotFound, ids)
	// }
	return c.JSON(http.StatusOK, ids)
}

// @Title get one Post content
// @Summary get one Post content.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id     query string  true "Post ID for its content"
// @Param   remote query boolean true "remote ip for media src?"
// @Success 200 "OK - get Post event successfully"
// @Failure 400 "Fail - incorrect query param id"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/one [get]
// @Security ApiKeyAuth
func GetOne(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		id      = c.QueryParam("id")
		remote  = c.QueryParam("remote")
	)

	lk.Log("Into GetOne, event id is %v", id)

	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "'id' is invalid (cannot be empty)")
	}

	event, err := em.FetchEvent(true, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if event == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Post not found @%s", id))
	}
	if len(event.RawJSON) == 0 {
		return c.JSON(http.StatusOK, fmt.Sprintf("Post has no content @%s", id))
	}

	////////////////////////////////////

	// set up event content, i.e. Post
	P := &Post{}
	if err := json.Unmarshal([]byte(event.RawJSON), P); err != nil {
		lk.Warn("Unmarshal Post Error, event is %v", event)
		return c.String(http.StatusInternalServerError, "convert RawJSON to [Post] Unmarshal error")
	}

	for i, p := range P.Content {

		// originally, path start with yyyy-mm
		path := p.Atch.Path

		// 1) update path for remote access
		P.Content[i].Atch.Path = filepath.Join(event.Owner, path)

		// 2) update type
		fpath := filepath.Join("data", "user-space", event.Owner, path)
		ftype := fdb.GetFileType(fpath)
		P.Content[i].Atch.Type = ftype

		sz := ""
		switch ftype {
		case "image":
			sz, err = fm.GetImageSize(fpath)
		case "video":
			sz, err = fm.GetVideoSize(fpath)
		default:
			// no need to get area size
		}
		if In(ftype, "image", "video") && err != nil {
			lk.Warn("get media area size error %v @ %s @ %s", err, ftype, fpath)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// 3) update area size
		P.Content[i].Atch.Size = sz
	}

	rmt, err := strconv.ParseBool(remote)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	P.GenVFX(uname, rmt)

	////////////////////////////////////

	PData, err := json.Marshal(P)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	event.RawJSON = string(PData)

	lk.Log("-->\n %v", event)

	return c.JSON(http.StatusOK, event)
}

// @Title delete one Post content
// @Summary delete one Post content.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id   query string true "Post ID for deleting"
// @Success 200 "OK - delete successfully"
// @Failure 400 "Fail - incorrect query param id"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/del/one [delete]
// @Security ApiKeyAuth
func DelOne(c echo.Context) error {
	var (
		id = c.QueryParam("id")
	)
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "'id' is invalid (cannot be empty)")
	}

	n, err := em.DelEvent(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, IF(n == 1, fmt.Sprintf("<%s> is deleted", id), fmt.Sprintf("<%s> is not existing, nothing to delete", id)))
}

// @Title erase one Post content
// @Summary erase one Post content permanently.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id   query string true "Post ID for erasing"
// @Success 200 "OK - erase successfully"
// @Failure 400 "Fail - incorrect query param id"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/post/erase/one [delete]
// @Security ApiKeyAuth
func EraseOne(c echo.Context) error {
	var (
		id = c.QueryParam("id")
	)
	if len(id) == 0 {
		return c.String(http.StatusBadRequest, "'id' is invalid (cannot be empty)")
	}

	n, err := em.EraseEvents(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, IF(n == 1, fmt.Sprintf("<%s> is erased permanently", id), fmt.Sprintf("<%s> is not existing, nothing to erase", id)))
}

// @Title get own Post id group in a specific period
// @Summary get own Post id group in one specific time period.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   period query string false "time period for query, format is 'yyyymm', e.g. '202206'. if missing, current yyyymm applies"
// @Success 200 "OK - get successfully"
// @Failure 400 "Fail - incorrect query param type"
// @Failure 404 "Fail - empty event ids"
// @Failure 500 "Fail - internal error"
// @Router /api/post/own/ids [get]
// @Security ApiKeyAuth
func OwnPosts(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		period  = c.QueryParam("period")
	)

	if len(period) == 0 {
		period = time.Now().Format("200601")
	}
	if _, err := time.Parse("200601", period); err != nil {
		return c.String(http.StatusBadRequest, "'period' format must be 'yyyymm', e.g. '202206'")
	}

	lk.Log("%s -- %s", uname, period)

	ids, err := em.FetchOwn(uname, period)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// if len(ids) == 0 {
	// 	return c.JSON(http.StatusNotFound, ids)
	// }
	return c.JSON(http.StatusOK, ids)
}

// @Title toggle a bookmark for a post
// @Summary add or remove a personal bookmark for a post.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id path string true "Post ID (event id) for toggling a bookmark"
// @Success 200 "OK - toggled bookmark successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/post/bookmark/{id} [patch]
// @Security ApiKeyAuth
func Bookmark(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		id      = c.Param("id")
	)
	bm, err := em.NewBookmark(uname, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	has, err := bm.ToggleEvent(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, has)
}

// @Title get current user's bookmark status for a post
// @Summary get current login user's bookmark status for a post.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id path string true "Post ID (event id) for checking bookmark status"
// @Success 200 "OK - get bookmark status successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/post/bookmark/status/{id} [get]
// @Security ApiKeyAuth
func BookmarkStatus(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		id      = c.Param("id")
	)
	bm, err := em.NewBookmark(uname, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bm.HasEvent(id))
}

// @Title get bookmarked Posts
// @Summary get all bookmarked Post ids.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   order query string false "order[desc asc] to get Post ids ordered by event time"
// @Success 200 "OK - get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/post/bookmark/bookmarked [get]
// @Security ApiKeyAuth
func BookmarkedPosts(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		order   = c.QueryParam("order")
	)
	bm, err := em.NewBookmark(uname, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bm.Bookmarks(order))
}

// @Title get a Post follower-Post ids
// @Summary get a specified Post follower-Post id group.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   followee query string true "followee Post ID"
// @Success 200 "OK - get successfully"
// @Failure 404 "Fail - empty follower ids"
// @Failure 500 "Fail - internal error"
// @Router /api/post/follower/ids [get]
// @Security ApiKeyAuth
func Followers(c echo.Context) error {
	var (
		flwee = c.QueryParam("followee")
	)

	flwers, err := em.Followers(flwee)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// if len(flwers) == 0 {
	// 	return c.JSON(http.StatusNotFound, flwers)
	// }
	return c.JSON(http.StatusOK, flwers)
}

// @Title add or remove a thumbsup for a post
// @Summary add or remove a personal thumbsup for a post.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id path string true "Post ID (event id) for adding or removing thumbs-up"
// @Success 200 "OK - added or removed thumb successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/post/thumbsup/{id} [patch]
// @Security ApiKeyAuth
func ThumbsUp(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		id      = c.Param("id")
	)
	ep, err := em.NewEventParticipate(id, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	has, err := ep.TogglePtp("ThumbsUp", uname)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	ptps, err := ep.Ptps("ThumbsUp")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	lk.Log("---> %v", ptps)

	return c.JSON(http.StatusOK, struct {
		ThumbsUp bool
		Count    int
	}{
		has, len(ptps),
	})
}

// @Title get current user's thumbsup status for a post
// @Summary get current login user's thumbsup status for a post.
// @Description
// @Tags    Post
// @Accept  json
// @Produce json
// @Param   id path string true "Post ID (event id) for checking thumbs-up status"
// @Success 200 "OK - get thumbs-up status successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/post/thumbsup/status/{id} [get]
// @Security ApiKeyAuth
func ThumbsUpStatus(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
		id      = c.Param("id")
	)
	ep, err := em.NewEventParticipate(id, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	has := ep.HasPtp("ThumbsUp", uname)
	ptps, err := ep.Ptps("ThumbsUp")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		ThumbsUp bool
		Count    int
	}{
		has, len(ptps),
	})
}
