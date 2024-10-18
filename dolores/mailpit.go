package dolores

type Message struct {
	ID string `json:"ID"`
}

type MailPitResponse struct {
	Messages []Message `json:"messages"`
}
