package post

import (
	"fmt"
	"strings"
)

type Post struct {
	Category string        `json:"category"`
	Topic    string        `json:"topic"`
	Keywords string        `json:"keywords"`
	Content  []Paragraph   `json:"content"`
	VFX      VisualEffects `json:"vfx"`
}

type Paragraph struct {
	Text     string        `json:"text"`
	Atch     Attachment    `json:"attachment"`
	VfxMedia VisualEffects `json:"vfx-media"`
	// VfxText  VisualEffects `json:"vfx-text"`
}

type Attachment struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size string `json:"size"` // real "width,height"
}

type VisualEffects struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (m Post) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\nCategory: %v\n", m.Category))
	sb.WriteString(fmt.Sprintf("Topic: %v\n", m.Topic))
	sb.WriteString(fmt.Sprintf("Keywords: %v\n", m.Keywords))
	sb.WriteString(fmt.Sprintf("VFX.Height: %v\n", m.VFX.Height))
	for _, cont := range m.Content {
		sb.WriteString(fmt.Sprintf("	Text: %v\n", cont.Text))
		sb.WriteString(fmt.Sprintf("	Attachment.Path: %v\n", cont.Atch.Path))
		sb.WriteString(fmt.Sprintf("	Attachment.Type: %v\n", cont.Atch.Type))
		sb.WriteString(fmt.Sprintf("	Attachment.Size: %v\n", cont.Atch.Size))
		sb.WriteString(fmt.Sprintf("	VFX Media.Width: %v\n", cont.VfxMedia.Width))
		sb.WriteString(fmt.Sprintf("	VFX Media.Height: %v\n", cont.VfxMedia.Height))
	}
	return sb.String()
}
