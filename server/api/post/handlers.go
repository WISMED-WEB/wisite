package post

import (
	"net/http"

	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'post.go' *** //

// @Title post template
// @Summary get post template for dev reference.
// @Description
// @Tags    post
// @Accept  json
// @Produce json
// @Success 200 "OK - upload successfully"
// @Router /api/post/template [get]
// @Security ApiKeyAuth
func GetTemplate(c echo.Context) error {
	return c.JSON(http.StatusOK, PostMeta{
		Title: "",
		Content: []struct {
			Text       string "json:\"text\""
			Attachment string "json:\"attachment\""
		}{
			{
				Text:       "",
				Attachment: "",
			},
		},
		Conclusion: "",
	})
}

// @Title upload a post
// @Summary upload a post by filling a post template.
// @Description
// @Tags    post
// @Accept  json
// @Produce json
// @Success 200 "OK - upload successfully"
// @Failure 400 "Fail - incorrect post-meta data format"
// @Failure 500 "Fail - internal error"
// @Router /api/post/upload [post]
// @Security ApiKeyAuth
func Upload(c echo.Context) error {

	meta := new(PostMeta)
	if err := c.Bind(meta); err != nil {
		return c.String(http.StatusBadRequest, "incorrect PostMeta data format: "+err.Error())
	}

	lk.Log("%v", meta)

	panic("TODO: link meta attachment & file")

	return c.String(http.StatusOK, "post uploaded successfully")
}

func GetPosts(c echo.Context) error {

	panic("TODO:")
	return nil
}
