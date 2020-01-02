package app

import (
	"context"
	"errors"
	"github.com/ksopin/notes/pkg/auth"
	"github.com/ksopin/notes/pkg/db"
	"github.com/ksopin/notes/pkg/notes"
	"sync"
)


type NotFoundErr error
type NotExistsErr error

var (
	application *App
	appOnce sync.Once
)

type App struct {
	storage *notes.Storage
	inputFilter *NoteInputFilter
	service *auth.Service
}

func InitApp(cfg *Config) {
	appOnce.Do(func() {
		application = NewApp(cfg)
	})
}



func NewApp(cfg *Config) *App {

	// all the magic should be here
	// init db

	conn := db.GetPersistentDB()

	return &App{
		storage:     notes.NewStorage(conn),
		inputFilter: NewNoteInputFilter(),
		service: auth.New(db.GetPersistentDB(), cfg.Auth.KeyPath),
	}
}

func Get() *App {
	if application == nil {
		panic("application has not initialized")
	}
	return application
}

func (m *App) GetTags(ctx context.Context) ([]*notes.Tag, error) {
	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	return m.storage.GetTags(u.Id)
}

func (m *App) GetNotes(ctx context.Context, lastId uint, tagIds []uint) ([]*notes.Note, error) {
	u, err := auth.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	return m.storage.GetNotes(u.Id, lastId, tagIds)
}

func (m *App) GetNote(ctx context.Context, id uint) (*notes.Note, error) {

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



func (m *App) Save(ctx context.Context, n *notes.Note) error {

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

func (m *App) Delete(ctx context.Context, id uint) error {

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

func (m *App) VerifySignature(ctx context.Context, token string) (*auth.User, error) {
	return m.service.VerifySignature(ctx, token)
}

func (a *App) Auth(ctx context.Context, creds *auth.Credentials) (string, error) {
	return a.service.Auth(ctx, creds)
}