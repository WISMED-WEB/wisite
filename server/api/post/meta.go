package post

import (
	"fmt"
	"strings"
)

type PostMeta struct {
	Title   string `json:"title"`
	Content []struct {
		Text       string `json:"text"`
		Attachment string `json:"attachment"`
	} `json:"content"`
	Conclusion string `json:"conclusion"`
}

func (m PostMeta) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("\nTitle: %v\n", m.Title))
	for _, cont := range m.Content {
		sb.WriteString(fmt.Sprintf("	Text: %v\n", cont.Text))
		sb.WriteString(fmt.Sprintf("	Attachment: %v\n", cont.Attachment))
	}
	sb.WriteString(fmt.Sprintf("Conclusion: %v\n", m.Conclusion))
	return sb.String()
}
