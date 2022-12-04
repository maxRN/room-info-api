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

func GetRoomsWithWebPage() []RawRooms {

	// ignore these rooms for now, as they don't have a plan: 2530 (E040), 2550 (E042), 2820 (E065), 2860 (E069)
	rooms := []string{"2060", "2100", "2110", "2120", "2130",
		"2140", "2150", "2310", "2660"}
	return getRoomWebPages(rooms)
}

func getRoomNumber(roomId string) (roomNumberInfo RoomNumberResponse) {
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

	roomNumberResponse := RoomNumberResponse{}
	if err := json.Unmarshal(body, &roomNumberResponse); err != nil {
		log.Println("error during marshaling")
		log.Println(err)
	}

	log.Println(roomNumberResponse)

	return roomNumberResponse
}

type RawRooms struct {
	Id       string
	Name     string
	Building string
	WebPage  string
}

func getRoomWebPages(rooms []string) (roomInfos []RawRooms) {
	loginCredentials := LoginCredentials{}

	login(&loginCredentials)
	webPage := ""
	BASE_URL := "https://navigator.tu-dresden.de/raum/542100."

	for _, room := range rooms {
		url := BASE_URL + room
		roomNumberResponse := getRoomNumber(room)
		webPage = fetchWebPage(url, &loginCredentials)
		roomInfos = append(roomInfos, RawRooms{Name: roomNumberResponse.Name, Id: room, Building: "APB", WebPage: webPage})
	}

	return roomInfos
}

type RoomNumberResponse struct {
	RoomId string `json:"roomid"`
	Name   string `json:"name"`
}

type LoginResponse struct {
	LoginToken string `json:"loginToken"`
}

type LoginCredentials struct {
	sessionId  string
	loginToken string
}

type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	University int    `json:"university"`
	From       string `json:"from"`
}

func login(loginCredentials *LoginCredentials) {
	username := base64.StdEncoding.EncodeToString([]byte(os.Getenv("TU_USERNAME")))
	password := base64.StdEncoding.EncodeToString([]byte(os.Getenv("TU_PASSWORD")))

	url := "https://navigator.tu-dresden.de/api/login"
	method := "POST"

	payloadStruct := LoginRequest{Username: strings.TrimSpace(string(username)), Password: strings.TrimSpace(string(password)), University: 1, From: ""}
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
			loginCredentials.sessionId = cookie.Value
		}
	}

	loginResponse := LoginResponse{}

	if err := json.Unmarshal(body, &loginResponse); err != nil {
		log.Println(err)
		log.Println("epic fail")
	}

	loginCredentials.loginToken = loginResponse.LoginToken

}

func fetchWebPage(url string, loginCredentials *LoginCredentials) (webPage string) {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Cookie", "JSESSIONID="+loginCredentials.sessionId+"; loginToken="+loginCredentials.loginToken)

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
			loginCredentials.loginToken = cookie.Value
		}
	}
	return string(body)
}
