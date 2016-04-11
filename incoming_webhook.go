package golang_slacktools

import (
	"encoding/json"
	"strings"
	"net/http"
	"fmt"
)


type IncomingWebhook struct {
	Text		string		`json:"text"`
	Attachments	[]Attachment	`json:"attachments"`
}

type Attachment struct{
	Fallback	string	`json:"fallback"`
	Color		string	`json:"color"`
	Pretext		string	`json:"pretext"`
	AuthorName	string	`json:"author_name"`
	AuthorLink	string	`json:"author_link"`
	AuthorIcon	string	`json:"author_icon"`
	Title		string	`json:"title"`
	TitleLink	string	`json:"title_link"`
	Text		string	`json:"text"`
	Fields		[]Field	`json:"fields"`
	ImageURL	string	`json:"image_url"`
	ThumbURL	string	`json:"thumb_url"`
}

type Field struct {
	Title		string	`json:"title"`
	Value		string	`json:"value"`
	Short		bool	`json:"short"`
}

func GetSimpleWebhook(text string) IncomingWebhook {
	return IncomingWebhook{Text:text}
}

func GetSimpleAttachment(text string) Attachment{
	return Attachment{Text:text}
}

func (hook *IncomingWebhook) AddAttachment (add Attachment){
	hook.Attachments=append(hook.Attachments,add)
}

func SendWebhook(hook IncomingWebhook, url string) error{
	if(url=="" || !strings.HasPrefix(url, "https://hooks.slack.com/services/")){
		return fmt.Errorf("WebhookUrl not configured. The program is unable to send a webhook")
	}
	marshaled, err :=json.Marshal(hook)
	if(err != nil) {
		return fmt.Errorf("Marshalling webhook failed: %q", err)
	}
	req, err :=http.NewRequest("POST",url,strings.NewReader(string(marshaled[:])))
	if(err!=nil) {
		return fmt.Errorf("Creating new HTTP-request for a webhook failed: %q", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp , err :=http.DefaultClient.Do(req)
	if(err!=nil){
		return fmt.Errorf("Sending webhook failed: %q", err)
	}
	if(resp.StatusCode!=http.StatusOK){
		return fmt.Errorf("Response is not status ok")
	}
	return nil
}


