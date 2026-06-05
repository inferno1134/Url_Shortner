package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"url-shortner/models"
	"url-shortner/service"
)

type UrlHandler struct {
	service *service.ShortnerService
}

func NewUrlHandler(service *service.ShortnerService) *UrlHandler {

	return &UrlHandler{
		service: service,
	}

}

func (h *UrlHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return

	}

	var req models.ShortenReqest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"Invalid Request Body",
			http.StatusBadRequest,
		)

		return

	}

	if req.URL == "" {
		http.Error(
			w,
			"Url Cannot be Empty",
			http.StatusBadRequest,
		)

		return
	}

	shortCode, err := h.service.CreateShortURL(req.URL)

	if err != nil {
		http.Error(
			w,
			"Failed to Create URL",
			http.StatusInternalServerError,
		)

		return
	}

	response := models.ShortenResponse{

		ShortURL: fmt.Sprintf(
			"http://localhost:8080/%s",
			shortCode,
		),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)

}
