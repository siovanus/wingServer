/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/restful"
	"github.com/siovanus/wingServer/http/service"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/manager"
	"github.com/urfave/cli"
)

var ConfigPath = ""

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "wing rest server"
	app.Action = startServer
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		config.LogLevelFlag,
		config.ConfigPathFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startServer(ctx *cli.Context) {
	logLevel := ctx.GlobalInt(config.GetFlagName(config.LogLevelFlag))
	log.InitLog(logLevel, log.PATH, log.Stdout)

	configPath := ctx.GlobalString(config.GetFlagName(config.ConfigPathFlag))
	if configPath != "" {
		ConfigPath = configPath
	}
	servConfig, err := config.NewConfig(ConfigPath)
	if err != nil {
		log.Errorf("parse config failed, err: %s", err)
		return
	}

	mgr := manager.NewQueryManager(servConfig)
	if mgr == nil {
		log.Errorf("query manager is nil")
		return
	}
	log.Infof("init svr success")
	serv := service.NewService(mgr)
	restServer := restful.InitRestServer(serv, servConfig.Port)
	go restServer.Start()

	go checkLogFile(logLevel)
	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("server received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}

func checkLogFile(logLevel int) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			isNeedNewFile := log.CheckIfNeedNewFile()
			if isNeedNewFile {
				log.ClosePrintLog()
				log.InitLog(logLevel, log.PATH, log.Stdout)
			}
		}
	}
}
