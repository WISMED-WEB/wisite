package post

import (
	"net/http"

	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
)

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
		Topic: "what's your topic",
		Content: []struct {
			Text       string "json:\"text\""
			AttachType string "json:\"attachtype\""
			Attachment string "json:\"attachment\""
		}{
			{
				Text:       "say some words for current part",
				AttachType: "what's your attachment type, e.g. image/video/audio/pdf/others",
				Attachment: "something attached for current part",
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

	P := new(Post)
	if err := c.Bind(P); err != nil {
		return c.String(http.StatusBadRequest, "incorrect Post format: "+err.Error())
	}

	lk.Log("%v", P)

	// panic("TODO: link Post attachment & file")

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
