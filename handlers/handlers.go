package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	compo *template.Template
}

func NewHandler(templs *template.Template) *Handler {
	return &Handler{
		compo: templs,
	}

}

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	err := h.compo.ExecuteTemplate(w, "index", nil)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HxOnUrlFormSubmit(w http.ResponseWriter, r *http.Request) {


	log.Println(r.FormValue("url"))

	var data struct{
		Url string
	};
	data.Url = r.FormValue("url")
	h.compo.ExecuteTemplate(w, "shorturl", data)
}
