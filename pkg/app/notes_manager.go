package app

import (
	"errors"
	"notes/pkg/notes"
	"notes/pkg/storage"
	"notes/pkg/storage/mysql"
	"sync"

)


type NotFoundErr error
type NotExistsErr error

var (
	manager *NotesManager
	managerOnce sync.Once
)

type NotesManager struct {
	storage storage.NotesStorage
	inputFilter *NoteInputFilter
}

func NewNotesManager() *NotesManager {
	return &NotesManager{
		storage: mysql.NewStorage(),
		inputFilter: NewNoteInputFilter(),
	}
}

func GetNotesManager() *NotesManager {
	managerOnce.Do(func() {
		manager = NewNotesManager()
	})
	return manager
}

func (m *NotesManager) GetTags() ([]*notes.Tag, error) {
	return m.storage.GetTags()
}

func (m *NotesManager) GetNotes(lastId uint, tagIds []uint) ([]*notes.Note, error) {
	return m.storage.GetNotes(lastId, tagIds)
}

func (m *NotesManager) Exists(id uint) bool {
	return m.storage.Exists(id)
}

func (m *NotesManager) GetNote(id uint) (*notes.Note, error) {
	if !m.storage.Exists(id) {
		return nil, NotFoundErr(errors.New("not found"))
	}

	return m.storage.GetNote(id)
}

func (m *NotesManager) Save(n *notes.Note) error {

	if n.Id > 0 && !m.storage.Exists(n.Id) {
		return NotExistsErr(errors.New("not exists"))
	}

	_, err := m.inputFilter.IsValid(n)
	if err != nil {
		return err
	}

	return m.storage.Save(n)
}

func (m *NotesManager) Delete(id uint) error {

	if !m.storage.Exists(id) {
		return NotExistsErr(errors.New("not exists"))
	}

	return m.storage.Delete(id)
}
