package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dangrt/playliststore/helpers"
	"github.com/joho/godotenv"
)

func main() {

	date := time.Now().Format("2006-01-02")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	clientID := os.Getenv("clientID")
	clientSecret := os.Getenv("clientSecret")
	db := helpers.InitDB("./spotify.db")
	client, err := helpers.InitialiseAuth(clientID, clientSecret)
	if err != nil {
		fmt.Printf("Error initializing authentication: %v\n", err)
		return
	}

	// Get the current user
	user, err := (*client).CurrentUser()
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
		return
	}

	fmt.Printf("Welcome, %s!\n", user.DisplayName)

	playlists, err := (*client).GetPlaylistsForUser(user.ID)
	if err != nil {
		fmt.Printf("Error getting playlists: %v\n", err)
		return
	}

	for _, playlist := range playlists.Playlists {
		pl := helpers.UserPlaylist{
			ID:        string(playlist.ID) + date,
			SpotifyID: string(playlist.ID),
			Name:      playlist.Name,
			Date:      date,
		}

		helpers.InsertPlaylist(db, pl)

		// Get tracks for playlist
		tracks, err := (*client).GetPlaylistTracks(playlist.ID)
		if err != nil {
			fmt.Printf("Error getting tracks for playlist: %v\n", err)
			return
		}

		helpers.InsertTracks(db, tracks.Tracks, pl)
	}
}
