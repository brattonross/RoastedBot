package http

import (
	"encoding/json"
	"net/http"

	"github.com/brattonross/roastedbot/pkg/twitch"
)

// NewHandler creates a new handler for the bot service.
func NewHandler(client *twitch.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/channels", channels(client))
	return mux
}

func channels(client *twitch.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		json.NewEncoder(w).Encode(client.Channels())
	}
}
