package post

import (
	"fmt"
	"strings"
)

type Paragraph struct {
	Text string `json:"text"`
	Path string `json:"path"`
}

type Post struct {
	Category string      `json:"category"`
	Topic    string      `json:"topic"`
	Content  []Paragraph `json:"content"`
	Summary  string      `json:"summary"`
}

func (m Post) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\nCategory: %v\n", m.Category))
	sb.WriteString(fmt.Sprintf("Title: %v\n", m.Topic))
	for _, cont := range m.Content {
		sb.WriteString(fmt.Sprintf("	Text: %v\n", cont.Text))
		sb.WriteString(fmt.Sprintf("	Path: %v\n", cont.Path))
	}
	sb.WriteString(fmt.Sprintf("Summary: %v\n", m.Summary))
	return sb.String()
}
