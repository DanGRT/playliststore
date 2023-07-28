package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	redirectURI = "http://localhost:8888/callback"
	tokenFile   = "/tmp/spotify-playlist-token"
)

func InitialiseAuth(clientID, clientSecret string) (*spotify.Client, error) {
	auth := spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistReadCollaborative, spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(clientID, clientSecret)

	token, err := loadAccessToken(tokenFile)
	if err != nil {
		token, err = Authenticate(auth)
		if err != nil {
			return nil, err
		}

		err = saveAccessToken(tokenFile, token)
		if err != nil {
			return nil, err
		}
	}

	client := auth.NewClient(token)
	return &client, nil
}

func Authenticate(auth spotify.Authenticator) (*oauth2.Token, error) {
	authURL := auth.AuthURL("playlistdb")
	fmt.Println("authenticate your Spotify account at:")
	fmt.Println(authURL)

	fmt.Print("url you were returned to: ")
	var url string
	fmt.Scan(&url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return auth.Token("playlistdb", req)
}

func loadAccessToken(filename string) (*oauth2.Token, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func saveAccessToken(filename string, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
