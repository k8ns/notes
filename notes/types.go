package notes


type Tag struct {
	Id uint `json:"id"`
	Name string `json:"name"`
}

type Note struct {
	Id uint `json:"id"`
	Body string `json:"body"`
	Tags []*Tag `json:"tags"`
}

func (n *Note) AddTag(t *Tag) {
	n.Tags = append(n.Tags, t)
}
