package tranlate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	zero.OnCommand("翻译").Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		err := ctx.Parse(&cmd)
		if err != nil {
			log.Errorf("fail parse translate context %v: %v", ctx, err)
			return
		}
		if cmd.Args == "" {
			ctx.Send(message.Message{message.Text("你想翻译个啥")})
			return
		}
		ctx.Send(translate(cmd.Args))
	})
}

const (
	baseUrl string = "http://translate.google.cn/translate_a/single?client=gtx&dt=t&dj=1&ie=UTF-8&sl=auto&tl=%s&q=%s"
)

func translate(origin string) string {
	localeToCh := map[string]string{
		"zh_CN": "中文",
		"ja":    "日语",
		"en_US": "英语",
		"fa":    "波斯语",
		"ru":    "俄语",
	}

	resChan := make(chan string, len(localeToCh))
	defer close(resChan)
	for k, v := range localeToCh {
		go func(url, locale string, resChan chan<- string) {
			var retSb strings.Builder

			client := http.Client{
				Timeout: 2 * time.Second,
			}
			resp, err := client.Get(url)
			if err != nil {
				log.Errorf("fail getting url %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("fail read resp.Body in url %s: %v", url, err)
				return
			}

			bodyMap := map[string]interface{}{}
			err = json.Unmarshal(respBytes, &bodyMap)
			if err != nil {
				log.Errorf("fail unmarshal bodyMap in url %s: %v", url, err)
				return
			}

			senetences, ok := bodyMap["sentences"].([]interface{})
			if !ok {
				log.Errorf("fail get sentences in url %s", url)
				return
			}

			retSb.WriteString(locale)
			retSb.WriteByte(':')
			retSb.WriteString(senetences[0].(map[string]interface{})["trans"].(string))
			resChan <- retSb.String()
		}(fmt.Sprintf(baseUrl, k, origin), v, resChan)
	}

	var retSb strings.Builder
loop:
	for i := 0; i < len(localeToCh); i++ {
		select {
		case res := <-resChan:
			retSb.WriteString(res)
			if i != len(localeToCh)-1 {
				retSb.WriteByte('\n')
			}
		case <-time.After(time.Second * 2):
			log.Warnf("translate word %s over time", origin)
			break loop
		}
	}
	return retSb.String()
}
