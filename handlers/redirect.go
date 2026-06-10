package handlers

import (
	"log"
	"net/http"
	// "url-shortner/service"
	// "url-shortner/models"
)

func (h *UrlHandler) RedirectURL(


	w http.ResponseWriter,
	r *http.Request,


){

	log.Printf("[Redirect]:Incoming Request for redirect")

	shortCode := r.URL.Path[1:]

	if shortCode=="" {
		http.Error(
			w,
			"Missing short Code",
			http.StatusBadRequest,
		)

		return
	}

	log.Printf("[Redirect] Extracted short code: '%s'", shortCode)

	


	longUrl, err := h.service.GetOriginalURL(shortCode)

	if err!=nil {
		http.Error(
			w,
			"Url not found",
			http.StatusNotFound,
		)

		return 
	}

	log.Printf("[Redirect] Short code found | Code: %s | LongURL: %s", shortCode, longUrl)


	log.Printf("[Redirect] Redirecting from %s to %s", shortCode, longUrl)

	http.Redirect(
		w,
		r,
		longUrl,
		http.StatusFound,
	
	)




}

