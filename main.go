package main

import (
	// "crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	// "strconv"
	"strings"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/joho/godotenv"
)

type Track struct {
	Name string `json:"name"`
}
type Item struct {
	Track Track `json:"track"`
}
type RecentlyPlayedResponse struct {
	Items []Item `json:"items"`
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/svg+xml")
	w.Header().Add("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
	client := &http.Client{}
	rcpls := getrecpls(client, gettoken(client))
	// AI code below
	s := svg.New(w)
	width, height := 700, 800
	padding := 40
	lineHeight := 30
	s.Start(width, height)
	s.Rect(0, 0, width, height, "fill:#121212")
	s.Rect(padding/2, padding/2, width-padding, height-padding, "fill:#1db954;stroke:#ffffff;stroke-width:4;rx:24;ry:24")
	// Heading
	s.Text(width/2, padding+30, "Recently Played", "text-anchor:middle;font-family:Arial, sans-serif;font-size:32px;fill:#ffffff;font-weight:bold")
	// Subtitle
	s.Text(width/2, padding+68, "Spotify Banner", "text-anchor:middle;font-family:Arial, sans-serif;font-size:18px;fill:#f3f3f3")

	if len(rcpls) == 0 {
		s.Text(width/2, height/2, "No recent tracks found", "text-anchor:middle;font-family:Arial, sans-serif;font-size:20px;fill:#ffffff")
	} else {
		maxItems := 20
		if len(rcpls) < maxItems {
			maxItems = len(rcpls)
		}
		startY := padding + 120
		for i := 0; i < maxItems; i++ {
			text := rcpls[i]
			if len(text) > 50 {
				text = text[:47] + "..."
			}
			x := padding + 24
			y := startY + i*lineHeight
			bgY := y - 24
			s.Rect(x-16, bgY, width-padding-64, lineHeight+10, "fill:rgba(0,0,0,0.15);rx:10;ry:10")
			s.Text(x, y, fmt.Sprintf("%d. %s", i+1, text), "font-family:Arial, sans-serif;font-size:18px;fill:#ffffff")
		}
	}
	// AI code above

	s.End()
}

func getcodeverifier() string {
	byts := make([]byte, 64)
	_, err := rand.Read(byts)
	if err != nil {
		panic(err)
	}
	token := base64.RawURLEncoding.EncodeToString(byts)
	return token
}

func getcodechallenge(codeverifier string) string {
	shaenc := sha256.Sum256([]byte(codeverifier))
	return base64.RawURLEncoding.EncodeToString(shaenc[:])
}
func getapitoken(client *http.Client, code string, codeverifier string) string {
	baseurl := "https://accounts.spotify.com/api/token"
	params := url.Values{
		"client_id":     {"cfe923b2d660439caf2b557b21f31221"},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://developer.spotify.com"},
		"code_verifier": {codeverifier},
	}
	req, err := http.NewRequest("POST", baseurl, strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	token := strings.Split(strings.Split(string(bodyText), "\"access_token\":\"")[1], "\",")[0]
	return token
}

func gettoken(client *http.Client) string {
	godotenv.Load()
	codever := getcodeverifier()
	codechal := getcodechallenge(codever)

	baseURL := "https://accounts.spotify.com/oauth2/v2/auth"
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {"cfe923b2d660439caf2b557b21f31221"},
		"scope":                 {"email openid profile user-self-provisioning playlist-modify-private playlist-modify-public playlist-read-collaborative playlist-read-private ugc-image-upload user-follow-modify user-follow-read user-library-modify user-library-read user-modify-playback-state user-read-currently-playing user-read-email user-read-playback-position user-read-playback-state user-read-private user-read-recently-played user-top-read user-personalized"},
		"redirect_uri":          {"https://developer.spotify.com"},
		"code_challenge":        {codechal},
		"code_challenge_method": {"S256"},
		"state":                 {"x1X0-Qo96pi118lEr8s0MZlVQ_lgVfCW"},
		"response_mode":         {"web_message"},
		"prompt":                {"none"},
	}
	authUrl := baseURL + "?" + params.Encode()
	req, err := http.NewRequest("GET", authUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	spdccok := os.Getenv("SPDC")
	req.Header.Set("Cookie", "sp_dc="+spdccok)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	code := strings.Split(strings.Split(string(bodyText), "\"code\": \"")[1], "\"")[0]
	token := getapitoken(client, code, codever)
	return token
}

func getrecpls(client *http.Client, token string) []string {
	murl := "https://api.spotify.com/v1/me/player/recently-played"
	params := url.Values{
		"limit": {"20"},
	}
	req, err := http.NewRequest("GET", murl+"?"+params.Encode(), nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var data RecentlyPlayedResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	var relol []string
	for _, item := range data.Items {
		relol = append(relol, item.Track.Name)
	}
	return relol
}

func main() {
	// client := &http.Client{}
	// fmt.Println("hii");
	// tok := gettoken(client)
	// getrecpls(client,tok)
	s := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(myHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
