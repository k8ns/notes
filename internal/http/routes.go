package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ksopin/notes/internal/app"
	"github.com/ksopin/notes/pkg/auth"
	"github.com/ksopin/notes/pkg/notes"
	"net/http"
	"strconv"
)

func InitRoutes(engine *gin.Engine) {

	engine.POST("/sign-in", signIn)

	authorized := engine.Group("/")
	authorized.Use(authMiddleware)
	{
		authorized.GET("/user", userInfo) // 200 || 404 || 500
		authorized.GET("/tags", tagsList) // 200 || 404 || 500
		authorized.GET("/notes", notesList) // 200 || 404 || 500
		authorized.POST("/notes", createNote) // 201 H:Location || 400 || 500
		authorized.GET("/notes/:id", getNote)
		authorized.PUT("/notes/:id", updateNote) // (200 || 204)  || 404 || 409 || 500
		authorized.DELETE("/notes/:id", deleteNote) // 204 || 404 || 405 H:Allow: GET || 503
	}
}

func InitWelcome(engine *gin.Engine, project *app.Config) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": map[string]string{
				"projectName": project.ProjectName,
				"code": project.Code,
				"version": project.Version,
				"env": project.Env,
			},
		})
	})
}


func userInfo(c *gin.Context) {
	u, err := auth.GetUser(c)
	if err != nil {
		writeErrResponse(c, err, http.StatusUnauthorized)
		return
	}

	writeOkResponse(c, u, http.StatusOK)
}

func signIn(c *gin.Context) {
	creds := &auth.Credentials{}
	err := c.BindJSON(creds)
	if err != nil {
		writeErrResponse(c, err, http.StatusBadRequest)
		return
	}

	token, err := app.GetAuthService().Auth(c, creds)
	if err != nil {
		writeErrResponse(c, err, http.StatusForbidden)
		return
	}
	writeOkResponse(c, token, http.StatusOK)
}

func tagsList(c *gin.Context) {
	list, err := app.GetNotesManager().GetTags(c)
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

	list, err := app.GetNotesManager().GetNotes(c, uint(lastId), tagIds)
	if err != nil {
		writeErrResponse(c, err, http.StatusInternalServerError)
		return
	}
	writeOkResponse(c, list, http.StatusOK)

}

func getNote(c *gin.Context) {
	paramId, _ := strconv.Atoi(c.Param("id"))
	note, err := app.GetNotesManager().GetNote(c, uint(paramId))
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

	err = app.GetNotesManager().Save(c, note)
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
	err = app.GetNotesManager().Save(c, note)
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
	err :=  app.GetNotesManager().Delete(c, uint(id))
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
