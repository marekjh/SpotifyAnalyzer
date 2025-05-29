package main

import "github.com/zmb3/spotify/v2/auth"

const RedirectURL = "https://127.0.0.1/auth-response"

func main() {
	auth := spotifyauth.New(RedirectURL)
}