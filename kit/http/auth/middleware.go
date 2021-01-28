package auth

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"github.com/Nerzal/gocloak/v7/pkg/jwx"
	"github.com/dgrijalva/jwt-go/v4"
	"net/http"
	"strings"
)

type Middleware interface {
	DecodeAndValidateToken(next http.Handler) http.Handler
	CheckToken(next http.Handler) http.Handler
	CheckTokenCustomHeader(next http.Handler) http.Handler
	CheckScope(next http.Handler) http.Handler
}

// NewMdw instantiates a new Middleware when using the Keycloak Direct Grant aka
// Resource Owner Password Credentials Flow
//
// see https://www.keycloak.org/docs/latest/securing_apps/index.html#_resource_owner_password_credentials_flow and
// https://tools.ietf.org/html/rfc6749#section-4.3 for more information about this flow
//noinspection GoUnusedExportedFunction
func NewMdw(ctx context.Context, gocloak gocloak.GoCloak, client *ClientSecurityInfo, allowedScope string, customHeaderName string) Middleware {
	return &middleware{
		gocloak:          gocloak,
		allowedScope:     allowedScope,
		customHeaderName: customHeaderName,
		ctx:              ctx,
		client:           client,
	}
}

type middleware struct {
	gocloak          gocloak.GoCloak
	client           *ClientSecurityInfo
	allowedScope     string
	customHeaderName string
	ctx              context.Context
}

func (m *middleware) tokenFromHeader(r *http.Request) string {
	token := ""

	if m.customHeaderName != "" {
		token = r.Header.Get(m.customHeaderName)
	}

	if token == "" {
		token = r.Header.Get("Authorization")
	}

	return token
}

// CheckTokenCustomHeader used to verify authorization tokens
func (m *middleware) CheckTokenCustomHeader(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := m.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		decodedToken, err := m.stripBearerAndCheckToken(token, m.client.Realm)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		if !decodedToken.Valid {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)

}

func (m *middleware) stripBearerAndCheckToken(accessToken string, realm string) (*jwt.Token, error) {
	accessToken = extractBearerToken(accessToken)
	decodedToken, _, err := m.gocloak.DecodeAccessToken(m.ctx, accessToken, realm, "")
	return decodedToken, err
}

func (m *middleware) DecodeAndValidateToken(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := m.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)

}

// CheckToken used to verify authorization tokens
func (m *middleware) CheckToken(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := m.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token = extractBearerToken(token)

		if token == "" {
			http.Error(w, "Bearer Token missing", http.StatusUnauthorized)
			return
		}

		result, err := m.gocloak.RetrospectToken(m.ctx, token, m.client.ID, m.client.Secret, m.client.Realm)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		if !*result.Active {
			http.Error(w, "Invalid or expired Token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func extractBearerToken(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func (m *middleware) CheckScope(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := m.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token = extractBearerToken(token)
		claims := &jwx.Claims{}
		_, err := m.gocloak.DecodeAccessTokenCustomClaims(m.ctx, token, m.client.Realm, "", claims)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		if !strings.Contains(claims.Scope, m.allowedScope) {
			http.Error(w, "Insufficient permissions to access the requested resource", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
