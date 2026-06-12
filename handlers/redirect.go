package handlers

import (
	"net/http"
	"errors"
	"database/sql"
	

	"github.com/tanishk-deore/url-shortner/logger"
	
	// "url-shortner/service"
	// "url-shortner/models"
)

func (h *UrlHandler) RedirectURL(

	w http.ResponseWriter,
	r *http.Request,

) {

	logger.Info("[Redirect]:Incoming Request for redirect")

	shortCode := r.URL.Path[1:]

	if shortCode == "" {
		http.Error(
			w,
			"Missing short Code",
			http.StatusBadRequest,
		)

		return
	}

	logger.Info("[Redirect] Extracted short code: '%s'", shortCode)

	// longUrl, err := h.service.GetOriginalURL(shortCode)

	// if err != nil {
	// 	http.Error(
	// 		w,
	// 		"Url not found",
	// 		http.StatusNotFound,
	// 	)

	// 	return
	// }


	longUrl, err := h.service.GetOriginalURL(shortCode)

	if err != nil {

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(
			w,
			"Url not found",
			http.StatusNotFound,
		)
		return
	}

	http.Error(
		w,
		"Internal Server Error",
		http.StatusInternalServerError,
	)
	return
	}

	logger.Info("[Redirect] Short code found | Code: %s | LongURL: %s", shortCode, longUrl)

	logger.Info("[Redirect] Redirecting from %s to %s", shortCode, longUrl)

	http.Redirect(
		w,
		r,
		longUrl,
		http.StatusFound,
	)

}
