package tests

import (
	"testing"
	"time"

	"example.com/user-management/internal/domain/entity"
)

func TestNewUser(t *testing.T) {
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
			name: "empty name",
			input: struct {
				name     string
				email    string
				password string
			}{
				name:     "",
				email:    "john@example.com",
				password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			input: struct {
				name     string
				email    string
				password string
			}{
				name:     "John Doe",
				email:    "",
				password: "password123",
			},
			wantErr: true,
		},
		{
			name: "short password",
			input: struct {
				name     string
				email    string
				password string
			}{
				name:     "John Doe",
				email:    "john@example.com",
				password: "pass",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.input.name, tt.input.email, tt.input.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if user.Name != tt.input.name {
					t.Errorf("NewUser().Name = %v, want %v", user.Name, tt.input.name)
				}
				if user.Email != tt.input.email {
					t.Errorf("NewUser().Email = %v, want %v", user.Email, tt.input.email)
				}
				if user.Password != tt.input.password {
					t.Errorf("NewUser().Password = %v, want %v", user.Password, tt.input.password)
				}
				if user.CreatedAt.IsZero() {
					t.Error("NewUser().CreatedAt is zero")
				}
				if user.UpdatedAt.IsZero() {
					t.Error("NewUser().UpdatedAt is zero")
				}
			}
		})
	}
}

func TestUserUpdate(t *testing.T) {
	user := &entity.User{
		ID:        "test-id",
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name  string
		input struct {
			name  string
			email string
		}
		wantErr bool
	}{
		{
			name: "valid update",
			input: struct {
				name  string
				email string
			}{
				name:  "New Name",
				email: "new@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: struct {
				name  string
				email string
			}{
				name:  "",
				email: "new@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			input: struct {
				name  string
				email string
			}{
				name:  "New Name",
				email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := user.UpdatedAt
			err := user.Update(tt.input.name, tt.input.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if user.Name != tt.input.name {
					t.Errorf("Update().Name = %v, want %v", user.Name, tt.input.name)
				}
				if user.Email != tt.input.email {
					t.Errorf("Update().Email = %v, want %v", user.Email, tt.input.email)
				}
				if !user.UpdatedAt.After(oldUpdatedAt) {
					t.Error("Update().UpdatedAt was not updated")
				}
			}
		})
	}
}
