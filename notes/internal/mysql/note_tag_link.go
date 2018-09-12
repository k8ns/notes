package mysql

import (
	"database/sql"
	"strings"
)


type NoteTagLinkRow struct {
	NoteId uint
	TagId uint
}


type NoteTagLinkTable struct {
	db *sql.DB
}

func NewNoteTagLinkTable() *NoteTagLinkTable {
	return &NoteTagLinkTable{
		db: GetPersistentDB(),
	}
}


func (t *NoteTagLinkTable) SelectByNoteIds(noteIds ...uint) ([]*NoteTagLinkRow, error) {

	if len(noteIds) == 0 {
		return nil, nil
	}

	sql := "SELECT * FROM notes_tags WHERE note_id IN (?"+
		strings.Repeat(",?", len(noteIds) - 1)+");"

	ids := make([]interface{}, len(noteIds))
	for k, v := range noteIds {
		ids[k] = v
	}

	rows, err := t.db.Query(sql, ids...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*NoteTagLinkRow, 0)
	for rows.Next() {
		e := &NoteTagLinkRow{}
		if err = rows.Scan(&e.NoteId, &e.TagId); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

//func (t *NoteTagLinkTable) SelectByTagsIds(tagIds ...uint) ([]*NoteTagLinkRow, error) {
//	rows, err := t.db.Query(sqlLinkByTagsIds, tagIds)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	list := make([]*NoteTagLinkRow, 0)
//	for rows.Next() {
//		e := &NoteTagLinkRow{}
//		if err = rows.Scan(&e.NoteId, &e.TagId); err != nil {
//			return nil, err
//		}
//		list = append(list, e)
//	}
//	return list, nil
//}

func (t *NoteTagLinkTable) CountByTagId(e Queryer, tagId uint) (int, error) {
	cnt := 0
	err := e.QueryRow("SELECT COUNT(note_id) as cnt FROM notes_tags WHERE tag_id = ?", tagId).Scan(&cnt)
	return cnt, err
}

func (t *NoteTagLinkTable) DeleteByTagId(tagId uint) error {
	res, err := t.db.Exec("DELETE FROM notes_tags WHERE tag_id = ?", tagId)
	if err == nil {
		_, err = res.RowsAffected()
	}
	return err
}

func (t *NoteTagLinkTable) DeleteByNoteId(noteId uint) error {
	return t.DeleteByNoteIdTx(t.db, noteId)
}

func (t *NoteTagLinkTable) DeleteByNoteIdTx(e Execer, noteId uint) error {
	res, err := e.Exec("DELETE FROM notes_tags WHERE note_id = ?", noteId)
	if err == nil {
		_, err = res.RowsAffected()
	}
	return err
}

func (t *NoteTagLinkTable) InsertTx(e Execer, r *NoteTagLinkRow) error {
	_, err := e.Exec("INSERT INTO notes_tags(note_id, tag_id) VALUES(?, ?)", r.NoteId, r.TagId)
	return err
}
