package platform

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateAccountOpts struct {
	Name string `json:"name"`
}
