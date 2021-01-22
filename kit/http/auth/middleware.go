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

type AuthMiddleware interface {
	DecodeAndValidateToken(next http.Handler) http.Handler
	CheckToken(next http.Handler) http.Handler
	CheckTokenCustomHeader(next http.Handler) http.Handler
	CheckScope(next http.Handler) http.Handler
}

// NewKeyCloakMdw instantiates a new AuthMiddleware when using the Keycloak Direct Grant aka
// Resource Owner Password Credentials Flow
//
// see https://www.keycloak.org/docs/latest/securing_apps/index.html#_resource_owner_password_credentials_flow and
// https://tools.ietf.org/html/rfc6749#section-4.3 for more information about this flow
//noinspection GoUnusedExportedFunction
func NewKeyCloakMdw(ctx context.Context, gocloak gocloak.GoCloak, client *AuthClient, allowedScope string, customHeaderName string) AuthMiddleware {
	return &keyCloakMiddleware{
		gocloak:          gocloak,
		realm:            client.Realm,
		allowedScope:     allowedScope,
		customHeaderName: customHeaderName,
		clientID:         client.ClientID,
		clientSecret:     client.ClientSecret,
		ctx:              ctx,
	}
}

type keyCloakMiddleware struct {
	gocloak          gocloak.GoCloak
	realm            string
	clientID         string
	clientSecret     string
	allowedScope     string
	customHeaderName string
	ctx              context.Context
}

func (auth *keyCloakMiddleware) tokenFromHeader(r *http.Request) string {
	token := ""

	if auth.customHeaderName != "" {
		token = r.Header.Get(auth.customHeaderName)
	}

	if token == "" {
		token = r.Header.Get("Authorization")
	}

	return token
}

// CheckTokenCustomHeader used to verify authorization tokens
func (auth *keyCloakMiddleware) CheckTokenCustomHeader(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := auth.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		decodedToken, err := auth.stripBearerAndCheckToken(token, auth.realm)
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

func (auth *keyCloakMiddleware) stripBearerAndCheckToken(accessToken string, realm string) (*jwt.Token, error) {
	accessToken = extractBearerToken(accessToken)
	decodedToken, _, err := auth.gocloak.DecodeAccessToken(auth.ctx, accessToken, realm, "")
	return decodedToken, err
}

func (auth *keyCloakMiddleware) DecodeAndValidateToken(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := auth.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)

}

// CheckToken used to verify authorization tokens
func (auth *keyCloakMiddleware) CheckToken(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := auth.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token = extractBearerToken(token)

		if token == "" {
			http.Error(w, "Bearer Token missing", http.StatusUnauthorized)
			return
		}

		result, err := auth.gocloak.RetrospectToken(auth.ctx, token, auth.clientID, auth.clientSecret, auth.realm)
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

func (auth *keyCloakMiddleware) CheckScope(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		token := auth.tokenFromHeader(r)

		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token = extractBearerToken(token)
		claims := &jwx.Claims{}
		_, err := auth.gocloak.DecodeAccessTokenCustomClaims(auth.ctx, token, auth.realm, "", claims)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid or malformed token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		if !strings.Contains(claims.Scope, auth.allowedScope) {
			http.Error(w, "Insufficient permissions to access the requested resource", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}