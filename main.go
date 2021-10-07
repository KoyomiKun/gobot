package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/urfave/cli/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	_ "github.com/Koyomikun/gobot/pkg/hhsh"
	_ "github.com/Koyomikun/gobot/pkg/searcher"
	_ "github.com/Koyomikun/gobot/pkg/setu"
	_ "github.com/Koyomikun/gobot/pkg/texas"
	_ "github.com/Koyomikun/gobot/pkg/tranlate"
)

type Config struct {
	Token      string   `json:"token,omitempty"`
	SuperUsers []string `json:"super_users,omitempty"`
}

var (
	cfgPath string
	addr    string
)

func main() {
	app := cli.App{
		Name:  "gobot",
		Usage: "golang bot based on ZeroBot.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cfg",
				Aliases:     []string{"c"},
				Usage:       "path to json configuration.",
				Value:       "./cfg/config.json",
				Destination: &cfgPath,
			},
			&cli.StringFlag{
				Name:        "addr",
				Aliases:     []string{"a"},
				Usage:       "address the bot listen on.",
				Value:       "127.0.0.1:6700",
				Destination: &addr,
			},
		},
		Action: func(c *cli.Context) error {
			// init logger
			log.SetFormatter(&easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
			})
			log.SetLevel(log.DebugLevel)

			// init config
			fileBytes, err := readFile(cfgPath)
			if err != nil {
				log.Errorf("fail init config: %v", err)
				return err
			}
			config := &Config{}
			err = json.Unmarshal(fileBytes, config)
			if err != nil {
				log.Errorf("fail init config: %v", err)
				return err
			}

			log.Infof("start %s", c.App.Name)

			// run zerobot
			zero.Run(zero.Config{
				NickName:      []string{"bot"},
				CommandPrefix: ".",
				SuperUsers:    config.SuperUsers,
				Driver: []zero.Driver{
					driver.NewWebSocketClient(fmt.Sprintf("ws://%s/", addr), config.Token),
				},
			})
			select {}
		},
	}
	app.Run(os.Args)

}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("fail reading file :%v", err)
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("fail reading file :%v", err)
		return nil, err
	}
	return fileBytes, nil
}
