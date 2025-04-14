package handlers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"html/template"
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

	originalUrl := r.FormValue("url")

	if !isValidURL(originalUrl) {
		h.compo.ExecuteTemplate(w, "EnterValidUrl", nil)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := sha256.Sum256([]byte(originalUrl))
	hashlen := 32
	hashStr := fmt.Sprintf("%x", hash)

	currentShortCodeIdx := 0
	shortCodelen := 6

	err := errors.New("ShortCode Not Created")
	var insertedData db.Shorturl

	for err != nil && currentShortCodeIdx < hashlen-shortCodelen {

		insertedData, err = h.queries.InsertSurl(r.Context(), db.InsertSurlParams{
			ShortCode:   hashStr[currentShortCodeIdx : currentShortCodeIdx+shortCodelen],
			OriginalUrl: originalUrl,
		})

		if err == nil {
			break
		}

		currentShortCodeIdx++

	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templ.PageInternalServerError.Execute(w, nil)
		return
	}

	var data struct {
		Url string
	}

	schema := "http"
	if r.TLS != nil {
		schema = "https"
	}

	data.Url = fmt.Sprintf("%s://%s/r/%s", schema, r.Host, insertedData.ShortCode)

	h.compo.ExecuteTemplate(w, "shorturl", data)

	fmt.Println("inserted : " + insertedData.ShortCode)

}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if id == "" {
		h.compo.ExecuteTemplate(w, "invalidId", nil)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.queries.FineOne(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.compo.ExecuteTemplate(w, "invalidOrExpired", nil)
		return
	}

	if data.OriginalUrl[:4] != "http" {

		data.OriginalUrl = "https://" + data.OriginalUrl
	}

	http.Redirect(w, r, data.OriginalUrl, http.StatusFound)

}
