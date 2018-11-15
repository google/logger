/*
Copyright 2018 AstroPay LLC. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package logger offers simple cross platform logging for Windows and Linux.
// Available logging endpoints are event log (Windows), syslog (Linux), and
// an io.Writer.
package logger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// Slack colors for messages
const (
	ColorGood    = "good"
	ColorDanger  = "danger"
	ColorWarning = "warning"
)

// SendAlert sends a notification to the specified slack channel
func SendAlert(channel, username, title, color, text string) (err error) {

	if channel != "" {

		// control parameters are valid
		if color == "" {
			color = ColorGood
		}

		template := `
		{
			"username": "$USERNAME",
			"attachments": [
				{
					"title": "$TITLE",
					"color": "$COLOR",
					"text": "$TEXT"
				}
			]
		}
		`

		// replace custom data
		msg := strings.Replace(template, "$USERNAME", username, 1)
		msg = strings.Replace(msg, "$TITLE", title, 1)
		msg = strings.Replace(msg, "$COLOR", color, 1)
		msg = strings.Replace(msg, "$TEXT", text, 1)

		// send message using HTTP
		agent := gorequest.New()

		response, _, errPost := agent.Post(channel).Send(msg).End()
		if response != nil {
			defer response.Body.Close()
		}

		if errPost == nil {
			if response.StatusCode != http.StatusOK {
				err = fmt.Errorf("Slack returned status %s", response.Status)
			}
		} else {
			err = fmt.Errorf("Error sending message to Slack: %s", errPost[0])
		}
	} else {
		err = fmt.Errorf("Invalid channel")
	}

	return
}
