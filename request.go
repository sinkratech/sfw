package sfw

// Request represent the base request in the application
type Request[T any] struct {
	Lang string `header:"Accept-Language"`
	Body T
}

// EmptyRequestBody is an alias for *struct{}
type EmptyRequestBody = struct {
	Lang string `header:"Accept-Language"`
}
