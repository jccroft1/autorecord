package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jccroft1/autorecord/internal/config"
)

const (
	// Environment names
	envClientID     = "SPOTIFY_CLIENTID"
	envClientSecret = "SPOTIFY_CLIENTSECRET"

	// Config key
	configExpiry       = "spotify_expiry"
	configAccessToken  = "spotify_access_token"
	configRefreshToken = "spotify_refresh_token"
	configPlayer       = "spotify_player"

	// App
	callbackURL = "http://autorecord.local/spotify/callback"

	// Spotify API
	authURL        = "https://accounts.spotify.com/authorize?%v"
	tokenURL       = "https://accounts.spotify.com/api/token"
	searchURL      = "https://api.spotify.com/v1/search?%v"
	listDevicesURL = "https://api.spotify.com/v1/me/player/devices"
	playerURL      = "https://api.spotify.com/v1/me/player/play?device_id=%v"

	defaultTimeFormat = "2006-01-02 15:04:05"
)

var (
	state          = ""
	requiredScopes = []string{
		"user-modify-playback-state", // Start or resume playback
		"user-read-playback-state",   // Get Players
		"user-read-private",          // Search for an item
	}
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	// TokenType string `json:"token_type"`
	// Scope string `json:"scope"`
	Expiry       int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	// TokenType string `json:"token_type"`
	// Scope string `json:"scope"`
	Expiry int `json:"expires_in"`
}

type SearchResponse struct {
	Albums AlbumList `json:"albums"`
	Tracks TrackList `json:"tracks"`
}

type AlbumList struct {
	Items []Album `json:"items"`
}

type TrackList struct {
	Items []Track `json:"items"`
}

type Album struct {
	URI string `json:"uri"`
}

type Track struct {
	URI   string `json:"uri"`
	Album Album  `json:"album"`
}

type PlayRequest struct {
	URI string `json:"context_uri"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID         string `json:"id"`
	Restricted bool   `json:"is_restricted"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
	if os.Getenv(envClientID) == "" || os.Getenv(envClientSecret) == "" {
		panic("unable to load spotify client credentials")
	}
}

func IsAuthed() bool {
	if config.Get(configAccessToken) == "" {
		return false
	}
	return true
}

func HasPlayer() bool {
	if config.Get(configPlayer) == "" {
		return false
	}
	return true
}

func GetAuthURL() (string, error) {
	client_id := os.Getenv(envClientID)
	if client_id == "" {
		return "", fmt.Errorf("Failed to get Spotify client_id from environment")
	}

	state = fmt.Sprint(rand.Int())

	v := url.Values{
		"response_type": []string{"code"},
		"client_id":     []string{client_id},
		"scope":         []string{strings.Join(requiredScopes, " ")},
		"redirect_uri":  []string{callbackURL},
		"state":         []string{state},
	}
	return fmt.Sprintf(authURL, v.Encode()), nil
}

func ProcessCallback(qs url.Values) error {
	if qs.Get("state") != state {
		return fmt.Errorf("state validation failed")
	}

	// extract query string stuff
	code := qs.Get("code")
	if code == "" {
		return fmt.Errorf("Failed to get code: %v", qs.Get("error"))
	}

	requestTime := time.Now()
	tokenData, err := getAuthToken(code)
	if err != nil {
		return err
	}
	if tokenData.AccessToken == "" || tokenData.RefreshToken == "" {
		return fmt.Errorf("spotify tokens not returned")
	}
	config.Set(configExpiry, requestTime.Add(time.Duration(tokenData.Expiry)*time.Second).Format(defaultTimeFormat))
	config.Set(configAccessToken, tokenData.AccessToken)
	config.Set(configRefreshToken, tokenData.RefreshToken)

	return nil
}

func ProcessPlayer(player string) {
	config.Set(configPlayer, player)
}

func GetPlayers() ([]Device, error) {
	err := checkToken()
	if err != nil {
		return []Device{}, err
	}

	req, err := http.NewRequest("GET", listDevicesURL, strings.NewReader(""))
	if err != nil {
		return []Device{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", config.Get(configAccessToken)))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Device{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Device{}, err
	}

	if res.StatusCode != http.StatusOK {
		// refresh token if 401?
		return []Device{}, fmt.Errorf("bad status code response: %v", string(body))
	}

	var data DevicesResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return []Device{}, err
	}

	return data.Devices, nil
}

func SearchAlbum(text string) (string, error) {
	err := checkToken()
	if err != nil {
		return "", err
	}

	qs := url.Values{}
	qs.Set("q", text)
	qs.Set("type", "album,track")
	// qs["type"] = []string{"artist", "album"}
	qs.Set("limit", "5") // TODO: Reduce to 1 eventually?

	url := fmt.Sprintf(searchURL, qs.Encode())
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", config.Get(configAccessToken)))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		// refresh token if 401?
		return "", fmt.Errorf("bad status code response: %v", string(body))
	}

	var data SearchResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if len(data.Albums.Items) > 0 {
		return data.Albums.Items[0].URI, nil
		// TODO: Check list of tracks
	}

	if len(data.Tracks.Items) > 0 {
		return data.Tracks.Items[0].Album.URI, nil
	}

	return "", fmt.Errorf("no results found for %v", text)
}

func PlayItem(uri string) error {
	err := checkToken()
	if err != nil {
		return err
	}
	if !HasPlayer() {
		return fmt.Errorf("no spotify player selected")
	}

	requestBody := PlayRequest{
		URI: uri,
	}
	url := fmt.Sprintf(playerURL, config.Get(configPlayer))
	b, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", config.Get(configAccessToken)))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("bad status code response: %v %v", res.StatusCode, string(body))
	}

	return nil
}

func getAuthToken(code string) (TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", callbackURL)

	var response TokenResponse
	err := getToken(form, &response)
	return response, err
}

func getRefreshToken(token string) (RefreshResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", token)

	var response RefreshResponse
	err := getToken(form, &response)
	return response, err
}

func getToken(form url.Values, data interface{}) error {
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv(envClientID), os.Getenv(envClientSecret))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	return nil
}

func checkToken() error {
	expiryTime, err := time.Parse(defaultTimeFormat, config.Get(configExpiry))
	if err != nil {
		return fmt.Errorf("unable to check expiry: %V", err)
	}

	if expiryTime.After(time.Now()) {
		// token has not yet expired, no need to refresh
		return nil
	}

	log.Println("time expired, refreshing")

	requestTime := time.Now()
	tokenData, err := getRefreshToken(config.Get(configRefreshToken))
	if err != nil {
		return err
	}
	if tokenData.AccessToken == "" {
		return fmt.Errorf("spotify token not returned")
	}
	config.Set(configExpiry, requestTime.Add(time.Duration(tokenData.Expiry)*time.Second).Format(defaultTimeFormat))
	config.Set(configAccessToken, tokenData.AccessToken)

	return nil
}
