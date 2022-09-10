package client

import "fmt"

var (
	mLayout = make(map[string]*Layout)
)

type Area struct {
	Width  int `json:"width" form:"width" query:"width"`
	Height int `json:"height" form:"height" query:"height"`
}

func (a Area) String() string {
	return fmt.Sprintf("width: %d,  height: %d", a.Width, a.Height)
}

type Layout struct {
	Viewport Area `json:"viewport"`
	Header   Area `json:"header"`
	Menu     Area `json:"menubar"`
	Content  Area `json:"content"`
}

func (lo Layout) String() string {
	return fmt.Sprintf(
		"\nViewport: %v\nHeader: %v\nMenu: %v\nContent: %v\n",
		lo.Viewport, lo.Header, lo.Menu, lo.Content,
	)
}

func newLayout(area *Area) *Layout {
	lo := new(Layout)

	lo.Viewport = *area

	hHeaderProp := 0.15
	lo.Header.Width = lo.Viewport.Width
	lo.Header.Height = int(float64(lo.Viewport.Height) * hHeaderProp)

	wMenuProp := 0.2
	lo.Menu.Width = int(float64(lo.Viewport.Width) * wMenuProp)
	lo.Menu.Height = lo.Viewport.Height

	wContentProp := 1.0 - wMenuProp
	hContentProp := 1.0 - hHeaderProp
	lo.Content.Width = int(float64(lo.Viewport.Width) * wContentProp)
	lo.Content.Height = int(float64(lo.Viewport.Height) * hContentProp)

	return lo
}

func GetLayout(uname string) *Layout {
	return mLayout[uname]
}

func AddLayout(uname string, lo *Layout) {
	mLayout[uname] = lo
}

func (lo *Layout) HeaderWidth() int {
	return lo.Header.Width
}

func (lo *Layout) HeaderHeight() int {
	return lo.Header.Height
}

func (lo *Layout) MenuWidth() int {
	return lo.Menu.Width
}

func (lo *Layout) MenuHeight() int {
	return lo.Menu.Height
}

func (lo *Layout) ContentWidth() int {
	return lo.Content.Width
}

func (lo *Layout) ContentHeight() int {
	return lo.Content.Height
}

func (lo *Layout) PostWidth() int {
	return int(float64(lo.ContentWidth()) * 0.9)
}

func (lo *Layout) PostTitleHeight() int {
	return int(float64(lo.ContentHeight()) * 0.06)
}

func (lo *Layout) PostContentHeight() int {
	return lo.PostTitleHeight() * 9
}
