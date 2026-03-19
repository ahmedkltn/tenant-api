package projects

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ahmedkltn/tenant-api/internal/store"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

type projectResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	projects := h.store.ListProjects(tenantID)

	resp := make([]projectResponse, 0, len(projects))
	for _, p := range projects {
		resp = append(resp, projectResponse{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")

	var req struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	p := h.store.CreateProject(req.Name, tenantID)
	writeJSON(w, http.StatusCreated, projectResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	id := chi.URLParam(r, "id")

	p, ok := h.store.GetProject(id, tenantID)
	if !ok {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, projectResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantId")
	id := chi.URLParam(r, "id")

	if ok := h.store.DeleteProject(id, tenantID); !ok {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
