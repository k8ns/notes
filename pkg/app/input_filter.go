package app

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

    . "github.com/ksopin/notes/pkg/notes"
)

type InputErr map[string]error

func (e *InputErr) Error() string {
	b := new(bytes.Buffer)
	for key, value := range *e {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

var (
    notesInputFilter *NoteInputFilter
    notesInputFilterOnce sync.Once
)

type NoteInputFilter struct {
    filters []func(n *Note)
    validators []func(n *Note) (string, error)
}

func (i *NoteInputFilter) IsValid(n *Note) (bool, error) {
    for _, f := range i.filters {
        f(n)
    }

    errs := make(InputErr)
    for _, v := range i.validators {
        if field, err := v(n); err != nil {
            errs[field] = err
        }
    }

    if len(errs) != 0 {
        return false, &errs
    }
    return true, nil
}

func NewNoteInputFilter() *NoteInputFilter {
    return &NoteInputFilter{
        filters: []func(n *Note){uniqueTagsFilter},
        validators: []func(n *Note) (string, error){cannotBeEmpty, mustHaveTags},
    }
}

func GetNoteInputFilter() *NoteInputFilter {
    notesInputFilterOnce.Do(func(){
        notesInputFilter = NewNoteInputFilter()
    })
    return notesInputFilter
}

func uniqueTagsFilter(n *Note) {
    tags := make([]*Tag, 0)
    m := make(map[string]bool)

    for _, t := range n.Tags {
        if _, ok := m[t.Name]; !ok {
            m[t.Name] = true
            tags = append(tags, t)
        }
    }

    n.Tags = tags
}

func mustHaveTags(n *Note) (string, error) {
    if len(n.Tags) > 0 {
        return "", nil
    }
    return "tags", errors.New("note must have at least one tag")
}

func cannotBeEmpty(n *Note) (string, error) {
    if n.Body == "" {
        return "body", errors.New("note must not have empty body")
    }

    return "", nil
}
