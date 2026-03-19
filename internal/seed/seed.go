package seed

import (
	"github.com/ahmedkltn/tenant-api/internal/store"
)

func Load(s *store.Store) {
	// Tenant 1 users
	s.AddUser(&store.User{
		ID:       "u1",
		Email:    "admin@tenant1.com",
		Password: "password123",
		TenantID: "tenant-1",
		Role:     store.RoleAdmin,
	})
	s.AddUser(&store.User{
		ID:       "u2",
		Email:    "viewer@tenant1.com",
		Password: "password123",
		TenantID: "tenant-1",
		Role:     store.RoleViewer,
	})

	// Tenant 2 users
	s.AddUser(&store.User{
		ID:       "u3",
		Email:    "admin@tenant2.com",
		Password: "password123",
		TenantID: "tenant-2",
		Role:     store.RoleAdmin,
	})

	// Tenant 1 projects
	s.CreateProject("Alpha Launch", "tenant-1")
	s.CreateProject("Beta Programme", "tenant-1")

	// Tenant 2 projects
	s.CreateProject("Go-to-Market", "tenant-2")
	s.CreateProject("Q3 Analytics", "tenant-2")
}
