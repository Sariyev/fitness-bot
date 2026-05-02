package webapp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/service"

	"github.com/jackc/pgx/v4"
)

type MediaHandler struct {
	mediaSvc *service.MediaService
}

func NewMediaHandler(mediaSvc *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaSvc: mediaSvc}
}

type presignRequest struct {
	ContentType   string `json:"content_type"`
	SizeBytes     int64  `json:"size_bytes"`
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   *int64 `json:"reference_id,omitempty"`
	IsPublic      bool   `json:"is_public,omitempty"`
}

// HandleMediaRoutes dispatches /app/api/media and /app/api/media/{id}{,/confirm}.
func (h *MediaHandler) HandleMediaRoutes(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/app/api/media")
	path = strings.TrimPrefix(path, "/")

	// /app/api/media/presign
	if path == "presign" && r.Method == http.MethodPost {
		var req presignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		res, err := h.mediaSvc.RequestUpload(r.Context(), user, service.UploadRequest{
			ContentType:   req.ContentType,
			SizeBytes:     req.SizeBytes,
			ReferenceType: req.ReferenceType,
			ReferenceID:   req.ReferenceID,
			IsPublic:      req.IsPublic,
		})
		if err != nil {
			handleMediaError(w, err)
			return
		}
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"upload_url": res.UploadURL,
			"media_id":   res.MediaID,
			"key":        res.Key,
		})
		return
	}

	// /app/api/media/{id}{,/confirm}
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid media id")
		return
	}

	// /app/api/media/{id}/confirm  (POST)
	if len(parts) == 2 && parts[1] == "confirm" && r.Method == http.MethodPost {
		m, err := h.mediaSvc.ConfirmUpload(r.Context(), user, id)
		if err != nil {
			handleMediaError(w, err)
			return
		}
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"success":    true,
			"media_id":   m.ID,
			"size_bytes": m.SizeBytes,
		})
		return
	}

	// /app/api/media/{id}  (GET, DELETE)
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			url, err := h.mediaSvc.GetURL(r.Context(), user, id)
			if err != nil {
				handleMediaError(w, err)
				return
			}
			jsonResponse(w, http.StatusOK, map[string]string{"url": url})
			return
		case http.MethodDelete:
			if err := h.mediaSvc.Delete(r.Context(), user, id); err != nil {
				handleMediaError(w, err)
				return
			}
			jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
			return
		}
	}

	jsonError(w, http.StatusNotFound, "not found")
}

func handleMediaError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		jsonError(w, http.StatusNotFound, "media not found")
	case errors.Is(err, service.ErrInvalidContentType):
		jsonError(w, http.StatusUnsupportedMediaType, "content-type not allowed")
	case errors.Is(err, service.ErrSizeExceedsTypeCap), errors.Is(err, service.ErrSizeNotPositive):
		jsonError(w, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, service.ErrQuotaExceeded):
		jsonError(w, http.StatusInsufficientStorage, "storage quota exceeded")
	case errors.Is(err, service.ErrPublicRequiresAdmin):
		jsonError(w, http.StatusForbidden, "admin only")
	case errors.Is(err, service.ErrNotOwner):
		jsonError(w, http.StatusForbidden, "not owner")
	case errors.Is(err, service.ErrSizeMismatch):
		jsonError(w, http.StatusBadRequest, err.Error())
	default:
		jsonError(w, http.StatusInternalServerError, err.Error())
	}
}
