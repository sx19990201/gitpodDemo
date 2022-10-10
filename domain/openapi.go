package domain

type Get struct {
	Description string   `json:"description"`
	Consumes    []string `json:"consumes"`
	Tags        []string `json:"tags"`
}

type Put struct {
	Description string   `json:"description"`
	Consumes    []string `json:"consumes"`
	Tags        []string `json:"tags"`
}

type Post struct {
	Description string   `json:"description"`
	Consumes    []string `json:"consumes"`
	Tags        []string `json:"tags"`
}

type Delete struct {
	Description string   `json:"description"`
	Consumes    []string `json:"consumes"`
	Tags        []string `json:"tags"`
}
