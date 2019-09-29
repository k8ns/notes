package mysql

import (
	"database/sql"
	. "notes/pkg/notes"

	"notes/pkg/storage"
)




type Storage struct {
	notesTable *NotesTable
	linksTable *NoteTagLinkTable
	tagsTable *TagsTable
}

func NewStorage() storage.NotesStorage {
	return &Storage{
		notesTable: NewNotesTable(),
		linksTable: NewNoteTagLinkTable(),
		tagsTable: NewTagsTable(),
	}
}

func (s *Storage) Search(search string) ([]*Note, error) {
	//sql := `SELECT * FROM notes WHERE MATCH (body) AGAINST (?)`

	return nil, nil
}

func (s *Storage) GetNote(id uint) (*Note, error) {
	noteRow, err := s.notesTable.GetById(id)
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

func (s *Storage) GetNotes(lastId uint, tagIds []uint) ([]*Note, error) {

	noteRows, err := s.notesTable.GetList(lastId, tagIds)
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

func (s *Storage) buildTagsMapFromLinks(links []*NoteTagLinkRow) (map[uint]*Tag, error) {

	tagIds := make([]uint, 0, len(links))
	for _, link := range links {
		tagIds = append(tagIds, link.TagId)
	}

	tagsRows, err := s.tagsTable.GetTagsByIds(tagIds...)
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


func (s *Storage) Exists(id uint) bool {
	return s.notesTable.Exists(id)
}

func (s *Storage) Save(note *Note) error {

	db := GetPersistentDB()
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

	oldTags, err := s.getTagsByNoteId(note.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.linkTagsToNote(tx, note)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.deleteUnusedTags(tx, oldTags)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}



func (s *Storage) Delete(id uint) error {
	note, err := s.GetNote(id)
	if err != nil {
		return err
	}

	db := GetPersistentDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = s.notesTable.DeleteTx(tx, id)
	if err == nil {
		err = s.linksTable.DeleteByNoteIdTx(tx, id)
	}

	if err == nil {
		err = s.deleteUnusedTags(tx, note.Tags)
    }

	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func (s *Storage) GetTags() ([]*Tag, error) {
	tagsRows, err := s.tagsTable.GetAll()
	if err != nil {
		return nil, err
	}
	list := make([]*Tag, 0, len(tagsRows))
	for _, r := range tagsRows {
		list = append(list, &Tag{r.Id, r.Name})
	}
	return list, err
}

func (s *Storage) getTagsByNoteId(id uint) ([]*Tag, error) {

	links, err := s.linksTable.SelectByNoteIds(id)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0, len(links))
	for _, l := range links {
		tags = append(tags, &Tag{Id: l.TagId})
	}
	return tags, nil
}

func (s *Storage) deleteUnusedTags(tx *sql.Tx, tags []*Tag) error {
	for _, tag := range tags {
		cnt, err := s.linksTable.CountByTagId(tx, tag.Id)
		if err == nil && cnt == 0 {
			err = s.tagsTable.DeleteTx(tx, tag.Id)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) saveNote(tx Execer, note *Note) error {
	noteRow := &NoteRow{note.Id, note.Body}
	err := s.notesTable.SaveTx(tx, noteRow)
	if err != nil {
		return err
	}
	note.Id = noteRow.Id
	return nil
}


func (s *Storage) saveTags(tx Execer, note *Note) (err error) {
	for _, tag := range note.Tags {
		tag.Id, err = s.insertTagIfNotExists(tx, tag)
		if err != nil {
			return err
		}
	}
	return err
}


func (s *Storage) insertTagIfNotExists(tx Execer, tag *Tag) (uint, error) {
	tagRow, err := s.tagsTable.GetTagByName(tag.Name)
	if err == nil {
		return tagRow.Id, nil
	}

	tagRow = &TagRow{Id: tag.Id, Name: tag.Name}
	err = s.tagsTable.InsertTx(tx, tagRow)
	return tagRow.Id, err
}


func (s *Storage) linkTagsToNote(tx Execer, note *Note) error {
	err := s.linksTable.DeleteByNoteId(note.Id)
	if err != nil {
		return err
	}

	for _, tag := range note.Tags {
		err = s.linksTable.InsertTx(tx, &NoteTagLinkRow{note.Id, tag.Id})
		if err != nil {
			return err
		}
	}
	return nil
}

