package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type rawRoom struct {
	Id       string
	Name     string
	Building string
	WebPage  string
}

// GetAllRooms fetches the raw HTML code for all rooms.
// The rawRooms returned from this function need to then be parsed into propers Rooms.
func GetAllRooms() []rawRoom {

	// ignore these rooms for now, as they don't have a plan: 2530 (E040), 2550 (E042), 2820 (E065), 2860 (E069)
	// roomsToFetch := []string{"2060", "2100", "2110", "2120", "2130",
	// 	"2140", "2150", "2310", "2660"}
	roomsToFetch := []string{"2060"}

	return fetchRooms(roomsToFetch)
}

// Returns the room number for the room with roomId.
// You don't need to be authenticated to make this call.
func getRoomName(roomId string) (roomName string) {
	url := "https://navigator.tu-dresden.de/api/roominfo/542100." + roomId + "?all=true"

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	roomNumberResponse := roomNumberResponse{}
	if err := json.Unmarshal(body, &roomNumberResponse); err != nil {
		log.Println("error during marshaling")
		log.Println(err)
	}

	log.Println(roomNumberResponse)

	return roomNumberResponse.Name
}

// Fetches all rooms specified by rooms.
func fetchRooms(rooms []string) (roomInfos []rawRoom) {
	credentials := loginCredentials{}

	login(&credentials)
	webPage := ""
	BASE_URL := "https://navigator.tu-dresden.de/raum/542100."

	for _, room := range rooms {
		url := BASE_URL + room
		roomNumberResponse := getRoomName(room)
		webPage = fetchWebPage(url, &credentials)
		roomInfos = append(roomInfos, rawRoom{Name: roomNumberResponse, Id: room, Building: "APB", WebPage: webPage})
	}

	return roomInfos
}

type roomNumberResponse struct {
	RoomId string `json:"roomid"`
	Name   string `json:"name"`
}

// Fetches the HTML contents of the webpage at the specified URL.
func fetchWebPage(url string, credentials *loginCredentials) (webPage string) {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Cookie", "JSESSIONID="+credentials.sessionId+"; loginToken="+credentials.loginToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	for _, cookie := range res.Cookies() {
		if cookie.Name == "loginToken" {
			credentials.loginToken = cookie.Value
		}
	}
	return string(body)
}

type loginResponse struct {
	LoginToken string `json:"loginToken"`
}

type loginCredentials struct {
	sessionId  string
	loginToken string
}

func login(credentials *loginCredentials) {
	type loginRequest struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		University int    `json:"university"`
		From       string `json:"from"`
	}

	username := base64.StdEncoding.EncodeToString([]byte(os.Getenv("TU_USERNAME")))
	password := base64.StdEncoding.EncodeToString([]byte(os.Getenv("TU_PASSWORD")))

	url := "https://navigator.tu-dresden.de/api/login"
	method := "POST"

	payloadStruct := loginRequest{Username: strings.TrimSpace(string(username)), Password: strings.TrimSpace(string(password)), University: 1, From: ""}
	payload, _ := json.Marshal(payloadStruct)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "JSESSIONID" {
			credentials.sessionId = cookie.Value
		}
	}

	loginResponse := loginResponse{}

	if err := json.Unmarshal(body, &loginResponse); err != nil {
		log.Println(err)
		log.Println("epic fail")
	}

	credentials.loginToken = loginResponse.LoginToken

}
