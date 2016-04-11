package golang_slacktools

import (
	"net/http"
	"github.com/gorilla/schema"
	"encoding/json"
)

type OutgoingWebhook struct{
	Token		string	`schema:"token"`
	TeamID		string	`schema:"team_id"`
	TeamDomain	string	`schema:"team_domain"`
	ChannelID	string	`schema:"channel_id"`
	ChannelName	string	`schema:"channel_name"`
	TimeStamp	string	`schema:"timestamp"`
	ServiceID	string	`schema:"service_id"`
	UserID		string	`schema:"user_id"`
	UserName	string	`schema:"user_name"`
	Text		string	`schema:"text"`
	TriggerWord	string	`schema:"trigger_word"`
}

type OutgoingWebhookResponse struct{
	Text		string		`json:"text"`
	Username	string		`json:"username"`
	Attachments 	[]Attachment	`json:"attachments"`
}

func (resp *OutgoingWebhookResponse) AddAttachment(add Attachment){
	resp.Attachments=append(resp.Attachments,add)
}

func GetSimpleOutgoingResponse(text string) OutgoingWebhookResponse{
	return OutgoingWebhookResponse{Text:text}
}

func HandleOutgoingWebhook(respond func(o OutgoingWebhook) OutgoingWebhookResponse, token string)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var o OutgoingWebhook
		err := r.ParseForm()
		if (err != nil) {
			http.Error(w, "Parameters cannot be parsed", http.StatusBadRequest)
			return
		}
		err = schema.NewDecoder().Decode(&o, r.PostForm)
		if (err != nil) {
			http.Error(w, "Decoding parameters failed", http.StatusInternalServerError)
			return
		}
		if(o.Token==""||o.TeamID==""||o.TeamDomain==""||o.ChannelID==""||o.ChannelName==""||o.UserID==""||o.UserName==""||o.TriggerWord==""||o.Text==""){
			http.Error(w,"Wrong format. Some parameters are empty",http.StatusBadRequest)
			return
		}
		if(o.Token!=token){
			http.Error(w,"Invalid token",http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(respond(o)); err != nil {
			http.Error(w, "Server failed to encode json", http.StatusInternalServerError)
		}
	})
}
