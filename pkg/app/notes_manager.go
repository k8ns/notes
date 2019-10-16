package app

import (
	"context"
	"errors"
	"github.com/ksopin/notes/pkg/auth"
	"github.com/ksopin/notes/pkg/notes"
	"sync"
)


type NotFoundErr error
type NotExistsErr error

var (
	manager *NotesManager
	managerOnce sync.Once
)

type NotesManager struct {
	storage *notes.Storage
	inputFilter *NoteInputFilter
}

func NewNotesManager() *NotesManager {
	return &NotesManager{
		storage:     notes.NewStorage(),
		inputFilter: NewNoteInputFilter(),
	}
}

func GetNotesManager() *NotesManager {
	managerOnce.Do(func() {
		manager = NewNotesManager()
	})
	return manager
}

func (m *NotesManager) GetTags(ctx context.Context) ([]*notes.Tag, error) {
	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	return m.storage.GetTags(u.Id)
}

func (m *NotesManager) GetNotes(ctx context.Context, lastId uint, tagIds []uint) ([]*notes.Note, error) {
	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	return m.storage.GetNotes(u.Id, lastId, tagIds)
}

//func (m *NotesManager) Exists(id uint) bool {
//	u, err := auth.GetUser(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	return m.storage.Exists(u.Id, id)
//}

func (m *NotesManager) GetNote(ctx context.Context, id uint) (*notes.Note, error) {

	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := m.storage.Exists(u.Id, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, NotExistsErr(errors.New("not exists"))
	}

	return m.storage.GetNote(u.Id, id)
}

func (m *NotesManager) Save(ctx context.Context, n *notes.Note) error {

	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil
	}

	exists, err := m.storage.Exists(u.Id, n.Id)
	if err != nil {
		return err
	}
	if n.Id > 0 && !exists {
		return NotExistsErr(errors.New("not exists"))
	}

	_, err = m.inputFilter.IsValid(n)
	if err != nil {
		return err
	}

	n.UserId = u.Id

	return m.storage.Save(n)
}

func (m *NotesManager) Delete(ctx context.Context, id uint) error {

	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil
	}

	exists, err := m.storage.Exists(u.Id, id)
	if err != nil {
		return err
	}
	if !exists {
		return NotExistsErr(errors.New("not exists"))
	}

	return m.storage.Delete(u.Id, id)
}
