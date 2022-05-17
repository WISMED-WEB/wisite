package post

import (
	"fmt"
	"strings"
)

type Post struct {
	// Category string `json:"category"`
	Topic   string `json:"topic"`
	Content []struct {
		Text       string `json:"text"`
		AttachType string `json:"attachtype"`
		Attachment string `json:"attachment"`
	} `json:"content"`
	Summary string `json:"conclusion"`
}

func (m Post) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\nTitle: %v\n", m.Topic))
	for _, cont := range m.Content {
		sb.WriteString(fmt.Sprintf("	Text: %v\n", cont.Text))
		sb.WriteString(fmt.Sprintf("	Attachment Type: %v\n", cont.AttachType))
		sb.WriteString(fmt.Sprintf("	Attachment: %v\n", cont.Attachment))
	}
	sb.WriteString(fmt.Sprintf("Summary: %v\n", m.Summary))
	return sb.String()
}
