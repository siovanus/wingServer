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
	"github.com/siovanus/wingServer/manager/ifpool"
	"github.com/siovanus/wingServer/store"
	"os"
	"os/signal"
	"runtime"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/siovanus/wingServer/config"
	"github.com/siovanus/wingServer/http/restful"
	"github.com/siovanus/wingServer/http/service"
	"github.com/siovanus/wingServer/log"
	"github.com/siovanus/wingServer/manager/flashpool"
	"github.com/siovanus/wingServer/manager/governance"
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

	store, err := store.ConnectToDb(servConfig.DatabaseURL)
	if err != nil {
		log.Errorf("store.ConnectToDb error: %s", err)
		return
	}
	defer store.Close()

	sdk := sdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(servConfig.JsonRpcAddress)

	govAddress, err := common.AddressFromHexString(servConfig.GovernanceAddress)
	if err != nil {
		log.Errorf("govAddress common.AddressFromHexString error: %s", err)
		return
	}
	fpAddress, err := common.AddressFromHexString(servConfig.FlashPoolAddress)
	if err != nil {
		log.Errorf("fpAddress common.AddressFromHexString error: %s", err)
		return
	}
	ifAddress, err := common.AddressFromHexString(servConfig.IFPoolAddress)
	if err != nil {
		log.Errorf("ifAddress common.AddressFromHexString error: %s", err)
		return
	}
	oracleAddress, err := common.AddressFromHexString(servConfig.OracleAddress)
	if err != nil {
		log.Errorf("oracleAddress common.AddressFromHexString error: %s", err)
		return
	}
	oscoreOracleAddress, err := common.AddressFromHexString(servConfig.OscoreOracleAddress)
	if err != nil {
		log.Errorf("oscoreOracleAddress common.AddressFromHexString error: %s", err)
		return
	}
	govMgr := governance.NewGovernanceManager(govAddress, servConfig.WingAddress, sdk, servConfig)
	if govMgr == nil {
		log.Errorf("governance manager is nil")
		return
	}
	fpMgr := flashpool.NewFlashPoolManager(fpAddress, oracleAddress, store, servConfig, govMgr, sdk)
	if fpMgr == nil {
		log.Errorf("flashpool manager is nil")
		return
	}
	ifMgr := ifpool.NewIFPoolManager(sdk, ifAddress, oscoreOracleAddress, store, servConfig, govMgr)
	if ifMgr == nil {
		log.Errorf("ifpool manager is nil")
		return
	}
	log.Infof("init svr success")
	serv := service.NewService(sdk, govMgr, fpMgr, ifMgr, store, servConfig)
	serv.AddListeningAddressList()
	restServer := restful.InitRestServer(serv, servConfig.Port)

	go serv.SnapshotDaily()
	go serv.SnapshotMinute()
	go serv.TrackEvent()
	go serv.MonitorAddrDebt()
	go restServer.Start()
	go checkLogFile(logLevel)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Info("Shutting down...")
	serv.Close()
	os.Exit(0)
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
