package main

import (
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

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	dbPath := os.Getenv("DB_PATH")

	db := helpers.InitDB(dbPath)
	client, err := helpers.InitialiseAuth(clientID, clientSecret)
	if err != nil {
		log.Printf("Error initializing authentication: %v\n", err)
		return
	}

	user, err := (*client).CurrentUser()
	if err != nil {
		log.Printf("Error getting current user: %v\n", err)
		return
	}

	log.Printf("Connected to spotify as, %s!\n", user.DisplayName)

	playlists, err := (*client).GetPlaylistsForUser(user.ID)
	if err != nil {
		log.Printf("Error getting playlists: %v\n", err)
		return
	}

	for _, playlist := range playlists.Playlists {
		pl := helpers.UserPlaylist{
			ID:   string(playlist.ID),
			Name: playlist.Name,
		}

		helpers.InsertPlaylist(db, pl)

		tracks, err := (*client).GetPlaylistTracks(playlist.ID)
		if err != nil {
			log.Printf("Error getting tracks for playlist: %v\n", err)
			return
		}

		helpers.InsertTracks(db, tracks.Tracks, pl, date)
	}
}
