package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"example.com/user-management/internal/infrastructure/persistence/postgres"
	"example.com/user-management/internal/interface/handler"
	"example.com/user-management/internal/usecase"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

type testServer struct {
	server *atreugo.Atreugo
	db     *sql.DB
}

func setupTestServer(t *testing.T) (*testServer, func()) {
	// Initialize database
	db, err := postgres.NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)

	// Initialize server
	server := atreugo.New(atreugo.Config{})

	// Setup routes
	server.Path("POST", "/users", userHandler.Create)
	server.Path("GET", "/users", userHandler.List)
	server.Path("GET", "/users/:id", userHandler.Get)
	server.Path("PUT", "/users/:id", userHandler.Update)
	server.Path("DELETE", "/users/:id", userHandler.Delete)

	ts := &testServer{
		server: server,
		db:     db,
	}

	// Start server
	serverShutdown := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Logf("Server error: %v", err)
		}
		close(serverShutdown)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	cleanup := func() {
		// Stop accepting new requests
		ts.db.Close()
		<-serverShutdown
	}

	return ts, cleanup
}

func performRequest(method, path string, body []byte) (int, []byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(path)
	if len(body) > 0 {
		req.Header.SetContentType("application/json")
		req.SetBody(body)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return 0, nil, err
	}

	return resp.StatusCode(), resp.Body(), nil
}

func TestUserAPI_Create(t *testing.T) {
	_, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse bool
	}{
		{
			name: "valid user",
			requestBody: map[string]interface{}{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "password123",
			},
			wantStatus:   201,
			wantResponse: true,
		},
		{
			name: "invalid user - missing name",
			requestBody: map[string]interface{}{
				"email":    "john@example.com",
				"password": "password123",
			},
			wantStatus:   400,
			wantResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			statusCode, respBody, err := performRequest("POST", "http://localhost:8080/users", body)
			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}

			if statusCode != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", statusCode, tt.wantStatus)
			}

			if tt.wantResponse {
				var response map[string]interface{}
				if err := json.Unmarshal(respBody, &response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["name"] != tt.requestBody["name"] {
					t.Errorf("Create() name = %v, want %v", response["name"], tt.requestBody["name"])
				}
				if response["email"] != tt.requestBody["email"] {
					t.Errorf("Create() email = %v, want %v", response["email"], tt.requestBody["email"])
				}
			}
		})
	}
}

func TestUserAPI_Update(t *testing.T) {
	_, cleanup := setupTestServer(t)
	defer cleanup()

	// Create a user first
	createBody := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(createBody)
	_, respBody, err := performRequest("POST", "http://localhost:8080/users", body)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var createdUser map[string]interface{}
	if err := json.Unmarshal(respBody, &createdUser); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	tests := []struct {
		name         string
		userID       string
		requestBody  map[string]interface{}
		wantStatus   int
		wantResponse bool
	}{
		{
			name:   "valid update",
			userID: createdUser["id"].(string),
			requestBody: map[string]interface{}{
				"name":  "John Updated",
				"email": "john.updated@example.com",
			},
			wantStatus:   200,
			wantResponse: true,
		},
		{
			name:   "non-existent user",
			userID: "non-existent-id",
			requestBody: map[string]interface{}{
				"name":  "John Updated",
				"email": "john.updated@example.com",
			},
			wantStatus:   404,
			wantResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			statusCode, respBody, err := performRequest("PUT", fmt.Sprintf("http://localhost:8080/users/%s", tt.userID), body)
			if err != nil {
				t.Fatalf("Failed to perform request: %v", err)
			}

			if statusCode != tt.wantStatus {
				t.Errorf("Update() status = %v, want %v", statusCode, tt.wantStatus)
			}

			if tt.wantResponse {
				var response map[string]interface{}
				if err := json.Unmarshal(respBody, &response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["name"] != tt.requestBody["name"] {
					t.Errorf("Update() name = %v, want %v", response["name"], tt.requestBody["name"])
				}
				if response["email"] != tt.requestBody["email"] {
					t.Errorf("Update() email = %v, want %v", response["email"], tt.requestBody["email"])
				}
			}
		})
	}
}
