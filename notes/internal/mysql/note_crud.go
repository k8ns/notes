package mysql


type NoteRow struct {
	Id uint
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


func NewNotesCrud() *Crud {
	return &Crud{
		db: GetPersistentDB(),
		prototype: &NoteRow{},
		sqlInsert: "INSERT INTO notes(body) VALUES(?)",
		sqlUpdate: "UPDATE notes SET body = ? WHERE id = ?",
		sqlDelete: "DELETE FROM notes WHERE id = ?",
	}
}
