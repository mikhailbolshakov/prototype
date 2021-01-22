package auth

type AuthClient struct {
	ClientID     string  `json:"clientID"`
	ClientSecret string  `json:"clientSecret"`
	Realm        string  `json:"realm"`
}

// Authenticate holds authentication information
type Authenticate struct {
	ClientID     string  `json:"clientID"`
	ClientSecret string  `json:"clientSecret"`
	Realm        string  `json:"realm,omitempty"`
	Scope        string  `json:"scope,omitempty"`
	UserName     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

type AuthenticateUser struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
}

// Refresh is used to refresh the JWT
type Refresh struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Realm        string `json:"realm,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// JWT is a JWT
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