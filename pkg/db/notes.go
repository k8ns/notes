package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
)

type NoteRow struct {
	Id uint
	UserId uint
	Body string
}

func (e *NoteRow) EmptyCopy() Row {
	return &TagRow{}
}

func (e *NoteRow) GetId() uint {
	return e.Id
}

func (e *NoteRow) SetId(id uint) {
	e.Id = id
}

func (e *NoteRow) InsertArgs() []interface{} {
	return []interface{}{e.Body}
}

func (e *NoteRow) UpdateArgs() []interface{} {
	return []interface{}{e.Body, e.Id}
}


type NotesTable struct {
	notesSelect *NoteSelect
	notesCrud *Crud
}

func NewNotesTable() *NotesTable {
	return &NotesTable{
		notesSelect: NewNoteSelect(),
		notesCrud: NewNotesCrud(),
	}
}

func (t *NotesTable) SaveTx(e Execer, r Row) error {
	return t.notesCrud.SaveTx(e, r)
}

func (t *NotesTable) Exists(userId, id uint) bool {
	return t.notesSelect.Exists(userId, id)
}

func (t *NotesTable) GetById(userId, id uint) (*NoteRow, error) {
	return t.notesSelect.GetById(userId, id)
}

func (t *NotesTable) GetList(userId uint, lastId uint, tagIds []uint) ([]*NoteRow, error) {
	return t.notesSelect.GetList(userId, lastId, tagIds)
}

func (t *NotesTable) DeleteTx(c Execer, id uint) error {
	return t.notesCrud.DeleteTx(c, id)
}

func NewNotesCrud() *Crud {
	return &Crud{
		db:        GetPersistentDB(),
		prototype: &NoteRow{},
		sqlInsert: "INSERT INTO notes(body) VALUES(?)",
		sqlUpdate: "UPDATE notes SET body = ? WHERE id = ?",
		sqlDelete: "DELETE FROM notes WHERE id = ?",
	}
}

type NoteSelect struct {
	db        *sql.DB
}

func NewNoteSelect() *NoteSelect {
	return &NoteSelect{
		GetPersistentDB(),
	}
}

func (t *NoteSelect) GetList(userId uint, lastId uint, tagIds []uint) ([]*NoteRow, error) {
	q, b := t.buildQuery(userId, lastId, tagIds)
	rows, err := t.db.Query(q, b...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*NoteRow, 0)
	for rows.Next() {
		e := &NoteRow{}
		if err := rows.Scan(&e.Id, &e.UserId, &e.Body); err != nil {
			return nil, err
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
	e := &NoteRow{}
	row := t.db.QueryRow("SELECT * FROM notes WHERE id = ? AND user_id = ?", id, userId)
	err := row.Scan(&e.Id, &e.Body)
	return e, err
}

func (t *NoteSelect) Exists(userId, id uint) bool {
	var iid int
	t.db.QueryRow("SELECT id FROM notes WHERE id = ? AND user_id = ?", id, userId).Scan(&iid)
	return iid > 0
}
