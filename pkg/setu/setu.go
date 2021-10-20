package setu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	baseUrl       = "https://api.lolicon.app/setu/v2"
	proxy         = "pixiv.runrab.workers.dev"
	msgInfoFormat = "\n%s : %s,"

	sizeKey = "regular"
)

func init() {
	zero.OnCommand("色图").Handle(func(ctx *zero.Ctx) {
		parseSetu(ctx, false)
	})
	zero.OnCommand("R18色图").Handle(func(ctx *zero.Ctx) {
		parseSetu(ctx, true)
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
	Err  string  `json:"error,omitempty"`
	Data []*Data `json:"data,omitempty"`
}

func parseSetu(ctx *zero.Ctx, r18 bool) {

	var cmd extension.CommandModel
	err := ctx.Parse(&cmd)
	if err != nil {
		log.Errorf("fail parse ctx: %s", ctx.Event.RawMessage)
		ctx.Send(message.Text("命令解析出错，请联系管理员"))
		return
	}

	if cmd.Args == "" {
		ctx.Send(message.Text("请附带关键字参数"))
		return
	}

	go ctx.Send(
		message.Message{
			message.At(ctx.Event.UserID),
			message.Text("开始查找..."),
		},
	)

	tags := strings.Split(cmd.Args, " ")
	setuResp := getSetu(tags, r18)
	if setuResp.Err != "" {
		ctx.Send(message.Text(setuResp.Err))
		return
	}
	setu := setuResp.Data
	if setu == nil {
		ctx.Send(message.Text("查找图片失败"))
		return
	}
	if len(setu) == 0 {
		ctx.Send(message.Text("该关键词下搜不到图片"))
		return
	}

	rand.Seed(time.Now().UnixNano())
	topSetu := setu[rand.Intn(len(setu))]

	picUrl, ok := topSetu.Urls[sizeKey]
	if !ok {
		log.Errorf("use key %s in dict %v", sizeKey, topSetu.Urls)
		ctx.Send(message.Text("Key解析失败，请联系管理员处理"))
		return
	}

	var msgInfoBuilder strings.Builder
	msgInfoBuilder.WriteString(fmt.Sprintf(msgInfoFormat, "标题", topSetu.Title))
	msgInfoBuilder.WriteString(fmt.Sprintf(msgInfoFormat, "作者", topSetu.Author))
	msgInfoBuilder.WriteString(fmt.Sprintf(msgInfoFormat, "标签", strings.Join(topSetu.Tags, ",")))

	ctx.Send(message.Message{
		message.At(ctx.Event.UserID),
		message.Text(msgInfoBuilder.String()),
		message.Image(picUrl),
	})
}

func getSetu(tags []string, r18 bool) *Resp {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		log.Errorf("fail getting setu: %v", err)
		return nil
	}

	query := req.URL.Query()
	for _, orTag := range tags {
		query.Add("tag", strings.ReplaceAll(orTag, "或", "|"))
	}
	query.Add("size", sizeKey)
	query.Add("num", "100")
	query.Add("proxy", proxy)
	if r18 {
		query.Add("r18", "1")
	} else {
		query.Add("r18", "0")
	}
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

	// TODO: refactoring <01-10-21, komikun> //
	if len(respBody.Data) == 0 {
		req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
		if err != nil {
			log.Errorf("fail getting setu: %v", err)
			return nil
		}

		query := req.URL.Query()
		query.Add("keyword", tags[0])
		query.Add("size", sizeKey)
		query.Add("num", "100")
		query.Add("proxy", proxy)
		if r18 {
			query.Add("r18", "1")
		} else {
			query.Add("r18", "0")
		}
		req.URL.RawQuery = query.Encode()

		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("fail getting setu: %v", err)
			return nil
		}
		defer resp.Body.Close()

		respBytes, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(respBytes, &respBody)
		if err != nil {
			log.Errorf("fail unmarshaling setu: %v", err)
			return nil
		}
		log.Debugf("From keyword: %s get %d", tags[0], len(respBody.Data))
	}

	return &respBody
}
