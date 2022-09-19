// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/drone/drone-go/plugin/webhook"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// New returns a new webhook extension.
func New(param1, param2 string) webhook.Plugin {
	return &plugin{
		// TODO replace or remove these configuration
		// parameters. They are for demo purposes only.
		param1: param1,
		param2: param2,
	}
}

type plugin struct {
	// TODO replace or remove these configuration
	// parameters. They are for demo purposes only.
	param1 string
	param2 string
}

func GenSign(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
func sendCard(status string, repoName string, repoLink string, commit string, build string) {
	webhookUrl := os.Getenv("PLUGIN_WEBHOOK")
	secret := os.Getenv("PLUGIN_SECRET")
	baseUrl := os.Getenv("PLUGIN_BASE")
	contentType := "application/json"
	currentTime := time.Now().Unix()
	sign, _ := GenSign(secret, currentTime)
	sendData := `
{
	"timestamp": ` + strconv.FormatInt(currentTime, 10) + `,
	"sign": "` + sign + `",
    "msg_type":"interactive",
    "card":{
        "config":{
            "wide_screen_mode":true
        },
        "elements":[
            {
                "tag":"div",
                "fields":[
                    {
                        "is_short":true,
                        "text":{
                            "tag":"lark_md",
                            "content":"**🗳RepoName：**\n[` + repoName + `](` + repoLink + `)"
                        }
                    },
                    {
                        "is_short":true,
                        "text":{
                            "tag":"lark_md",
                            "content":"**📝Status：**\n` + status + `"
                        }
                    }
                ]
            },
            {
                "tag":"action",
                "actions":[
                    {
                        "tag":"button",
                        "text":{
                            "tag":"lark_md",
                            "content":"[drone](` + baseUrl + build + `)"
                        },
                        "type":"primary"
                    },
                    {
                        "tag":"button",
                        "text":{
                            "tag":"lark_md",
                            "content":"[commit](` + commit + `)"
                        },
                        "type":"default"
                    }
                ]
            }
        ]
    }
}
`
	resp, err := http.Post(webhookUrl, contentType, strings.NewReader(sendData))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Get failed with error: ", resp.Status, string(body))
	}
}
func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Event == webhook.EventBuild {
		fmt.Printf("%+v\n", req)
		link := req.Repo.Link
		slug := req.Repo.Slug
		commit := req.Build.Link
		action := req.Action
		build := slug + "/" + strconv.Itoa(int(req.Build.Number))
		if action == webhook.ActionCreated {
			sendCard(action, slug, link, commit, build)
		}
		if action == webhook.ActionUpdated {
			if req.Build.Status == "success" {
				fmt.Println("success")
				sendCard(req.Build.Status, slug, link, commit, build)
			}
			if req.Build.Status == "failure" {
				sendCard(req.Build.Status, slug, link, commit, build)
			}
		}
	}
	return nil
}
