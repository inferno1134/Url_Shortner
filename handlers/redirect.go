package handlers

import(
	"net/http"


	// "url-shortner/service"
	// "url-shortner/models"
)

func (h *UrlHandler) RedirectURL(


	w http.ResponseWriter,
	r *http.Request,


){

	shortCode := r.URL.Path[1:]

	if shortCode=="" {
		http.Error(
			w,
			"Missing short Code",
			http.StatusBadRequest,
		)

		return
	}


	longUrl, exists := h.service.GetOriginalURL(shortCode)

	if !exists {
		http.Error(
			w,
			"Url not found",
			http.StatusNotFound,
		)

		return 
	}

	http.Redirect(
		w,
		r,
		longUrl,
		http.StatusFound,
	
	)




}

