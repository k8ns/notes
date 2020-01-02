package notes

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ksopin/notes/pkg/db"
	"github.com/pkg/errors"
	"strconv"
)

type NoteRow struct {
	Id uint
	UserId uint
	Body string
}

func (e *NoteRow) EmptyCopy() db.Row {
	return &TagRow{}
}

func (e *NoteRow) GetId() uint {
	return e.Id
}

func (e *NoteRow) SetId(id uint) {
	e.Id = id
}

func (e *NoteRow) InsertArgs() []interface{} {
	return []interface{}{e.UserId, e.Body}
}

func (e *NoteRow) UpdateArgs() []interface{} {
	return []interface{}{e.Body, e.Id}
}


type NotesTable struct {
	notesSelect *NoteSelect
	notesCrud *db.Crud
}

func NewNotesTable(conn *sql.DB) *NotesTable {
	return &NotesTable{
		notesSelect: NewNoteSelect(),
		notesCrud: NewNotesCrud(conn),
	}
}

func (t *NotesTable) SaveTx(e db.Execer, r db.Row) error {
	return t.notesCrud.SaveTx(e, r)
}

func (t *NotesTable) Exists(userId, id uint) (bool, error) {
	return t.notesSelect.Exists(userId, id)
}

func (t *NotesTable) GetById(userId, id uint) (*NoteRow, error) {
	return t.notesSelect.GetById(userId, id)
}

func (t *NotesTable) GetList(userId uint, lastId uint, tagIds []uint) ([]*NoteRow, error) {
	return t.notesSelect.GetList(userId, lastId, tagIds)
}

func (t *NotesTable) DeleteTx(c db.Execer, id uint) error {
	return t.notesCrud.DeleteTx(c, id)
}

func NewNotesCrud(conn *sql.DB) *db.Crud {
	return db.NewCrud(
		conn,
		&NoteRow{},
		"INSERT INTO notes(user_id, body) VALUES(?, ?)",
		"UPDATE notes SET body = ? WHERE id = ?",
		"DELETE FROM notes WHERE id = ?",
	)
}

type NoteSelect struct {
	conn        *sql.DB
}

func NewNoteSelect() *NoteSelect {
	return &NoteSelect{
		conn: db.GetPersistentDB(),
	}
}

func (t *NoteSelect) GetList(userId uint, lastId uint, tagIds []uint) ([]*NoteRow, error) {
	q, b := t.buildQuery(userId, lastId, tagIds)
	rows, err := t.conn.Query(q, b...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	list := make([]*NoteRow, 0)
	for rows.Next() {
		e := &NoteRow{}
		if err := rows.Scan(&e.Id, &e.UserId, &e.Body); err != nil {
			return nil, errors.WithStack(err)
		}
		list = append(list, e)
	}

	return list, nil
}

func (t *NoteSelect) buildQuery(userId uint, lastId uint, tagIds []uint) (string, []interface{}){
	bind := make([]interface{}, 0)
    var query string

	tids := make([]interface{}, len(tagIds))
	for k, v := range tagIds {
		tids[k] = v
	}

	if len(tids) == 0 {
        query = "SELECT * FROM notes"
        if lastId > 0 {
            bind = append(bind, lastId)
            query += " WHERE id < ?"
        }
        query += " ORDER BY id DESC LIMIT 10"
    } else {

        buf := &bytes.Buffer{}

        bind = append(bind, tids...)
        buf.WriteString("SELECT n.* FROM notes_tags t1")
        where := ""
        for i := range tids {
            if i == 0 {
                continue
            }
            k := i + 1
            buf.WriteString(" LEFT JOIN notes_tags t")
            buf.WriteString(strconv.Itoa(k))
            buf.WriteString(" ON t")
            buf.WriteString(strconv.Itoa(k))
            buf.WriteString(".note_id = t1.note_id")
            where += fmt.Sprintf(" AND t%d.tag_id = ?", k)
        }
        buf.WriteString(" LEFT JOIN notes n ON n.id = t1.note_id")
        buf.WriteString(" WHERE t1.tag_id = ?")
        buf.WriteString(where)

        bind = append(bind, userId)
        buf.WriteString(" AND n.user_id = ?")

        if lastId > 0 {
            buf.WriteString(" AND n.id < ?")
            bind = append(bind, lastId)
        }

        buf.WriteString(" ORDER BY id DESC LIMIT 10")

        query = buf.String()
    }
	return query, bind
}

func (t *NoteSelect) GetById(userId, id uint) (*NoteRow, error) {
	r := &NoteRow{}
	row := t.conn.QueryRow("SELECT * FROM notes WHERE id = ? AND user_id = ?", id, userId)
	err := row.Scan(&r.Id, &r.UserId, &r.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r, nil
}

func (t *NoteSelect) Exists(userId, id uint) (bool, error) {
	var iid int
	err := t.conn.QueryRow("SELECT id FROM notes WHERE id = ? AND user_id = ?", id, userId).Scan(&iid)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return iid > 0, nil
}
