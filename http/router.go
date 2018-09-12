package http

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "notes/notes"
    "strconv"
)

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

func tagsList(c *gin.Context) {
	list, err := notes.GetNotesStorage().AllTags()
	writeResponse(c, list, err, http.StatusOK)
}

func notesList(c *gin.Context) {
	lastId, _ := strconv.Atoi(c.Query("last_id"))

    tagIds := make([]uint, 0, len(c.QueryArray("tag")))

	if tags, ok := c.GetQueryArray("tag"); ok {
	    for _, tagId := range tags {
	        tid, err := strconv.Atoi(tagId)
	        if err != nil {
                writeResponse(c, nil, err, http.StatusBadRequest)
            }
            tagIds = append(tagIds, uint(tid))
        }
    }

	list, err := notes.GetNotesStorage().GetNotes(uint(lastId), tagIds)
	writeResponse(c, list, err, http.StatusOK)
}

func getNote(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id"))
	if !notes.GetNotesStorage().Exists(uint(paramId)) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	note, err := notes.GetNotesStorage().GetNote(uint(paramId))
	writeResponse(c, note, err, http.StatusOK)
}

func createNote(c *gin.Context) {
	note := &notes.Note{}
	err := c.BindJSON(note)
	if err == nil {
        errs := notes.GetNoteInputFilter().IsValid(note)
        if len(errs) > 0 {
            writeConflictResponse(c, errs)
            return
        }

        err = notes.GetNotesStorage().Save(note)
	}

	c.Header("Location", c.Request.URL.Host+"/notes/"+strconv.Itoa(int(note.Id)))
	writeResponse(c, note, err, http.StatusCreated)
}

func updateNote(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id"))
	id := uint(paramId)
	if !notes.GetNotesStorage().Exists(id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	note := &notes.Note{}
	err := c.BindJSON(note)
	if err == nil {
		note.Id = id

		errs := notes.GetNoteInputFilter().IsValid(note)
		if len(errs) > 0 {
            writeConflictResponse(c, errs)
            return
        }

        err = notes.GetNotesStorage().Save(note)
	}

	writeResponse(c, note, err, http.StatusOK)
}

func deleteNote(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if ! notes.GetNotesStorage().Exists(uint(id)) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err :=  notes.GetNotesStorage().Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func writeConflictResponse(c *gin.Context, errs map[string]error) {

    out := make(map[string]string)
    for field, err := range errs {
        out[field] = err.Error()
    }

    c.JSON(http.StatusConflict, gin.H{
        "error": out,
    })
}

func writeResponse(c *gin.Context, data interface{}, err error, status int) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(status, gin.H{
		"data": data,
	})
}

func ok(c *gin.Context) {
	c.Status(http.StatusOK)
}
