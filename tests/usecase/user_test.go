package usecase_test

import (
	"context"
	"errors"
	"testing"

	"example.com/user-management/internal/domain/entity"
	"example.com/user-management/internal/domain/repository"
	"example.com/user-management/internal/usecase"
)

type mockUserRepository struct {
	users map[string]*entity.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *mockUserRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if user, exists := r.users[email]; exists {
		return user, nil
	}
	return nil, repository.ErrNotFound
}

func (r *mockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if _, err := r.Get(ctx, user.ID); err != nil {
		return err
	}
	r.users[user.Email] = user
	return nil
}

func (r *mockUserRepository) Delete(ctx context.Context, id string) error {
	for email, user := range r.users {
		if user.ID == id {
			delete(r.users, email)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (r *mockUserRepository) List(ctx context.Context) ([]*entity.User, error) {
	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func TestUserUseCase_Create(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUseCase(repo)
	ctx := context.Background()

	tests := []struct {
		name  string
		input struct {
			name     string
			email    string
			password string
		}
		wantErr bool
	}{
		{
			name: "valid user",
			input: struct {
				name     string
				email    string
				password string
			}{
				name:     "John Doe",
				email:    "john@example.com",
				password: "password123",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			input: struct {
				name     string
				email    string
				password string
			}{
				name:     "John Doe",
				email:    "john@example.com",
				password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := uc.Create(ctx, tt.input.name, tt.input.email, tt.input.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if user.Name != tt.input.name {
					t.Errorf("Create().Name = %v, want %v", user.Name, tt.input.name)
				}
				if user.Email != tt.input.email {
					t.Errorf("Create().Email = %v, want %v", user.Email, tt.input.email)
				}
			}
		})
	}
}

func TestUserUseCase_Update(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUseCase(repo)
	ctx := context.Background()

	// Create a user first
	user, err := uc.Create(ctx, "John Doe", "john@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name  string
		input struct {
			id    string
			name  string
			email string
		}
		wantErr bool
	}{
		{
			name: "valid update",
			input: struct {
				id    string
				name  string
				email string
			}{
				id:    user.ID,
				name:  "John Updated",
				email: "john.updated@example.com",
			},
			wantErr: false,
		},
		{
			name: "non-existent user",
			input: struct {
				id    string
				name  string
				email string
			}{
				id:    "non-existent-id",
				name:  "John Updated",
				email: "john.updated@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedUser, err := uc.Update(ctx, tt.input.id, tt.input.name, tt.input.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if updatedUser.Name != tt.input.name {
					t.Errorf("Update().Name = %v, want %v", updatedUser.Name, tt.input.name)
				}
				if updatedUser.Email != tt.input.email {
					t.Errorf("Update().Email = %v, want %v", updatedUser.Email, tt.input.email)
				}
			}
		})
	}
}
