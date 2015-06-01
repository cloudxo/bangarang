package main

import (
	"flag"
	"log"

	_ "github.com/eliothedeman/bangarang/alarm/console"
	_ "github.com/eliothedeman/bangarang/alarm/pd"
	"github.com/eliothedeman/bangarang/api"
	"github.com/eliothedeman/bangarang/config"
	"github.com/eliothedeman/bangarang/pipeline"
	_ "github.com/eliothedeman/bangarang/provider/http"
	_ "github.com/eliothedeman/bangarang/provider/tcp"
)

var (
	confFile = flag.String("conf", "/etc/bangarang/conf.json", "path main config file")
)

func main() {
	flag.Parse()
	log.Println(flag.Arg(0))
	ac, err := config.LoadConfigFile(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	p := pipeline.NewPipeline(ac)
	p.Start()
	apiServer := api.NewServer(8081, p)
	apiServer.Serve()
	<-make(chan struct{})
}
