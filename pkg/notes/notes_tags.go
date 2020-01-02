package notes

import (
	"database/sql"
	"github.com/ksopin/notes/pkg/db"
	"github.com/pkg/errors"
	"strings"
)


type NoteTagLinkRow struct {
	NoteId uint
	TagId uint
}


type NoteTagLinkTable struct {
	conn *sql.DB
}

func NewNoteTagLinkTable(conn *sql.DB) *NoteTagLinkTable {
	return &NoteTagLinkTable{
		conn: conn,
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

	rows, err := t.conn.Query(sql, ids...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	list := make([]*NoteTagLinkRow, 0)
	for rows.Next() {
		e := &NoteTagLinkRow{}
		if err = rows.Scan(&e.NoteId, &e.TagId); err != nil {
			return nil, errors.WithStack(err)
		}
		list = append(list, e)
	}
	return list, nil
}


func (t *NoteTagLinkTable) CountByTagId(e db.Queryer, tagId uint) (int, error) {
	cnt := 0
	err := e.QueryRow("SELECT COUNT(note_id) as cnt FROM notes_tags WHERE tag_id = ?", tagId).Scan(&cnt)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return cnt, nil
}

func (t *NoteTagLinkTable) DeleteByTagId(tagId uint) error {
	res, err := t.conn.Exec("DELETE FROM notes_tags WHERE tag_id = ?", tagId)
	if err == nil {
		_, err = res.RowsAffected()
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func (t *NoteTagLinkTable) DeleteByNoteId(noteId uint) error {
	return t.DeleteByNoteIdTx(t.conn, noteId)
}

func (t *NoteTagLinkTable) DeleteByNoteIdTx(e db.Execer, noteId uint) error {
	res, err := e.Exec("DELETE FROM notes_tags WHERE note_id = ?", noteId)
	if err == nil {
		_, err = res.RowsAffected()
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (t *NoteTagLinkTable) InsertTx(e db.Execer, r *NoteTagLinkRow) error {
	_, err := e.Exec("INSERT INTO notes_tags(note_id, tag_id) VALUES(?, ?)", r.NoteId, r.TagId)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
