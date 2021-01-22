package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Nerzal/gocloak/v7"
)

// AuthenticationHandler is used to authenticate with the api
type AuthenticationHandler interface {
	AuthenticateClient(Authenticate) (*JWT, error)
	AuthenticateUser(*AuthenticateUser) (*JWT, error)
	RefreshToken(Refresh) (*JWT, error)
}

type authenticationHandler struct {
	gocloak gocloak.GoCloak
	client  *AuthClient
	ctx     context.Context
}

// NewAuthenticationHandler instantiates a new AuthenticationHandler
// Setting realm is optional
//noinspection GoUnusedExportedFunction
func NewAuthenticationHandler(ctx context.Context, gocloak gocloak.GoCloak, client *AuthClient) AuthenticationHandler {
	return &authenticationHandler{
		gocloak: gocloak,
		client:  client,
		ctx:     ctx,
	}
}

func (handler *authenticationHandler) AuthenticateClient(requestData Authenticate) (*JWT, error) {
	realm := requestData.Realm
	if realm == "" {
		realm = handler.client.Realm
	}

	response, err := handler.gocloak.LoginClient(handler.ctx, requestData.ClientID, requestData.ClientSecret, realm)
	if err != nil {
		return nil, gocloak.APIError{
			Code:    403,
			Message: err.Error(),
		}
	}

	if response.AccessToken == "" {
		return nil, errors.New("authentication failed")
	}

	return &JWT{
		AccessToken:      response.AccessToken,
		ExpiresIn:        response.ExpiresIn,
		NotBeforePolicy:  response.NotBeforePolicy,
		RefreshExpiresIn: response.RefreshExpiresIn,
		RefreshToken:     response.RefreshToken,
		Scope:            response.Scope,
		SessionState:     response.SessionState,
		TokenType:        response.TokenType,
	}, nil
}

func (handler *authenticationHandler) AuthenticateUser(requestData *AuthenticateUser) (*JWT, error) {

	response, err := handler.gocloak.Login(handler.ctx, handler.client.ClientID, handler.client.ClientSecret, handler.client.Realm, requestData.UserName, requestData.Password)
	if err != nil {
		return nil, gocloak.APIError{
			Code:    http.StatusForbidden,
			Message: err.Error(),
		}
	}

	if response.AccessToken == "" {
		return nil, errors.New("authentication failed")
	}

	return &JWT{
		AccessToken:      response.AccessToken,
		ExpiresIn:        response.ExpiresIn,
		NotBeforePolicy:  response.NotBeforePolicy,
		RefreshExpiresIn: response.RefreshExpiresIn,
		RefreshToken:     response.RefreshToken,
		Scope:            response.Scope,
		SessionState:     response.SessionState,
		TokenType:        response.TokenType,
	}, nil
}

func (handler *authenticationHandler) RefreshToken(requestData Refresh) (*JWT, error) {
	realm := requestData.Realm
	if realm == "" {
		realm = handler.client.Realm
	}

	response, err := handler.gocloak.RefreshToken(handler.ctx, requestData.RefreshToken, requestData.ClientID, requestData.ClientSecret, requestData.Realm)
	if err != nil {
		return nil, gocloak.APIError{
			Code:    http.StatusForbidden,
			Message: "Failed to refresh token",
		}
	}

	if response.AccessToken == "" {
		return nil, errors.New("authentication failed")
	}

	return &JWT{
		AccessToken:      response.AccessToken,
		ExpiresIn:        response.ExpiresIn,
		NotBeforePolicy:  response.NotBeforePolicy,
		RefreshExpiresIn: response.RefreshExpiresIn,
		RefreshToken:     response.RefreshToken,
		Scope:            response.Scope,
		SessionState:     response.SessionState,
		TokenType:        response.TokenType,
	}, nil
}
