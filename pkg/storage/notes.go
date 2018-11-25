package storage

import (
	. "notes/pkg/notes"
)

type NotesStorage interface {
	GetNote(id uint) (*Note, error)
	GetNotes(lastId uint, tagIds []uint) ([]*Note, error)
	Exists(id uint) bool
	Save(note *Note) error
	Delete(id uint) error
	GetTags() ([]*Tag, error)
}
