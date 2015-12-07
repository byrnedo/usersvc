package main

import (
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase"
	"net/http"
	"fmt"
	"github.com/byrnedo/usersvc/routers"
	"github.com/byrnedo/apibase/db/mongo"
	"github.com/byrnedo/apibase/natsio"
	"time"
	encBson "github.com/maxwellhealth/encryptedbson"
	"github.com/byrnedo/apibase/helpers/env"
)

func init() {

	apibase.Init()

	mongoUrl := env.GetOr("MONGO_URL", apibase.Conf.GetDefaultString("mongo.url", ""))
	Info.Println("Attempting to connect to [" + mongoUrl + "]")

	mongo.Init(mongoUrl, Trace)

	natsOpts := natsio.NewNatsOptions(func(n *natsio.NatsOptions) error {
		n.Url = env.GetOr("NATS_URL", apibase.Conf.GetDefaultString("nats.url", "nats://localhost:4222"))
		Info.Println("Attempting to connect to [" + n.Url + "]")
		n.Timeout = 10 * time.Second
		if appName, err := apibase.Conf.GetString("app-name"); err == nil && len(appName) > 0 {
			n.Name = appName
		} else {
			panic("must set app-name in conf.")
		}

		Trace.Printf("Nats Opts: %+v", n)

		return nil
	})

	natsCon, err := natsOpts.Connect()
	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}

	encryptionKey, err := apibase.Conf.GetString("encryption-key")
	if err != nil {
		panic("Failed to get encryption-key:" + err.Error())
	}
	copy(encBson.EncryptionKey[:], encryptionKey)

	routers.InitMq(natsCon)

	routers.InitWeb()

}

func main() {

	var (
		host string
		port int
		err error
	)

	host = apibase.Conf.GetDefaultString("http.host", "localhost")
	if port, err = env.GetOrInt("PORT", int(apibase.Conf.GetDefaultInt("http.port", 9999))); err != nil {
		panic(err.Error())
	}

	var listenAddr = fmt.Sprintf("%s:%d", host, port)
	Info.Printf("listening on " + listenAddr)
	if err = http.ListenAndServe(listenAddr, nil); err != nil {
		panic("Failed to start server:" + err.Error())
	}
}