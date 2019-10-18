package app

import (
	"errors"
	"strings"
	"testing"
	// "github.com/stretchr/testyfy/mock"

	. "github.com/ksopin/notes/pkg/notes"
)


func TestValidation(t *testing.T) {

	note1 := &Note{
		Id: 1,
		Body: "(valid-body)",
		Tags: make([]*Tag, 0),
	}

	note2 := &Note{
		Id: 2,
		Body: "(invalid-body)",
		Tags: make([]*Tag, 0),
	}

	filter := &NoteInputFilter{
		filters: []func(n *Note){func(n *Note){
			n.Body = strings.Trim(n.Body, "()")
		}},
		validators: []func(n *Note) (string, error){func(n *Note) (string, error) {
			if n.Body != "valid-body" {
				return "body", errors.New("test")
			}
			return "", nil
		}},
	}

	test1, _ := filter.IsValid(note1)
	if !test1 {
		t.Error("note 1 is invalid, must be valid")
	}

	if note1.Body != "valid-body" {
		t.Error("filter didn't work")
	}

	test2, _ := filter.IsValid(note2)
	if test2 {
		t.Error("note 2 is valid, must be invalid")
	}

	if note2.Body != "invalid-body" {
		t.Error("filter didn't work")
	}
}
