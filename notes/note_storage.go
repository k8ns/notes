package notes

import (
	"sync"
	"notes/notes/internal/mysql"
)

var (
	notesStorage *NotesStorage
	notesStorageOnce sync.Once
)


type NotesStorage struct {
	notesSelect *mysql.NoteSelect
	notesCrud *mysql.Crud
	linksTable *mysql.NoteTagLinkTable
	tagsSelect *mysql.TagSelect
	tagsCrud *mysql.Crud
}

func NewNotesStorage() *NotesStorage {
	return &NotesStorage{
		notesSelect: mysql.NewNoteSelect(),
		notesCrud: mysql.NewNotesCrud(),
		linksTable: mysql.NewNoteTagLinkTable(),
		tagsSelect: mysql.NewTagSelect(),
		tagsCrud: mysql.NewTagsCrud(),
	}
}

func GetNotesStorage() *NotesStorage {
	notesStorageOnce.Do(func(){
		notesStorage = NewNotesStorage()
	})
	return notesStorage
}

func (s *NotesStorage) GetNote(id uint) (*Note, error) {
	noteRow, err := s.notesSelect.GetById(id)
	if err != nil  {
		return nil, err
	}

	links, err := s.linksTable.SelectByNoteIds([]uint{id}...)
	if err != nil {
		return nil, err
	}

	tagsM, err := s.buildTagsMapFromLinks(links)
	if err != nil {
		return nil, err
	}

	note := &Note{Id:noteRow.Id, Body:noteRow.Body}
	for _, link := range links {
		if tag, ok := tagsM[link.TagId]; ok {
			note.AddTag(tag)
		}
	}

	return note, err
}

func (s *NotesStorage) GetNotes(lastId uint, tagIds []uint) ([]*Note, error) {

	noteRows, err := s.notesSelect.GetList(lastId, tagIds)
	if err != nil {
		return nil, err
	}

	l := len(noteRows)
	m := make(map[uint]*Note, l)
	noteIds := make([]uint, 0, l)
	notes := make([]*Note, 0, l)

	for _, r := range noteRows {
		note := &Note{r.Id, r.Body, make([]*Tag, 0)}
		m[r.Id] = note
		noteIds = append(noteIds, r.Id)
		notes = append(notes, note)
	}

	links, err := s.linksTable.SelectByNoteIds(noteIds...)
	if err != nil {
		return nil, err
	}

	tagsM, err := s.buildTagsMapFromLinks(links)
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		if tag, ok := tagsM[link.TagId]; ok {
			m[link.NoteId].AddTag(tag)
		}
	}

	return notes, nil
}

func (s *NotesStorage) buildTagsMapFromLinks(links []*mysql.NoteTagLinkRow) (map[uint]*Tag, error) {

	tagIds := make([]uint, 0, len(links))
	for _, link := range links {
		tagIds = append(tagIds, link.TagId)
	}

	tagsRows, err := s.tagsSelect.GetTagsByIds(tagIds...)
	tagsM := make(map[uint]*Tag)
	if err != nil {
		return nil, err
	}
	for _, r := range tagsRows {
		tag := &Tag{r.Id, r.Name}
		tagsM[tag.Id] = tag
	}
	return tagsM, nil
}


func (s *NotesStorage) Exists(id uint) bool {
	return s.notesSelect.Exists(id)
}

func (s *NotesStorage) Save(note *Note) error {

	db := mysql.GetPersistentDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = s.saveNote(tx, note)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.saveTags(tx, note)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.linkTagsToNote(tx, note)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}


func (s *NotesStorage) Delete(id uint) error {
	note, err := s.GetNote(id)
	if err != nil {
		return err
	}

	db := mysql.GetPersistentDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = s.notesCrud.DeleteTx(tx, id)
	if err == nil {
		err = s.linksTable.DeleteByNoteIdTx(tx, id)
	}

	if err == nil {
        for _, tag := range note.Tags {
            cnt, err := s.linksTable.CountByTagId(tx, tag.Id)
            if err == nil && cnt == 0 {
                err = s.tagsCrud.DeleteTx(tx, tag.Id)
            }

            if err != nil {
                break
            }
        }
    }


	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}


func (s *NotesStorage) saveNote(tx mysql.Execer, note *Note) error {
	noteRow := &mysql.NoteRow{note.Id, note.Body}
	err := s.notesCrud.SaveTx(tx, noteRow)
	if err != nil {
		return err
	}
	note.Id = noteRow.Id
	return nil
}


func (s *NotesStorage) saveTags(tx mysql.Execer, note *Note) (err error) {
	for _, tag := range note.Tags {
		tag.Id, err = s.insertTagIfNotExists(tx, tag)
		if err != nil {
			return err
		}
	}
	return err
}


func (s *NotesStorage) insertTagIfNotExists(tx mysql.Execer, tag *Tag) (uint, error) {
	tagRow, err := s.tagsSelect.GetTagByName(tag.Name)
	if err == nil {
		return tagRow.Id, nil
	}

	tagRow = &mysql.TagRow{Id: tag.Id, Name: tag.Name}
	err = s.tagsCrud.InsertTx(tx, tagRow)
	return tagRow.Id, err
}


func (s *NotesStorage) linkTagsToNote(tx mysql.Execer, note *Note) error {
	err := s.linksTable.DeleteByNoteId(note.Id)
	if err != nil {
		return err
	}

	for _, tag := range note.Tags {
		err = s.linksTable.InsertTx(tx, &mysql.NoteTagLinkRow{note.Id, tag.Id})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NotesStorage) AllTags() ([]*Tag, error) {
	tagsRows, err := s.tagsSelect.GetAll()
	if err != nil {
		return nil, err
	}
	list := make([]*Tag, 0, len(tagsRows))
	for _, r := range tagsRows {
		list = append(list, &Tag{r.Id, r.Name})
	}
	return list, err
}