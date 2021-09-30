package setu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	baseUrl string = "https://api.lolicon.app/setu/v2"
)

func init() {
	zero.OnCommand("色图").Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		err := ctx.Parse(&cmd)
		if err != nil {
			log.Errorf("fail parse ctx: %s", ctx.Event.RawMessage)
			return
		}

		if cmd.Args == "" {
			ctx.Send("请附带关键字参数")
			return
		}

		setu := getSetu(cmd.Args)
		if setu == nil {
			return
		}

		ctx.Send(message.At(ctx.Event.UserID))
		ctx.Send(message.Image(setu.Urls["origin"]))
	})
}

type Data struct {
	Pid        int               `json:"pid,omitempty"`
	P          int               `json:"p,omitempty"`
	Uid        int               `json:"uid,omitempty"`
	Title      string            `json:"title,omitempty"`
	Author     string            `json:"author,omitempty"`
	R18        bool              `json:"r18,omitempty"`
	Width      int               `json:"width,omitempty"`
	Height     int               `json:"height,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	Ext        string            `json:"ext,omitempty"`
	UploadData int               `json:"uploadData,omitempty"`
	Urls       map[string]string `json:"urls,omitempty"`
}

type Resp struct {
	Err  string `json:"error,omitempty"`
	Data *Data  `json:"data,omitempty"`
}

func getSetu(keyword string) *Data {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		log.Errorf("fail getting setu: %v", err)
		return nil
	}
	query := req.URL.Query()
	query.Add("keyword", keyword)
	query.Add("r18", "1")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("fail getting setu: %v", err)
		return nil
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	respBody := Resp{}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		log.Errorf("fail unmarshaling setu: %v", err)
		return nil
	}

	if respBody.Err != "" {
		log.Errorf("fail getting setu: %s", respBody.Err)
		return nil
	}
	return respBody.Data
}
