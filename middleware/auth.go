package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sipatir08/cinema-ticket-api/models"
)

type TokenVerifier interface {
	Verify(token string) (models.AuthClaims, error)
}

type ctxKey struct{}

// Auth: wajib header "Authorization: Bearer <token>" yang valid.
func Auth(v TokenVerifier, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			deny(w, http.StatusUnauthorized, models.ErrUnauthorized)
			return
		}
		claims, err := v.Verify(token)
		if err != nil {
			deny(w, http.StatusUnauthorized, models.ErrUnauthorized)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, claims)))
	}
}

// RequireRole: dipasang di dalam Auth, misalnya hanya admin yang boleh ubah jadwal.
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := r.Context().Value(ctxKey{}).(models.AuthClaims)
		if claims.Role != role {
			deny(w, http.StatusForbidden, models.ErrForbidden)
			return
		}
		next(w, r)
	}
}

func deny(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}