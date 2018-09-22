package mysql

import (
    "bytes"
    "database/sql"
    "fmt"
    "strconv"
)

type NoteSelect struct {
	db        *sql.DB
}

func NewNoteSelect() *NoteSelect {
	return &NoteSelect{
		GetPersistentDB(),
	}
}

func (t *NoteSelect) GetList(lastId uint, tagIds []uint) ([]*NoteRow, error) {
	q, b := t.buildQuery(lastId, tagIds)
	rows, err := t.db.Query(q, b...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*NoteRow, 0)
	for rows.Next() {
		e := &NoteRow{}
		if err := rows.Scan(&e.Id, &e.Body); err != nil {
			return nil, err
		}
		list = append(list, e)
	}

	return list, nil
}

func (t *NoteSelect) buildQuery(lastId uint, tagIds []uint) (string, []interface{}){
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

        if lastId > 0 {
            buf.WriteString(" AND n.id < ?")
            bind = append(bind, lastId)
        }

        buf.WriteString(" ORDER BY id DESC LIMIT 10")

        query = buf.String()
    }
	return query, bind
}

func (t *NoteSelect) GetById(id uint) (*NoteRow, error) {
	e := &NoteRow{}
	row := t.db.QueryRow("SELECT * FROM notes WHERE id = ?", id)
	err := row.Scan(&e.Id, &e.Body)
	return e, err
}

func (t *NoteSelect) Exists(id uint) bool {
	var iid int
	t.db.QueryRow("SELECT id FROM notes WHERE id = ?", id).Scan(&iid)
	return iid > 0
}
