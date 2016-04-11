package golang_slacktools

import (
	"testing"
)

const Webhook_URL string=""

func TestSendWebhook(t *testing.T) {
	if(Webhook_URL==""){ //in case theres no real slack integration
		t.SkipNow()
	}
	var a IncomingWebhook =GetSimpleWebhook("testhook")

	a.AddAttachment(GetSimpleAttachment("attachment1"))
	a.AddAttachment(GetSimpleAttachment("attachment2"))
	if err:= SendWebhook(a,Webhook_URL);err!=nil{
		t.Error(err)
	}
}

func TestIncAddAttachment(t *testing.T) {
	var a IncomingWebhook =GetSimpleWebhook("testhook")
	a.AddAttachment(GetSimpleAttachment("attachment1"))
	if(len(a.Attachments)==0){
		t.Error("Adding attachments failed")
	}
}