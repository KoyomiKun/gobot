package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"

	ctxlog "github.com/Koyomikun/gobot/utils/logger"
	logger "github.com/Koyomikun/gobot/utils/logger/beego"
)

var (
	ip, port string
	logPath  string
)

func main() {
	app := &cli.App{
		Name:    "kbot",
		Version: "v1.0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-path",
				Aliases:     []string{"l"},
				Usage:       "log file path",
				Destination: &logPath,
			},
			&cli.StringFlag{
				Name:        "ip",
				Usage:       "binding ip address",
				Value:       "127.0.0.1",
				Destination: &ip,
			},
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "binding port",
				Value:       "8080",
				Destination: &port,
			},
		},
		Action: func(c *cli.Context) error {
			// init logger
			ctx := context.WithValue(context.Background(), ctxlog.RequestID, ctxlog.GetUUID())
			bLogger := logger.NewLogger(logPath)
			ctxlog.SetLogger(bLogger)

			ctxlog.Infof(ctx, "start %s", c.App.Name)

			// init websocket
			var upgrade = websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			http.HandleFunc("/cqhttp/ws", func(rw http.ResponseWriter, r *http.Request) {
				selfId, err := strconv.Atoi(r.Header.Get("X-Self-ID"))
				if err != nil {
					fmt.Printf("fail to convert self id: %v \n", err)
					return
				}
				role := r.Header.Get("X-Client-Role")
				host := r.Header.Get("Host")
				fmt.Printf("%d %s %s \n", selfId, role, host)
				conn, err := upgrade.Upgrade(rw, r, nil)
				if err != nil {
					fmt.Printf("fail to start connection \n")
				}
				defer conn.Close()
				// read data
				_, data, err := conn.ReadMessage()
				if err != nil {
					fmt.Printf("fail to read data from connection \n")
					return
				}
				fmt.Printf("Read data: %s \n", data)
			})
			return http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil)
		},
	}
	app.Run(os.Args)

}
