package client

import (
	"fmt"
	"sync"

	lk "github.com/digisan/logkit"
)

var (
	mLayout = &sync.Map{}
)

type Area struct {
	Width  int `json:"width" form:"width" query:"width"`
	Height int `json:"height" form:"height" query:"height"`
}

func (a Area) String() string {
	return fmt.Sprintf("{ width: %04d, height: %04d }", a.Width, a.Height)
}

type Layout struct {
	Viewport Area `json:"viewport"` // whole page area
	Header   Area `json:"header"`   // header area for title etc.
	Menu     Area `json:"menubar"`  // left area for menu
	Content  Area `json:"content"`  // right area for content
	Footer   Area `json:"footer"`   // bottom area for extra message
}

func (lo Layout) String() string {
	return fmt.Sprintf(
		"%s%-12s %v    %s%-12s %v    %s%-12s %v    %s%-12s %v    %s%-12s %v",
		lk.LF, "Viewport:", lo.Viewport,
		lk.LF, "Header:", lo.Header,
		lk.LF, "Menu:", lo.Menu,
		lk.LF, "Content:", lo.Content,
		lk.LF, "Footer:", lo.Footer,
	)
}

func newLayout(area *Area) *Layout {
	lo := new(Layout)

	lo.Viewport = *area

	hHeaderProp := 0.15
	lo.Header.Width = lo.Viewport.Width
	lo.Header.Height = int(float64(lo.Viewport.Height) * hHeaderProp)

	hFooterProp := 0.02
	lo.Footer.Width = lo.Viewport.Width
	lo.Footer.Height = int(float64(lo.Viewport.Height) * hFooterProp)

	wMenuProp := 0.2
	lo.Menu.Width = int(float64(lo.Viewport.Width) * wMenuProp)
	lo.Menu.Height = lo.Viewport.Height

	wContentProp := 1.0 - wMenuProp
	hContentProp := 1.0 - hHeaderProp - hFooterProp
	lo.Content.Width = int(float64(lo.Viewport.Width) * wContentProp)
	lo.Content.Height = int(float64(lo.Viewport.Height) * hContentProp)

	return lo
}

func GetLayout(uname string) *Layout {
	if lo, ok := mLayout.Load(uname); ok {
		return lo.(*Layout)
	}
	return nil
}

func AddLayout(uname string, lo *Layout) {
	mLayout.Store(uname, lo)
}

////////////////////////////////////////////

func (lo *Layout) PostWidth() int {
	return int(float64(lo.Content.Width) * 0.9)
}

func (lo *Layout) PostTitleHeight() int {
	return int(float64(lo.Content.Height) * 0.06)
}

func (lo *Layout) PostContentHeight() int {
	return lo.PostTitleHeight() * 9
}
