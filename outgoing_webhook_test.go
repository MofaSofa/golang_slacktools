package golang_slacktools

import (
	"testing"
	"fmt"
	"net/http"
	"strings"
	"net/http/httptest"
)

func TestHandleOutgoingWebhook(t *testing.T){
	request, err:=http.NewRequest("POST","localhost:8080",strings.NewReader("token=aoDvjEaQr9wQ0lqC85plXkMG&team_id=T0001&team_domain=example&channel_id=C2147483705&channel_name=test&timestamp=1355517523.000005&user_id=U2147483697&user_name=Steve&text=googlebot: What is the air-speed velocity of an unladen swallow?&trigger_word=googlebot:"))
	request.Header.Add("Content-Type","application/x-www-form-urlencoded")
	if(err!=nil){
		t.Errorf("error in builing a post request to test with: %q", err)
	}
	response := httptest.NewRecorder()
	HandleOutgoingWebhook(outgoingwebhookresponseexample, "aoDvjEaQr9wQ0lqC85plXkMG").ServeHTTP(response,request)
	if(response.Code!=http.StatusOK){
		t.Errorf("Response has not status 200(OK): %d %q",response.Code,response.Body)
	}
}

func TestOutAddAttachment(t *testing.T) {
	a:=GetSimpleOutgoingResponse("exampleresponse")
	a.AddAttachment(GetSimpleAttachment("attachment1"))
	if(len(a.Attachments)==0){
		t.Error("Adding attachments failed")
	}
}

func outgoingwebhookresponseexample(o OutgoingWebhook) OutgoingWebhookResponse{
	b :=OutgoingWebhookResponse{Text:fmt.Sprintf("%q %q %q",o.Token,o.TimeStamp,o.Text),Username:o.Text}
	b.Attachments=append(b.Attachments,GetSimpleAttachment("bla"))
	return b
}