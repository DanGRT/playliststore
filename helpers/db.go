package helpers

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
)

type UserPlaylist struct {
	ID        string
	SpotifyID string
	Name      string
	Date      string
}

type Song struct {
	SongID     string
	Name       string
	Artist     string
	DurationMS int
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Println(err)
	}
	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createTable(db *sql.DB) error {
	sqlTable := `CREATE TABLE IF NOT EXISTS Playlists(
		ID TEXT NOT NULL PRIMARY KEY,
		Name TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS Songs(
		SongID TEXT NOT NULL,
		Name TEXT,
		Artist TEXT,
		DurationMS INTEGER,
		PRIMARY KEY(SongID)
	);
	
	CREATE TABLE IF NOT EXISTS PlaylistSongMapping(
		PlaylistID TEXT NOT NULL,
		SongID TEXT NOT NULL,
		Date TEXT NOT NULL,
		FOREIGN KEY (PlaylistID) REFERENCES Playlists(ID),
		FOREIGN KEY (SongID) REFERENCES Songs(SongID),
		PRIMARY KEY(PlaylistID, SongID, Date)
	);
	`

	_, err := db.Exec(sqlTable)
	if err != nil {
		return err
	}

	return nil
}

func InsertPlaylist(db *sql.DB, pl UserPlaylist) error {
	sqlAddItem := `
	INSERT OR IGNORE INTO Playlists(
		ID,
		Name
	) values(?, ?)
	`

	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(pl.ID, pl.Name)
	if err2 != nil {
		return err2
	}

	return nil
}

func InsertTracks(db *sql.DB, tracks []spotify.PlaylistTrack, pl UserPlaylist, date string) error {
	sqlAddItem := `
	INSERT OR IGNORE INTO Songs(
		SongID,
		Name,
		Artist,
		DurationMS
	) values(?, ?, ?, ?)
	`

	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, track := range tracks {
		s := Song{
			SongID:     string(track.Track.ID),
			Artist:     track.Track.Artists[0].Name,
			Name:       track.Track.Name,
			DurationMS: track.Track.Duration,
		}

		_, err2 := stmt.Exec(s.SongID, s.Name, s.Artist, s.DurationMS)
		if err2 != nil {
			return err2
		}
		InsertPlaylistSongMapping(db, tracks, pl.ID, date)
	}

	return nil
}

func InsertPlaylistSongMapping(db *sql.DB, tracks []spotify.PlaylistTrack, playlistID string, date string) error {
	sqlAddMapping := `
	INSERT INTO PlaylistSongMapping(
		PlaylistID,
		SongID,
		Date
	) values(?, ?, ?)
	`

	stmt, err := db.Prepare(sqlAddMapping)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, track := range tracks {
		_, err2 := stmt.Exec(playlistID, string(track.Track.ID), date)
		if err2 != nil {
			return err2
		}
	}

	return nil
}
