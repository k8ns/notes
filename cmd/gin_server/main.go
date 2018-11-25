package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notes/pkg/app"
	"notes/pkg/notes"

	"strconv"
)

func main() {
	r := gin.Default()

	r.Use(addHeaders)
	InitRoutes(r)
	r.Run(":80")
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}


func InitRoutes(engine *gin.Engine) {
	engine.GET("/tags", tagsList) // 200 || 404 || 500
	engine.GET("/notes", notesList) // 200 || 404 || 500
	engine.POST("/notes", createNote) // 201 H:Location || 400 || 500
	engine.OPTIONS("/notes", ok)

	engine.GET("/notes/:id", getNote)
	engine.PUT("/notes/:id", updateNote) // (200 || 204)  || 404 || 409 || 500
	engine.DELETE("/notes/:id", deleteNote) // 204 || 404 || 405 H:Allow: GET || 503
	engine.OPTIONS("/notes/:id", ok)
}

func ok(c *gin.Context) {
	c.Status(http.StatusOK)
}

func tagsList(c *gin.Context) {
	list, err := app.GetNotesManager().GetTags()
	if err != nil {
		writeErrResponse(c, err, http.StatusBadRequest)
		return
	}
	writeOkResponse(c, list, http.StatusOK)
}

func notesList(c *gin.Context) {
	lastId, _ := strconv.Atoi(c.Query("last_id"))

	tagIds := make([]uint, 0, len(c.QueryArray("tag")))
	if tags, ok := c.GetQueryArray("tag"); ok {
		for _, tagId := range tags {
			tid, err := strconv.Atoi(tagId)
			if err != nil {
				writeErrResponse(c, err, http.StatusBadRequest)
				return
			}
			tagIds = append(tagIds, uint(tid))
		}
	}

	list, err := app.GetNotesManager().GetNotes(uint(lastId), tagIds)
	if err != nil {
		writeErrResponse(c, err, http.StatusInternalServerError)
		return
	}
	writeOkResponse(c, list, http.StatusOK)

}

func getNote(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id"))
	note, err := app.GetNotesManager().GetNote(uint(paramId))
	if err != nil {
		switch err.(type) {
		case app.NotFoundErr:
			writeErrResponse(c, err, http.StatusNotFound)
		default:
			writeErrResponse(c, err, http.StatusInternalServerError)
		}
		return
	}
	writeOkResponse(c, note, http.StatusOK)
}

func createNote(c *gin.Context) {
	note := &notes.Note{}
	err := c.BindJSON(note)
	if err != nil {
		writeErrResponse(c, err, http.StatusBadRequest)
		return
	}

	err = app.GetNotesManager().Save(note)
	if err != nil {
		switch err.(type) {
		case *app.InputErr:
			inputErrs := err.(*app.InputErr)
			errs := map[string]error(*inputErrs)
			writeMapErrResponse(c, errs, http.StatusConflict)
		default:
			writeErrResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	c.Header("Location", c.Request.URL.Host+"/notes/"+strconv.Itoa(int(note.Id)))
	writeOkResponse(c, note, http.StatusCreated)
}

func updateNote(c *gin.Context) {
	paramId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeErrResponse(c, err, http.StatusBadRequest)
		return
	}

	id := uint(paramId)
	note := &notes.Note{}
	err = c.BindJSON(note)
	if err != nil {
		writeErrResponse(c, err, http.StatusBadRequest)
		return
	}

	note.Id = id
	err = app.GetNotesManager().Save(note)
	if err != nil {
		switch err.(type) {
		case app.NotExistsErr:
			writeErrResponse(c, err, http.StatusBadRequest)
		case *app.InputErr:
			inputErrs := err.(*app.InputErr)
			errs := map[string]error(*inputErrs)
			writeMapErrResponse(c, errs, http.StatusConflict)
		default:
			writeErrResponse(c, err, http.StatusInternalServerError)
		}
		return
	}
	writeOkResponse(c, note, http.StatusOK)
}

func deleteNote(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err :=  app.GetNotesManager().Delete(uint(id))
	if err != nil {
		switch err.(type) {
		case app.NotExistsErr:
			writeErrResponse(c, err, http.StatusBadRequest)
		default:
			writeErrResponse(c, err, http.StatusInternalServerError)
		}
		return
	}
	c.Status(http.StatusNoContent)
}



func writeOkResponse(c *gin.Context, data interface{}, status int) {
	c.JSON(status, gin.H{
		"data": data,
	})
}

func writeErrResponse(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}

func writeMapErrResponse(c *gin.Context, errs map[string]error, status int) {
	m := make(map[string]string, len(errs))
	for key, err := range errs {
		m[key] = err.Error()
	}
	c.JSON(status, gin.H{
		"error": m,
	})
}
