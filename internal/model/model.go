package model

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID      int64
	Content string
	Author  string
}

type LogFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
