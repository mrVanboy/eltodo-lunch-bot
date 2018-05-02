package webhook

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"
	"fmt"
	"eltodo-lunch-bot/cfg"
)

type Attachment struct {
	Fallback  string `json:"fallback"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	Color	  string `json:"color"`
	Footer	  string `json:"footer"`
	MarkdownIn []string `json:"mrkdwn_in"`
}

type Body struct {
	Username    string `json:"username"`
	Attachments []Attachment `json:"attachments"`
}

var attachments []Attachment

func NewAttachment(attachment Attachment) {
	attachment.MarkdownIn = []string{`text`}
	attachments = append(attachments, attachment)
}

func BuildJSON() ([]byte, error) {
	addColorToAttachments()
	addFooter()
	b:= Body{
		Username: cfg.Get().BotName,
		Attachments: attachments,
	}
	return json.Marshal(b)
}
func addFooter() {
	for i := range attachments {
		attachments[i].Footer = cfg.Get().Footer
	}
}

func addColorToAttachments(){
	colors := getColors(len(attachments))
	for i := range attachments {
		attachments[i].Color = colors[i]
	}
}

func Send() error {
	url := cfg.Get().MainWebhookUrl
	body, err := BuildJSON()
	if err != nil {
		return err
	}
	err = sendPostRequest(url, body)
	if err != nil {
		return err
	}
	return nil
}
func sendPostRequest(url string, body []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		respBody, _ := ioutil.ReadAll(resp.Body)
		return errors.New(`response status code is not 200 OK, but ` + resp.Status + `, response: ` + string(respBody))
	}
	return nil
}

func NotifyAboutErrors(errArr []error) error {
	var data string
	for i, v := range errArr {
		data += fmt.Sprintf("*Error %d:* %s\n", i, v.Error())
	}
	payload := fmt.Sprintf(`{"text": "%s"}`, data)
	bPayload := []byte(payload)
	url := cfg.Get().DebugWebhookUrl

	err := sendPostRequest(url, bPayload)
	if err != nil {
		return err
	}
	return nil
}