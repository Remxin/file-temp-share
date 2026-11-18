//go:build integration

package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/dto"
	"file-temp-share/backend/internal/handler"
	"file-temp-share/backend/internal/platform/database"
	"file-temp-share/backend/internal/platform/httpx"
	"file-temp-share/backend/internal/queue"
	"file-temp-share/backend/internal/repository"
	"file-temp-share/backend/internal/server"
	"file-temp-share/backend/internal/service"
)

func TestAPI_Health(t *testing.T) {
	engine := newTestAPI(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status field = %q, want ok", body["status"])
	}
}

func TestAPI_RegisterLoginAndMe(t *testing.T) {
	engine := newTestAPI(t)
	email := "int-" + uuid.NewString() + "@filetemp.local"
	password := "password123"

	reg := postJSON(t, engine, "/api/v1/auth/register", dto.RegisterRequest{
		Email:    email,
		Password: password,
	}, "")
	if reg.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want 201; body=%s", reg.Code, reg.Body.String())
	}
	var authRes dto.AuthResponse
	decodeJSON(t, reg, &authRes)
	if authRes.Token == "" {
		t.Fatal("register: empty token")
	}
	if authRes.User.Email != email {
		t.Fatalf("register email = %q, want %q", authRes.User.Email, email)
	}

	login := postJSON(t, engine, "/api/v1/auth/login", dto.LoginRequest{
		Email:    email,
		Password: password,
	}, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d, want 200; body=%s", login.Code, login.Body.String())
	}
	var loginRes dto.AuthResponse
	decodeJSON(t, login, &loginRes)
	if loginRes.Token == "" {
		t.Fatal("login: empty token")
	}

	me := getJSON(t, engine, "/api/v1/auth/me", loginRes.Token)
	if me.Code != http.StatusOK {
		t.Fatalf("me status = %d, want 200; body=%s", me.Code, me.Body.String())
	}
	var user dto.UserResponse
	decodeJSON(t, me, &user)
	if user.Email != email {
		t.Fatalf("me email = %q, want %q", user.Email, email)
	}
	if !user.IsActive {
		t.Fatal("me: expected active user")
	}
}

func TestAPI_TransfersRequireAuth(t *testing.T) {
	engine := newTestAPI(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transfers", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body=%s", w.Code, w.Body.String())
	}
	var errBody httpx.ErrorBody
	decodeJSON(t, w, &errBody)
	if errBody.Error != "unauthorized" {
		t.Fatalf("error = %q, want unauthorized", errBody.Error)
	}
}

func TestAPI_LoginRejectsBadPassword(t *testing.T) {
	engine := newTestAPI(t)
	email, password := uniqueCreds()
	_ = registerUser(t, engine, email, password)

	w := postJSON(t, engine, "/api/v1/auth/login", dto.LoginRequest{
		Email:    email,
		Password: "wrong-password",
	}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body=%s", w.Code, w.Body.String())
	}
}

func TestAPI_RegisterConflictOnDuplicateEmail(t *testing.T) {
	engine := newTestAPI(t)
	email, password := uniqueCreds()
	_ = registerUser(t, engine, email, password)

	w := postJSON(t, engine, "/api/v1/auth/register", dto.RegisterRequest{
		Email:    email,
		Password: password,
	}, "")
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", w.Code, w.Body.String())
	}
}

func TestAPI_CreateTransferAndGetPublic(t *testing.T) {
	engine := newTestAPI(t)
	email, password := uniqueCreds()
	token := registerUser(t, engine, email, password)

	title := "int transfer"
	create := postJSON(t, engine, "/api/v1/transfers", dto.CreateTransferRequest{
		Title:          &title,
		ExpiresInHours: 24,
	}, token)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body=%s", create.Code, create.Body.String())
	}
	var transfer dto.TransferResponse
	decodeJSON(t, create, &transfer)
	if transfer.ShortCode == "" {
		t.Fatal("create: empty short_code")
	}
	if transfer.Status != "active" {
		t.Fatalf("status = %q, want active", transfer.Status)
	}

	pub := getJSON(t, engine, "/api/v1/public/transfers/"+transfer.ShortCode, "")
	if pub.Code != http.StatusOK {
		t.Fatalf("public status = %d, want 200; body=%s", pub.Code, pub.Body.String())
	}
	var public dto.PublicTransferResponse
	decodeJSON(t, pub, &public)
	if public.ShortCode != transfer.ShortCode {
		t.Fatalf("public short_code = %q, want %q", public.ShortCode, transfer.ShortCode)
	}
	if public.Title == nil || *public.Title != title {
		t.Fatalf("public title = %v, want %q", public.Title, title)
	}
	if public.SenderEmail != email {
		t.Fatalf("sender_email = %q, want %q", public.SenderEmail, email)
	}
}

func TestAPI_ChangePasswordThenLogin(t *testing.T) {
	engine := newTestAPI(t)
	email, password := uniqueCreds()
	token := registerUser(t, engine, email, password)
	newPassword := "new-password-99"

	change := postJSON(t, engine, "/api/v1/auth/password", dto.ChangePasswordRequest{
		CurrentPassword: password,
		NewPassword:     newPassword,
	}, token)
	if change.Code != http.StatusNoContent {
		t.Fatalf("change password status = %d, want 204; body=%s", change.Code, change.Body.String())
	}

	oldLogin := postJSON(t, engine, "/api/v1/auth/login", dto.LoginRequest{
		Email:    email,
		Password: password,
	}, "")
	if oldLogin.Code != http.StatusUnauthorized {
		t.Fatalf("old password login status = %d, want 401; body=%s", oldLogin.Code, oldLogin.Body.String())
	}

	newLogin := postJSON(t, engine, "/api/v1/auth/login", dto.LoginRequest{
		Email:    email,
		Password: newPassword,
	}, "")
	if newLogin.Code != http.StatusOK {
		t.Fatalf("new password login status = %d, want 200; body=%s", newLogin.Code, newLogin.Body.String())
	}
}

func TestAPI_PublicTransferNotFound(t *testing.T) {
	engine := newTestAPI(t)

	w := getJSON(t, engine, "/api/v1/public/transfers/doesnotexist", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String())
	}
}

func newTestAPI(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	pool, err := database.NewPool(ctx, databaseURL())
	if err != nil {
		t.Skipf("postgres unavailable (start with task infra:up:db): %v", err)
	}
	t.Cleanup(pool.Close)

	tokens := auth.NewTokenManager("integration-test-jwt-secret", time.Hour)
	users := repository.NewUserRepository(pool)
	resets := repository.NewPasswordResetRepository(pool)
	transfers := repository.NewTransferRepository(pool)

	authSvc := service.NewAuthService(users, resets, tokens, noopPublisher{}, "http://localhost:3000")
	transferSvc := service.NewTransferService(transfers, users, noopStorage{}, noopPublisher{}, "http://localhost:3000", 1<<20)
	adminSvc := service.NewAdminService(users, transfers, "http://localhost:3000")

	return server.New(server.Dependencies{
		AuthHandler:     handler.NewAuthHandler(authSvc),
		TransferHandler: handler.NewTransferHandler(transferSvc, 1<<20),
		AdminHandler:    handler.NewAdminHandler(adminSvc),
		TokenManager:    tokens,
	})
}

func databaseURL() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://filetemp:filetemp@localhost:5432/filetemp?sslmode=disable"
}

func uniqueCreds() (email, password string) {
	return "int-" + uuid.NewString() + "@filetemp.local", "password123"
}

func registerUser(t *testing.T, engine http.Handler, email, password string) string {
	t.Helper()
	w := postJSON(t, engine, "/api/v1/auth/register", dto.RegisterRequest{
		Email:    email,
		Password: password,
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
	var authRes dto.AuthResponse
	decodeJSON(t, w, &authRes)
	if authRes.Token == "" {
		t.Fatal("register: empty token")
	}
	return authRes.Token
}

func postJSON(t *testing.T, engine http.Handler, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	engine.ServeHTTP(w, req)
	return w
}

func getJSON(t *testing.T, engine http.Handler, path, token string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	engine.ServeHTTP(w, req)
	return w
}

func decodeJSON(t *testing.T, w *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), dst); err != nil {
		t.Fatalf("decode json: %v; body=%s", err, w.Body.String())
	}
}

type noopPublisher struct{}

func (noopPublisher) PublishPasswordReset(context.Context, queue.PasswordResetEmailMessage) error {
	return nil
}
func (noopPublisher) PublishTransferInvite(context.Context, queue.TransferInviteEmailMessage) error {
	return nil
}
func (noopPublisher) Close() error { return nil }

type noopStorage struct{}

func (noopStorage) Put(context.Context, string, io.Reader, int64, string) error { return nil }
func (noopStorage) Get(context.Context, string) (io.ReadCloser, error) {
	return nil, errors.New("noop storage")
}
func (noopStorage) Delete(context.Context, string) error { return nil }
