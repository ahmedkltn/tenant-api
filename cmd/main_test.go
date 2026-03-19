package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedkltn/tenant-api/internal/auth"
	mw "github.com/ahmedkltn/tenant-api/internal/middleware"
	"github.com/ahmedkltn/tenant-api/internal/projects"
	"github.com/ahmedkltn/tenant-api/internal/seed"
	"github.com/ahmedkltn/tenant-api/internal/store"
	"github.com/go-chi/chi/v5"
)

func setupServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := store.New()
	seed.Load(s)
	h := projects.New(s)

	r := chi.NewRouter()
	r.Post("/auth/login", auth.NewLoginHandler(s))

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Route("/tenants/{tenantId}/projects", func(r chi.Router) {
			r.Use(mw.TenantGuard)
			r.Get("/", h.List)
			r.With(mw.RoleCheck(store.RoleAdmin)).Post("/", h.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get)
				r.With(mw.RoleCheck(store.RoleAdmin)).Delete("/", h.Delete)
			})
		})
	})

	return httptest.NewServer(r)
}

func login(t *testing.T, srv *httptest.Server, email, password string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"email": email, "password": password})
	resp, _ := http.Post(srv.URL+"/auth/login", "application/json", bytes.NewReader(body))
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["token"]
}

// test Auth
func TestNoToken_Returns401(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/tenants/tenant-1/projects")
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

// test cross tenant : tenant A trying to access tenant B data
func TestCrossTenant_IsBlocked(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	token := login(t, srv, "admin@tenant1.com", "password123")
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/tenants/tenant-2/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", resp.StatusCode)
	}
}

// RBAC test (delete with viewer role)
func TestViewer_CannotDelete(t *testing.T) {
	srv := setupServer(t)
	defer srv.Close()

	token := login(t, srv, "viewer@tenant1.com", "password123")
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/tenants/tenant-1/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)

	// get a project id
	var projects []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&projects)
	id := projects[0]["id"].(string)

	req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/tenants/tenant-1/projects/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", resp.StatusCode)
	}
}
