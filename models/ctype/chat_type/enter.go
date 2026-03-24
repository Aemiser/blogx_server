package chat_type

type ContentMsg struct {
	Content string
}

type ImagetMsg struct {
	Src string
}

type MarkdownMsg struct {
	Content string
}
type ChatMsg struct {
	ContentMsg  *ContentMsg  `json:"contentMsg,omitempty"`
	ImagetMsg   *ImagetMsg   `json:"imagetMsg,omitempty"`
	MarkdownMsg *MarkdownMsg `json:"markdownMsg,omitempty"`
}
