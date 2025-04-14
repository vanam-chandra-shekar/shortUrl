package handlers

import (
	"html/template"
	"net/http"

	"short/db"
)

type Handler struct {
	compo   *template.Template
	queries *db.Queries
}

func NewHandler(templs *template.Template, databaseConn db.DBTX) *Handler {
	return &Handler{
		compo:   templs,
		queries: db.New(databaseConn),
	}

}

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	err := h.compo.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err)
	}
}
