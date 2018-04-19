package golang_slacktools

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"strings"
	"time"
)

type SlashCommand struct {
	Token       string `schema:"token"`
	TeamID      string `schema:"team_id"`
	TeamDomain  string `schema:"team_domain"`
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
	Command     string `schema:"command"`
	Text        string `schema:"text"`
	ResponseURL string `schema:"response_url"`
}

type SlashCommandResponse struct {
	Response_type string       `json:"response_type"`
	Text          string       `json:"text"`
	Attachments   []Attachment `json:"attachments"`
}

func GetSimpleSlashCommandResponse(text string) SlashCommandResponse {
	return SlashCommandResponse{Response_type: "in_channel", Text: text}
}

func (resp *SlashCommandResponse) AddAttachment(add Attachment) {
	resp.Attachments = append(resp.Attachments, add)
}

func HandleSlashCommand(respond func(s SlashCommand) SlashCommandResponse, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var s SlashCommand
		switch r.Method {
		case "POST":
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Parameters cannot be parsed", http.StatusBadRequest)
				return
			}
			err = schema.NewDecoder().Decode(&s, r.PostForm)
			if err != nil {
				http.Error(w, "Decoding parameters failed", http.StatusInternalServerError)
				return
			}
		case "GET":
			err := schema.NewDecoder().Decode(&s, r.URL.Query())
			if err != nil {
				http.Error(w, "Decoding parameters failed", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Unsupported http method", http.StatusBadRequest)
		}
		if s.Token == "" || s.TeamID == "" || s.TeamDomain == "" || s.ChannelID == "" || s.ChannelName == "" || s.UserID == "" || s.UserName == "" || s.Command == "" || s.ResponseURL == "" {
			http.Error(w, "Wrong format. Some parameters are empty", http.StatusBadRequest)
			return
		}
		if s.Token != token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		t := time.Now()
		end := false
		go func() {
			d, _ := time.ParseDuration("200ms")
			for time.Since(t).Seconds() < 2.8 && !end {
				time.Sleep(d)
			}
			if !end {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("Command recieved. Response has delay"))
				if err != nil {
					http.Error(w, "Writing a response failed", http.StatusInternalServerError)
					return
				}
				fl, ok := w.(http.Flusher)
				if !ok {
					http.Error(w, "Flushing failed", http.StatusInternalServerError)
					return
				}
				fl.Flush()
			}
		}()

		resp := respond(s)
		end = true
		if time.Since(t).Seconds() < 3 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, "Server failed to encode json", http.StatusInternalServerError)
			}
		} else {
			marshaled, err := json.Marshal(resp)
			if err != nil {
				log.Printf("Marshalling delayed response failed: %q", err)
				return
			}
			req, err := http.NewRequest("POST", s.ResponseURL, strings.NewReader(string(marshaled[:])))
			if err != nil {
				log.Printf("Creating new HTTP-request for a delayed response failed: %q", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Sending delayed response failed: %q", err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("Response is not status ok")
				return
			}
		}
	})
}
