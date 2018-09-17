package mysql

import (
	"database/sql"
	"strings"
)

type TagRow struct {
	Id uint
	Name string
}

func (e *TagRow) EmptyCopy() Row {
	return &TagRow{}
}

func (e *TagRow) GetId() uint {
	return e.Id
}

func (e *TagRow) SetId(id uint) {
	e.Id = id
}

func (e *TagRow) InsertArgs() []interface{} {
	return []interface{}{e.Name}
}

func (e *TagRow) UpdateArgs() []interface{} {
	return []interface{}{e.Name, e.Id}
}

func NewTagsCrud() *Crud {
	return &Crud{
		db: GetPersistentDB(),
		prototype: &TagRow{},
		sqlInsert: "INSERT INTO tags(name) VALUES(?)",
		sqlUpdate: "UPDATE tags SET name = ? WHERE id = ?",
		sqlDelete: "DELETE FROM tags WHERE id = ?",
	}
}




type TagSelect struct {
	db *sql.DB
}

func NewTagSelect() *TagSelect {
	return &TagSelect{
		db: GetPersistentDB(),
	}
}

func (t *TagSelect) GetAll() ([]*TagRow, error) {
	return t.rows("SELECT * FROM tags")
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
		return nil, err
	}
	defer rows.Close()

	var list []*TagRow
	for rows.Next() {
		e := &TagRow{}

		if err = rows.Scan(&e.Id, &e.Name); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, err
}

func (t *TagSelect) GetTagByName(name string) (*TagRow, error) {
	e := &TagRow{}
	row := t.db.QueryRow("SELECT * FROM tags WHERE `name` = ?", name)
	err := row.Scan(&e.Id, &e.Name)
	return e, err
}
