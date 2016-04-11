package golang_slacktools

import (
	"testing"
	"net/http"
	"strings"
	"net/http/httptest"
	"encoding/json"
	"io/ioutil"
	"io"
	//"time"
)

func TestPOSTHandleSlashCommand(t *testing.T) {
	request, err:=http.NewRequest("POST","localhost:8080",strings.NewReader("token=Aywr3sizXdooZn9WV2O00NHD&team_id=T0001&team_domain=example&channel_id=C2147483705&channel_name=test&user_id=U2147483697&user_name=Steve&command=/weather&text=94070&response_url=https://hooks.slack.com/commands/1234/5678&"))
	request.Header.Add("Content-Type","application/x-www-form-urlencoded")
	if(err!=nil){
		t.Errorf("error in builing a post request to test with: %q", err)
	}
	response := httptest.NewRecorder()
	HandleSlashCommand(response_example,"Aywr3sizXdooZn9WV2O00NHD").ServeHTTP(response,request)
	if(response.Code!=http.StatusOK){
		t.Errorf("Response has not status 200(OK): %d %q",response.Code,response.Body)
	}
	var resp SlashCommandResponse
	body , _ :=ioutil.ReadAll(io.LimitReader(response.Body, 2000000))
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Error("Unmarshalling json body failed")
	}
	if(resp.Response_type!="in_channel" || resp.Text!="exampletext"){
		t.Error("Invalid response")
	}
}

func TestGETHandleSlashCommand(t *testing.T){
	request, err := http.NewRequest("GET","localhost:8080/?token=Aywr3sizXdooZn9WV2O00NHD&team_id=T0001&team_domain=example&channel_id=C2147483705&channel_name=test&user_id=U2147483697&user_name=Steve&command=/weather&text=94070&response_url=https://hooks.slack.com/commands/1234/5678&",strings.NewReader(""))
	if(err!=nil){
		t.Errorf("error in building a get request to test with: %q", err)
	}
	response := httptest.NewRecorder()
	HandleSlashCommand(response_example,"Aywr3sizXdooZn9WV2O00NHD").ServeHTTP(response,request)
	if(response.Code!=http.StatusOK){
		t.Errorf("Response has not status 200(OK): %d %q",response.Code,response.Body)
	}
	var resp SlashCommandResponse
	body , _ :=ioutil.ReadAll(io.LimitReader(response.Body, 2000000))
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Error("Unmarshalling json body failed")
	}
	if(resp.Response_type!="in_channel" || resp.Text!="exampletext"){
		t.Error("Invalid response")
	}

}

func TestSlashAddAttachment(t *testing.T) {
	a:=GetSimpleSlashCommandResponse("exampleresponse")
	a.AddAttachment(GetSimpleAttachment("attachment1"))
	if(len(a.Attachments)==0){
		t.Error("Adding attachments failed")
	}
}

/*func TestHandleSlashCommand(t *testing.T){
	http.Handle("/",handleSlashCommand(responsexample))
	log.Fatal(http.ListenAndServe(":8080",nil))
}*/

func response_example(s SlashCommand) SlashCommandResponse {
	//d, _ := time.ParseDuration("5000ms")
	//time.Sleep(d)
	return GetSimpleSlashCommandResponse("exampletext")
}