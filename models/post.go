package models

type Post struct {
	Id      string
	Title   string
	Content string
}

type Client struct {
	Id        string
	Email     string
	firstName string
	lastName  string
	Surname   string
	birthDate string
	Password  string
}

func NewPost(id, title, content string) *Post {
	return &Post{id, title, content}
}

func NewClient(Id, Email, firstName, lastName, Surname, birthDate, Password string) *Client {
	return &Client{Id, Email, firstName, lastName, Surname, birthDate, Password}
}
