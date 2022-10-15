package post

import (
	"fmt"
	"strings"

	. "github.com/digisan/go-generics/v2"
	nt "github.com/digisan/gotk/net-tool"
	clt "github.com/wismed-web/wisite-api/server/api/client"
)

type Post struct {
	Category string      `json:"category"`
	Topic    string      `json:"topic"`
	Keywords string      `json:"keywords"`
	Content  []Paragraph `json:"content"`
}

type Paragraph struct {
	Text     string     `json:"text"`
	RichText string     `json:"richtext"`
	Atch     Attachment `json:"atch"`
	Ele      Element    `json:"ele"`
}

type Attachment struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size string `json:"size"` // real dimension "width,height"
}

type Element struct {
	Video string `json:"html-video"` // <video>
	Image string `json:"html-image"` // <image>
	Audio string `json:"html-audio"` // <audio>
	HP    string `json:"html-hp"`    // <h1-h6> or <p>
}

func (p Post) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\nCategory: %v\n", p.Category))
	sb.WriteString(fmt.Sprintf("Topic: %v\n", p.Topic))
	sb.WriteString(fmt.Sprintf("Keywords: %v\n", p.Keywords))
	for _, cont := range p.Content {
		sb.WriteString(fmt.Sprintf("\t\tText: %v\n", cont.Text))
		sb.WriteString(fmt.Sprintf("\t\tAttachment.Path: %v\n", cont.Atch.Path))
		sb.WriteString(fmt.Sprintf("\t\tAttachment.Type: %v\n", cont.Atch.Type))
		sb.WriteString(fmt.Sprintf("\t\tAttachment.Size: %v\n", cont.Atch.Size))
	}
	return sb.String()
}

func (p *Post) GenVFX(uname string, remoteMode bool) {

	lo := clt.GetLayout(uname)
	width := lo.PostWidth() / 2
	height := lo.PostContentHeight() / 2
	ip := IF(remoteMode, nt.PublicIP(), "127.0.0.1")
	port := 1323

	for _, para := range p.Content {

		ele := para.Ele

		if len(para.RichText) != 0 {
			ele.HP = para.RichText
		} else if len(para.Text) != 0 {
			ele.HP = fmt.Sprintf(`<p>%s</p>`, para.Text)
		}

		src := fmt.Sprintf(`http://%s:%d%s`, ip, port, para.Atch.Path)

		switch para.Atch.Type {
		case "image":
			ele.Image = fmt.Sprintf(`<img src="%s" alt="" width="%d" height="%d">`, src, width, height)
		case "video":
			ele.Video = fmt.Sprintf(`<video width="%d" height="%d" controls autoplay muted>
				<source src="%s" type="video/mp4">
			</video>`, width, height, src)
		case "":
			ele.Audio = fmt.Sprintf(`<audio controls>
				<source src="%s" type="audio/ogg">
			</audio>`, src)
		default:

		}
	}
}
