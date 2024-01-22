package Application

import (
	"fmt"
	"github.com/zalando/go-keyring"
	"os"
)

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

func Login(username, password string) (Authentication, error) {
	RefreshToken, err := GetRefreshToken(Credentials{Username: username, Password: password})

	if err != nil {
		return Authentication{}, err
	}

	auth := &Authentication{
		userPass: Credentials{
			Username: username,
			Password: password,
		},
		refreshToken: RefreshToken,
	}

	return *auth, nil
}

func (auth *Authentication) RefreshCredentials() error {
	refreshToken, err := GetRefreshToken(auth.userPass)
	if err != nil {
		return err
	}
	auth.refreshToken = refreshToken
	return nil
}

func (auth *Authentication) GetAccessToken() (AccessToken, error) {
	accessToken, err := GetAccessToken(auth.refreshToken)
	if err != nil {
		return AccessToken{}, err
	}
	return accessToken, nil
}

func saveUsername(username string) error {
	f, err := os.Create("username.txt")
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	_, err = f.WriteString(username)
	if err != nil {
		return err
	}

	return nil
}

func loadUsername() (string, error) { // <-- Added error return type
	f, err := os.Open("username.txt")
	if err != nil {
		return "", err
	}
	defer f.Close() // <-- Use defer to ensure file closure

	var storedUsername string
	_, err = fmt.Fscanf(f, "%s", &storedUsername)
	if err != nil {
		return "", err
	}

	return storedUsername, nil
}

// save credentials to keyring
func (auth *Authentication) SaveCredentials() error {
	// save credentials to keyring
	if auth.RefreshCredentials() != nil {
		return fmt.Errorf("invalid credentials")
	}

	err := keyring.Set("SentinelDownload", auth.userPass.Username, auth.userPass.Password)
	if err != nil {
		return err
	}

	err2 := saveUsername(auth.userPass.Username)
	if err2 != nil {
		return err2
	}

	return nil
}

// load credentials from keyring
func LoadCredentials() (Authentication, error) {
	// load username
	username, err := loadUsername()
	if err != nil {
		return Authentication{}, err
	}

	// load password
	password, err := keyring.Get("SentinelDownload", username)
	if err != nil {
		return Authentication{}, err
	}

	login, err := Login(username, password)
	if err != nil {
		return Authentication{}, err
	}

	return login, nil
}
