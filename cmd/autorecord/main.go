package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jccroft1/autorecord/internal/camera"
	"github.com/jccroft1/autorecord/internal/spotify"
	"github.com/jccroft1/autorecord/internal/vision"
	"github.com/jccroft1/autorecord/internal/web"
)

const (
	skipCamera        = false
	skipImageSearch   = false
	skipSpotifySearch = false

	errorText = "Oops, something went wrong..."

	spotifyAuthText   = "Setup your Spotify account. We'll redirect you to login to Spotify so you can approve this app."
	spotifyPlayerText = "We need to choose a default player for the music playback. You'll need to be signed into your Spotify account on that device."
)

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/spotify/auth", spotifyAuth)
	http.HandleFunc("/spotify/callback", spotifyCallback)
	http.HandleFunc("/spotify/player/options", spotifyPlayerOptions)
	http.HandleFunc("/spotify/player/select", spotifyPlayerSelect)
	http.HandleFunc("/do", doHandler)

	log.Println("starting server")
	fmt.Println(http.ListenAndServe(":80", nil))
}

func spotifyAuth(w http.ResponseWriter, req *http.Request) {
	url, err := spotify.GetAuthURL()
	if err != nil {
		log.Panicln(err)
		fmt.Fprint(w, errorText)
		return
	}

	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func spotifyCallback(w http.ResponseWriter, req *http.Request) {
	err := spotify.ProcessCallback(req.URL.Query())
	if err != nil {
		log.Print("failed to process Spotify callback:", err)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func spotifyPlayerOptions(w http.ResponseWriter, req *http.Request) {
	devices, err := spotify.GetPlayers()
	if err != nil {
		log.Panicln(err)
		fmt.Fprint(w, errorText)
		return
	}

	items := []web.Item{}
	for _, device := range devices {
		if device.Restricted {
			continue
		}
		items = append(items,
			web.Item{Text: fmt.Sprintf("%v (%v)", device.Name, device.Type), Path: fmt.Sprintf("/spotify/player/select?id=%v", device.ID)},
		)
	}

	web.Show(w, web.Page{
		Title:      "Choose a player below...",
		Questions:  items,
		ShowButton: false,
	})
}

func spotifyPlayerSelect(w http.ResponseWriter, req *http.Request) {
	spotify.ProcessPlayer(req.URL.Query().Get("id"))

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func doHandler(w http.ResponseWriter, req *http.Request) {
	var image string
	var err error
	if skipCamera {
		image, err = camera.OpenImage("file2.jpg")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		image, err = camera.Snap()
		if err != nil {
			log.Fatal(err)
		}
	}

	var text string
	if skipImageSearch {
		text = "parachutes coldplay"
	} else {
		text, err = vision.Search(image)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("image search result", text)
	}

	web.Show(w, web.Page{
		Title:      fmt.Sprintln("Found: ", text),
		Questions:  []web.Item{},
		ShowButton: true,
	})

	var result string
	if skipSpotifySearch {
		result = "spotify:album:6ZG5lRT77aJ3btmArcykra"
	} else {
		result, err = spotify.SearchAlbum(text)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, errorText)
			return
		}
		log.Println("Spotify result: ", result)
	}

	err = spotify.PlayItem(result)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, errorText)
		return
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	todo := []web.Item{}

	if !spotify.IsAuthed() {
		todo = append(todo, web.Item{Text: spotifyAuthText, Path: "/spotify/auth"})
	} else {
		if !spotify.HasPlayer() {
			todo = append(todo, web.Item{Text: spotifyPlayerText, Path: "/spotify/player/options"})
		}
	}

	if len(todo) > 0 {
		web.Show(w, web.Page{
			Title:      "We need to sort out some stuff...",
			Questions:  todo,
			ShowButton: false,
		})
		return
	}

	web.Show(w, web.Page{
		Title:      "You're good to go!",
		Questions:  []web.Item{},
		ShowButton: true,
	})
}
