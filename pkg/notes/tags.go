package notes

import (
	"database/sql"
	"github.com/ksopin/notes/pkg/db"
	"github.com/pkg/errors"
	"strings"
)

type TagRow struct {
	Id uint
	UserId uint
	Name string
}

func (e *TagRow) EmptyCopy() db.Row {
	return &TagRow{}
}

func (e *TagRow) GetId() uint {
	return e.Id
}

func (e *TagRow) SetId(id uint) {
	e.Id = id
}

func (e *TagRow) InsertArgs() []interface{} {
	return []interface{}{e.UserId, e.Name}
}

func (e *TagRow) UpdateArgs() []interface{} {
	return []interface{}{e.Name, e.Id}
}


type TagsTable struct {
	tagSelect *TagSelect
	tagCrud *db.Crud
}

func NewTagsTable(conn *sql.DB) *TagsTable {
	return &TagsTable{
		tagSelect: NewTagSelect(conn),
		tagCrud: NewTagsCrud(conn),
	}
}

func (t *TagsTable) GetTagsByIds(tagIds ...uint) ([]*TagRow, error) {
	return t.tagSelect.GetTagsByIds(tagIds...)
}

func (t *TagsTable) GetAll(userId uint) ([]*TagRow, error) {
	return t.tagSelect.GetAll(userId)
}

func (t *TagsTable) DeleteTx(c db.Execer, id uint) error {
	return t.tagCrud.DeleteTx(c, id)
}

func (t *TagsTable) GetTagByName(name string) (*TagRow, error) {
	return t.tagSelect.GetTagByName(name)
}

func (t *TagsTable) InsertTx(e db.Execer, r db.Row) error {
	return t.tagCrud.InsertTx(e, r)
}

func NewTagsCrud(conn *sql.DB) *db.Crud {
	return db.NewCrud(
		conn,
		&TagRow{},
		"INSERT INTO tags(user_id, name) VALUES(?, ?)",
		"UPDATE tags SET name = ? WHERE id = ?",
		"DELETE FROM tags WHERE id = ?",
	)
}


type TagSelect struct {
	db *sql.DB
}

func NewTagSelect(conn *sql.DB) *TagSelect {
	return &TagSelect{
		db: conn,
	}
}

func (t *TagSelect) GetAll(userId uint) ([]*TagRow, error) {
	return t.rows("SELECT * FROM tags WHERE user_id = ?", userId)
}

func (t *TagSelect) GetTagsByIds(tagIds ...uint) ([]*TagRow, error) {

	if len(tagIds) == 0 {
		return nil, nil
	}

	sql := "SELECT * FROM tags WHERE id IN (?"+
		strings.Repeat(",?", len(tagIds) - 1)+");"

	ids := make([]interface{}, len(tagIds))
	for k, v := range tagIds {
		ids[k] = v
	}

	return t.rows(sql, ids...)
}

func (t *TagSelect) rows(sql string, binds ...interface{}) ([]*TagRow, error) {
	rows, err := t.db.Query(sql, binds...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var list []*TagRow
	for rows.Next() {
		e := &TagRow{}

		if err = rows.Scan(&e.Id, &e.UserId, &e.Name); err != nil {
			return nil, errors.WithStack(err)
		}
		list = append(list, e)
	}
	return list, nil
}

func (t *TagSelect) GetTagByName(name string) (*TagRow, error) {
	r := &TagRow{}
	row := t.db.QueryRow("SELECT * FROM tags WHERE `name` = ?", name)
	err := row.Scan(&r.Id, &r.UserId, &r.Name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r, err
}
