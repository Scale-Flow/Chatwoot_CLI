package clientapi

type Contact struct {
	SourceID    string `json:"source_id"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PubsubToken string `json:"pubsub_token,omitempty"`
}

type CreateContactOpts struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone_number,omitempty"`
}
