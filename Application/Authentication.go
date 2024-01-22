package Application

type RefreshToken struct {
	Token        string `json:"refresh_token"`
	SessionState string `json:"session_state"`
}

type AccessToken struct {
	Token        string `json:"access_token"`
	SessionState string `json:"session_state"`
}

type Credentials struct {
	Username string
	Password string
}

type Authentication struct {
	userPass     Credentials
	refreshToken RefreshToken
}
