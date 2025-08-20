package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/marekjh/spotifyanalyzer/internal/server"
)

func main() {
	ctx := context.Background()
	s := server.NewServer(ctx)

	s.Engine.Run()
}

func init() {
	data, err := os.ReadFile("../test/env.json")
	if err != nil {
		log.Fatal(err)
	}

	env := struct {
		ClientID     string
		ClientSecret string
	}{}

	err = json.Unmarshal(data, &env)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv("SPOTIFY_ID", env.ClientID)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("SPOTIFY_SECRET", env.ClientSecret)
	if err != nil {
		log.Fatal(err)
	}
}
