package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"net/http"
	_ "net/http/pprof"

	"github.com/Sirupsen/logrus"
	"github.com/eliothedeman/bangarang/api"
	"github.com/eliothedeman/bangarang/config"
	_ "github.com/eliothedeman/bangarang/escalation/console"
	_ "github.com/eliothedeman/bangarang/escalation/email"
	_ "github.com/eliothedeman/bangarang/escalation/grafana-graphite-annotation"
	_ "github.com/eliothedeman/bangarang/escalation/pd"
	"github.com/eliothedeman/bangarang/pipeline"
	_ "github.com/eliothedeman/bangarang/provider/http"
	_ "github.com/eliothedeman/bangarang/provider/tcp"
)

var (
	confFile    = flag.String("conf", "conf.db", "path main config file")
	dev         = flag.Bool("dev", false, "puts bangarang in a dev testing mode")
	showVersion = flag.Bool("version", false, "display the version of this binary")
	confType    = flag.String("conf-type", "db", `type of configuration used ["db", "json"]`)
	apiPort     = flag.Int("api-port", 8081, "port to serve the http api on")
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	tf := &logrus.TextFormatter{}
	tf.FullTimestamp = true
	tf.ForceColors = true
	logrus.SetFormatter(tf)
}

func handleSigs() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Kill, os.Interrupt)

	done := <-stop
	logrus.Fatal(done.String())
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	flag.Parse()

	// display the current version and exit
	if *showVersion {
		fmt.Print(config.Current)
		os.Exit(0)
	}

	// load configuration
	cp := config.GetProvider(*confType, *confFile)
	if cp == nil {
		logrus.Fatalf("Unable to load config of type %s at location %s", *confType, *confFile)
	}
	ac, err := cp.GetCurrent()
	if err != nil {
		logrus.Fatal(err)
	}

	if ac.LogLevel == "" {
		ac.LogLevel = "info"
	}

	ll, err := logrus.ParseLevel(ac.LogLevel)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.SetLevel(ll)
	}

	logrus.Infof("Starting processing pipeline with %d policie(s)", len(ac.Policies))
	// create and start up a new pipeline
	p := pipeline.NewPipeline(ac)
	p.Start()

	logrus.Infof("Serving the http api on port %d", 8081)
	// create and start a new api server
	apiServer := api.NewServer(*apiPort, p)
	go apiServer.Serve()

	handleSigs()
}
