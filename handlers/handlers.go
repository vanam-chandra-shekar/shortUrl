package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"short/db"
	"short/templ"
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

func isValidURL(input string) bool {
	// Parse the URL
	parsedURL, err := url.ParseRequestURI(input)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		return true // It's a valid full URL
	}

	// Regex to match domain-like URLs (e.g., example.com, sub.example.com)
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)

	return domainRegex.MatchString(input) // Accepts "example.com" but rejects "example"
}

func (h *Handler) HxOnUrlFormSubmit(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Url string
	}
	data.Url = r.FormValue("url")

	if !isValidURL(data.Url) {
		h.compo.ExecuteTemplate(w, "EnterValidUrl", nil)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urldata, err := h.queries.InsertSurl(context.Background(), data.Url)

	if err != nil {
		log.Printf("[Error] : %v\n", err)
		templ.PageInternalServerError.Execute(w, nil)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = h.queries.UpdateShortCode(context.Background(), urldata.Sid)

	if err != nil {

		h.queries.DeleteSurl(context.Background(), urldata.Sid)
		log.Printf("[Error] : %v\n", err)
		templ.PageInternalServerError.Execute(w, nil)
		http.Error(w, "", http.StatusInternalServerError)

		return
	}

	log.Println(urldata.Sid, urldata.ShortCode, urldata.OriginalUrl, urldata.CreatedAt.Time)

	h.compo.ExecuteTemplate(w, "shorturl", data)
}
