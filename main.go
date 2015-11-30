package main

import (
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase"
	"net/http"
	"fmt"
	"github.com/byrnedo/apibase/env"
	"github.com/byrnedo/usersvc/routers"
	"github.com/byrnedo/apibase/db/mongo"
	"github.com/byrnedo/apibase/natsio"
	"time"
)

func main() {

	var (
		host string
		port int
		err error
	)

	apibase.Init()

	mongo.Init(env.GetOr("MONGO_URL", apibase.Conf.GetDefaultString("mongo.url", "")), Trace)

	natsOpts := natsio.NewNatsOptions(func(n *natsio.NatsOptions) error {
		n.Url = env.GetOr("NATS_URL", apibase.Conf.GetDefaultString("nats.url", "nats://localhost:4222"))
		n.Timeout = 10 * time.Second
		return nil
	})

	natsCon, err := natsOpts.Connect()
	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}

	routers.InitMq(natsCon)

	webRouter := routers.InitWeb()

	host = apibase.Conf.GetDefaultString("http.host", "localhost")
	if port, err = env.GetOrInt("PORT", apibase.Conf.GetDefaultInt("http.port", 9999)); err != nil {
		panic(err.Error())
	}

	var listenAddr = fmt.Sprintf("%s:%d", host, port)
	Info.Printf("listening on " + listenAddr)
	http.ListenAndServe(listenAddr, webRouter)
}