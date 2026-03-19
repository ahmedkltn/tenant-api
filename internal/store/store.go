package store

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleViewer Role = "viewer"
)

type User struct {
	ID       string
	Email    string
	Password string
	TenantID string
	Role     Role
}

type Project struct {
	ID        string
	Name      string
	TenantID  string
	CreatedAt time.Time
}

type Store struct {
	mu       sync.RWMutex
	users    map[string]*User    // key : userId
	projects map[string]*Project // key : projectId
}

func New() *Store {
	return &Store{
		users:    make(map[string]*User),
		projects: make(map[string]*Project),
	}
}

// Users
func (s *Store) AddUser(u *User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.Email] = u
}

func (s *Store) FindUserByEmail(email string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[email]
	return u, ok
}

// Projects CRUD
func (s *Store) CreateProject(name, tenantID string) *Project {
	p := &Project{
		ID:        uuid.NewString(),
		Name:      name,
		TenantID:  tenantID,
		CreatedAt: time.Now().UTC(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projects[p.ID] = p
	return p
}

func (s *Store) ListProjects(tenantID string) []*Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Project
	for _, p := range s.projects {
		if p.TenantID == tenantID {
			result = append(result, p)
		}
	}
	return result
}

func (s *Store) GetProject(id, tenantID string) (*Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.projects[id]
	if !ok || p.TenantID != tenantID {
		return nil, false
	}
	return p, true
}

func (s *Store) DeleteProject(id, tenantID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[id]
	if !ok || p.TenantID != tenantID {
		return false
	}
	delete(s.projects, id)
	return true
}
