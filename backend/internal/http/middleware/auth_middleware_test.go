package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"sumerki/backend/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func TestAuthRejectsMissingAuthorizationHeader(t *testing.T) {
	rec := runAuthMiddleware(t, "", service.NewAuthService(nil, "test-secret"))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "missing_authorization_header") {
		t.Fatalf("expected missing_authorization_header, got %s", rec.Body.String())
	}
}

func TestAuthRejectsInvalidToken(t *testing.T) {
	rec := runAuthMiddleware(t, "Bearer bad-token", service.NewAuthService(nil, "test-secret"))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid_token") {
		t.Fatalf("expected invalid_token, got %s", rec.Body.String())
	}
}

func TestAuthRejectsExpiredToken(t *testing.T) {
	auth := service.NewAuthService(nil, "test-secret")
	claims := service.Claims{
		UserID: "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-26 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	rec := runAuthMiddleware(t, "Bearer "+token, auth)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "expired_token") {
		t.Fatalf("expected expired_token, got %s", rec.Body.String())
	}
}

func runAuthMiddleware(t *testing.T, authorization string, auth *service.AuthService) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := Auth(auth)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	if err := handler(c); err != nil {
		t.Fatalf("middleware returned error: %v", err)
	}

	return rec
}
