package handlers

import "net/http"

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
