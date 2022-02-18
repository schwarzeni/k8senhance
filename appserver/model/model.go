package model

type MenuItem struct {
	ID   string
	Name string
	URL  string
}

type User struct {
	ID   string
	Name string
}

type TemplateNavStatus struct {
	Menu    []MenuItem
	User    User
	CurrURL string
}
