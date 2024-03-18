package spotify

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	logg "go-spooty/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	Client     string
	Key        string
	Refresh    string
	Token      string
	State      string
	httpClient *http.Client
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

type Song struct {
	Playing bool   `json:"is_playing"`
	Shuffle bool   `json:"shuffle_state"`
	Repeat  string `json:"repeat_state"`
	Item    struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		Duration  int       `json:"duration_ms"`
		API_href  string    `json:"href"`
		Preview   string    `json:"preview_url"`
		Track     int       `json:"track_number"`
		URI       string    `json:"uri"`
		Album     Album     `json:"album"`
		Artists   []Artists `json:"artists"`
		Available []string  `json:"available_markets"`
		External  External  `json:"external_urls"`
	}
}

type Album struct {
	Type       string    `json:"album_type"`
	Artists    []Artists `json:"artists"`
	Available  []string  `json:"available_markets"`
	External   External  `json:"external_urls"`
	API_href   string    `json:"href"`
	Id         string    `json:"id"`
	Images     []Images  `json:"images"`
	Album_Name string    `json:"name"`
	Release    string    `json:"release_date"`
}

type Artists struct {
	External External `json:"external_urls"`
	API_href string   `json:"href"`
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	URI      string   `json:"uri"`
}

type External struct {
	URL string `json:"spotify"`
}

type Images struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func init() {
	httpClient = createHTTPClient()
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

/* GET API TOKEN */
func GetAccessToken() (access_token string) {
	param := url.Values{}
	param.Add("grant_type", "refresh_token")
	param.Add("refresh_token", Refresh)
	req, err := http.NewRequest("POST", Token, strings.NewReader(param.Encode()))
	if err != nil {
		logg.SystemLogger.Fatalf("Error Occured. %+v", err)
	}

	str := fmt.Sprintf("%s:%s", Client, Key)
	hdr := fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(str)))
	req.Header.Set("Authorization", hdr)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		logg.SystemLogger.Fatalf("Error sending request to API endpoint. %+v", err)
	} else {
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			logg.SystemLogger.Fatalf("Couldn't parse response body. %+v", err)
		}

		getJson := string(body)
		var result map[string]string
		json.Unmarshal([]byte(getJson), &result)

		return result["access_token"]
	}
	return
}

/* GET PLAYBACK STATE */
func GetPlaybackState() (query string) {
	req, err := http.NewRequest("GET", State, nil)
	if err != nil {
		logg.SystemLogger.Fatalf("Error Occured. %+v", err)
	}

	token := GetAccessToken()
	hdr := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", hdr)

	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		logg.SystemLogger.Fatalf("Error sending request to API endpoint. %+v", err)
	} else {
		defer response.Body.Close()

		switch status := response.StatusCode; status {
		case 503:
			logg.SystemLogger.Println("Service Unavailable")
		case 502:
			logg.SystemLogger.Println("Bad Gateway")
		case 500:
			logg.SystemLogger.Println("Internal Server Error")
		case 429:
			logg.SystemLogger.Println("Too Many Requests")
		case 404:
			logg.SystemLogger.Println("Not Found")
		case 403:
			logg.SystemLogger.Println("Forbidden")
		case 401:
			logg.SystemLogger.Println("Unauthorized")
		case 400:
			logg.SystemLogger.Println("Bad Request")
		case 204:
			logg.SystemLogger.Println("No Content")
		case 200:
			body, err := io.ReadAll(response.Body)
			if err != nil {
				logg.SystemLogger.Fatalf("Couldn't parse response body. %+v", err)
			}

			var song Song
			json.Unmarshal([]byte(body), &song)

			/*
					This block created for debugging Spotify API Request
				  If you need change/update usable fields, use this LoggBlock
			*/

			// *  *	 *  Start Debbuging Block  *	*	 *
			// logg.SystemLogger.Println("Id:", song.Item.Id)
			// logg.SystemLogger.Println("Name:", song.Item.Name)
			// logg.SystemLogger.Println("Duration:", song.Item.Duration)
			// logg.SystemLogger.Println("API_href:", song.Item.API_href)
			// logg.SystemLogger.Println("Preview:", song.Item.Preview)
			// logg.SystemLogger.Println("Track:", song.Item.Track)
			// logg.SystemLogger.Println("URI:", song.Item.URI)
			// logg.SystemLogger.Println("Album Type:", song.Item.Album.Type)
			// logg.SystemLogger.Println("Album Artists External:", song.Item.Album.External.URL)
			// for i, art := range song.Item.Album.Artists {
			// 	text := fmt.Sprintf("Album Artist %d:", i+1)
			// 	logg.SystemLogger.Println(text, art)
			// 	logg.SystemLogger.Println("Album Artists Api Href:", art.API_href)
			// 	logg.SystemLogger.Println("Album Artists Id:", art.Id)
			// 	logg.SystemLogger.Println("Album Artists Name:", art.Name)
			// 	logg.SystemLogger.Println("Album Artists Type:", art.Type)
			// 	logg.SystemLogger.Println("Album Artists URI:", art.URI)
			// }
			// logg.SystemLogger.Println("Album Available:", song.Item.Album.Available)
			// logg.SystemLogger.Println("Album External:", song.Item.Album.External.URL)
			// logg.SystemLogger.Println("Album API Href:", song.Item.Album.API_href)
			// logg.SystemLogger.Println("Album Id:", song.Item.Album.Id)
			// for i, img := range song.Item.Album.Images {
			// 	text := fmt.Sprintf("Album Image %d:", i+1)
			// 	logg.SystemLogger.Println(text, img)
			// 	logg.SystemLogger.Println("Album Images URL:", img.URL)
			// 	logg.SystemLogger.Println("Album Images Height:", img.Height)
			// 	logg.SystemLogger.Println("Album Images Width:", img.Width)
			// }
			// logg.SystemLogger.Println("Album Name:", song.Item.Album.Album_Name)
			// logg.SystemLogger.Println("Album Release:", song.Item.Album.Release)
			// logg.SystemLogger.Println("Available:", song.Item.Available)
			// logg.SystemLogger.Println("External:", song.Item.External.URL)
			// *	*	 *  End Debbuging Block  *	*	 *

			if song.Playing {
				return string(body)
			}
		}
	}
	return
}
