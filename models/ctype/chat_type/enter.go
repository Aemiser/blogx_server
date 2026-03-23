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
	ContentMsg  *ContentMsg  `json:"contentMsg"`
	ImagetMsg   *ImagetMsg   `json:"imagetMsg"`
	MarkdownMsg *MarkdownMsg `json:"markdownMsg"`
}
