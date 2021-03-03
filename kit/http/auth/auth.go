package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v7"
)

type ClientSecurityInfo struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Realm  string `json:"realm"`
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Authenticate struct {
	Client *ClientSecurityInfo `json:"client"`
	Scope  string              `json:"scope,omitempty"`
	User   *User               `json:"user"`
}

type Refresh struct {
	Client       *ClientSecurityInfo `json:"client"`
	RefreshToken string              `json:"refreshToken,omitempty"`
}

type JWT struct {
	AccessToken      string `json:"accessToken"`
	ExpiresIn        int    `json:"expiresIn"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	TokenType        string `json:"tokenType"`
	NotBeforePolicy  int    `json:"notBeforePolicy"`
	SessionState     string `json:"sessionState"`
	Scope            string `json:"scope"`
}

type Service interface {
	AuthClient(*Authenticate) (*JWT, error)
	AuthUser(*User) (*JWT, error)
	RefreshToken(*Refresh) (*JWT, error)
	CheckAccessToken(accessToken string) error
}

type serviceImpl struct {
	gocloak gocloak.GoCloak
	client  *ClientSecurityInfo
	ctx     context.Context
}

// New instantiates a new Service
// Setting realm is optional
//noinspection GoUnusedExportedFunction
func New(ctx context.Context, gocloak gocloak.GoCloak, client *ClientSecurityInfo) Service {
	return &serviceImpl{
		gocloak: gocloak,
		client:  client,
		ctx:     ctx,
	}
}

func (s *serviceImpl) AuthClient(requestData *Authenticate) (*JWT, error) {
	realm := requestData.Client.Realm
	if realm == "" {
		realm = s.client.Realm
	}

	response, err := s.gocloak.LoginClient(s.ctx, requestData.Client.ID, requestData.Client.Secret, realm)
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

func (s *serviceImpl) AuthUser(requestData *User) (*JWT, error) {

	response, err := s.gocloak.Login(s.ctx, s.client.ID, s.client.Secret, s.client.Realm, requestData.UserName, requestData.Password)
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

func (s *serviceImpl) RefreshToken(requestData *Refresh) (*JWT, error) {
	realm := requestData.Client.Realm
	if realm == "" {
		realm = s.client.Realm
	}

	response, err := s.gocloak.RefreshToken(s.ctx, requestData.RefreshToken, requestData.Client.ID, requestData.Client.Secret, requestData.Client.Realm)
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

func (s *serviceImpl) CheckAccessToken(accessToken string) error {

	result, err := s.gocloak.RetrospectToken(s.ctx, accessToken, s.client.ID, s.client.Secret, s.client.Realm)
	if err != nil {
		return err
	}

	if !*result.Active {
		return fmt.Errorf("invalid or expired token")
	}

	return nil
}