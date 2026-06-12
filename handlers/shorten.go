package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tanishk-deore/url-shortner/logger"
	"github.com/tanishk-deore/url-shortner/models"
	"github.com/tanishk-deore/url-shortner/service"
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

	logger.Info("[ShortenUrl] Recieved Request | Method : %s | Content type: %s", r.Method, r.Header.Get("Content-Type"))

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
		logger.Error("[ShortenUrl]: Failed to decode the request : %v", err)
		http.Error(
			w,
			"Invalid Request Body",
			http.StatusBadRequest,
		)

		return

	}

	if req.URL == "" {

		logger.Error("[ShorternUrl]: Url field is Empty : %v", err)
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

	logger.Info("[ShortenUrl]: Short Url is Created")

	logger.Info("[ShortenUrl] Sending Response...")

	response := models.ShortenResponse{

		ShortURL: fmt.Sprintf(
			"http://%s/%s",
			r.Host,
			shortCode,
		),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)

}
