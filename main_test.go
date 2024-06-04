package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	goLoggerHTTPClient "github.com/pobyzaarif/go-logger/http/client"
	"github.com/pobyzaarif/go-mailer/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// setup test
	config := config.LoadConfig("./config.json")
	conf := GoMailerConfig{
		Provider: config.Provider,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		From:     config.Username,
	}

	mailService := NewGoMailer(conf)
	msgSubject := "This is test mail"
	msgBody := fmt.Sprintf("msg %v", time.Now().Unix())

	err := mailService.SendEmail(GoMailerForm{
		To:      []string{config.TestMailTo},
		Subject: msgSubject,
		Body:    msgBody,
	})
	assert.Equal(t, nil, err)

	var gotResponseFrom, gotResponseBody, gotResponseSubject string
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Second * 5) // add delay, and maybe need some time to deliver an email.
		// construct payload
		rawPayload := map[string]interface{}{"from": config.TestMailTo}
		payload := new(bytes.Buffer)
		err = json.NewEncoder(payload).Encode(rawPayload)
		if err != nil {
			t.Error(err)
		}

		// construct request
		request, err := http.NewRequest(http.MethodPost, config.TestMailReaderURL, payload)
		if err != nil {
			t.Error(err)
		}
		request.Header.Add("Accept", "application/json")
		request.Header.Add("Content-Type", "application/json")

		timeout := time.Second * 15

		type response struct {
			Rc    int    `json:"rc"`
			Error string `json:"error"`
			Data  struct {
				From    string `json:"from"`
				Subject string `json:"subject"`
				Body    string `json:"body"`
			} `json:"data"`
		}
		var resp response
		_, err = goLoggerHTTPClient.Call(
			context.TODO(),
			request,
			timeout,
			goLoggerHTTPClient.JSONResponseBodyFormat,
			&resp,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		t.Logf("read email counter %v", i)
		t.Logf("response: %+v", resp)
		if resp.Data.Body != "" {
			gotResponseFrom = resp.Data.From
			gotResponseSubject = resp.Data.Subject
			gotResponseBody = resp.Data.Body
			break
		}
	}

	assert.Equal(t, config.Username, gotResponseFrom)
	assert.Equal(t, msgSubject, gotResponseSubject)
	assert.Contains(t, gotResponseBody, msgBody)
}
