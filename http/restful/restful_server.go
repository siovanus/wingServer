// Package restful privides restful server router
package restful

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/siovanus/wingServer/http/common"
	"github.com/siovanus/wingServer/log"
)

type ApiServer interface {
	Start() error
	Stop()
}

type handler func(map[string]interface{}) map[string]interface{}

type Action struct {
	sync.RWMutex
	name    string
	handler handler
}

type restServer struct {
	router   *Router
	port     uint64
	listener net.Listener
	server   *http.Server
	postMap  map[string]Action //post method map
	getMap   map[string]Action //get method map
}

//init restful server
func InitRestServer(web Web, port uint64) ApiServer {
	rt := &restServer{
		port: port,
	}

	rt.router = NewRouter()
	rt.getMap = make(map[string]Action)
	rt.postMap = make(map[string]Action)
	rt.registryRestServerAction(web)
	rt.initGetHandler()
	rt.initPostHandler()
	return rt
}

//resigtry handler method
func (this *restServer) registryRestServerAction(web Web) {

	postMethodMap := map[string]Action{
		common.USERFLASHPOOLOVERVIEW: {name: common.ACTION_USERFLASHPOOLOVERVIEW, handler: web.UserFlashPoolOverview},
		common.ASSETPRICE:            {name: common.ACTION_ASSETPRICE, handler: web.AssetPrice},
		common.ASSETPRICELIST:        {name: common.ACTION_ASSETPRICELIST, handler: web.AssetPriceList},
		common.CLAIMWING:             {name: common.ACTION_CLAIMWING, handler: web.ClaimWing},
		common.LIQUIDATIONLIST:       {name: common.ACTION_LIQUIDATIONLIST, handler: web.LiquidationList},

		common.IFPOOLINFO: {name: common.ACTION_IFPOOLINFO, handler: web.IFPoolInfo},
		common.IFHOSTORY:  {name: common.ACTION_IFHISTORY, handler: web.IFHistory},
	}
	getMethodMap := map[string]Action{
		common.FLASHPOOLMARKETDISTRIBUTION: {name: common.ACTION_FLASHPOOLMARKETDISTRIBUTION, handler: web.FlashPoolMarketDistribution},
		common.IFPOOLMARKETDISTRIBUTION:    {name: common.ACTION_IFPOOLMARKETDISTRIBUTION, handler: web.IfPoolMarketDistribution},
		common.POOLDISTRIBUTION:            {name: common.ACTION_POOLDISTRIBUTION, handler: web.PoolDistribution},
		common.GOVBANNEROVERVIEW:           {name: common.ACTION_GOVBANNEROVERVIEW, handler: web.GovBannerOverview},
		common.GOVBANNER:                   {name: common.ACTION_GOVBANNER, handler: web.GovBanner},
		common.RESERVES:                    {name: common.ACTION_RESERVES, handler: web.Reserves},
		common.IFRESERVES:                  {name: common.ACTION_IFRESERVES, handler: web.IfReserves},
		common.FLASHPOOLDETAIL:             {name: common.ACTION_FLASHPOOLDETAIL, handler: web.FlashPoolDetail},
		common.FLASHPOOLBANNER:             {name: common.ACTION_FLASHPOOLBANNER, handler: web.FlashPoolBanner},
		common.FLASHPOOLALLMARKET:          {name: common.ACTION_FLASHPOOLALLMARKET, handler: web.FlashPoolAllMarket},
		common.BORROWADDRESSLIST:           {name: common.ACTION_BORROWADDRESSLIST, handler: web.BorrowAddressList},
		common.WINGAPYS:                    {name: common.ACTION_WINGAPYS, handler: web.WingApys},
	}
	this.postMap = postMethodMap
	this.getMap = getMethodMap
}

//start server
func (this *restServer) Start() error {
	retPort := this.port
	if retPort == 0 {
		log.Fatal("Not configure HttpRestPort port ")
		return nil
	}

	var err error
	this.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", retPort))
	if err != nil {
		log.Fatal("net.Listen: ", err.Error())
		return err
	}
	log.Infof("server start, listen 0.0.0.0:%d", retPort)
	this.server = &http.Server{Handler: this.router}
	err = this.server.Serve(this.listener)

	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
		return err
	}

	return nil
}

func (this *restServer) getPath(url string) string {
	return url
}

func (this *restServer) getUrlParams(r *http.Request) map[string]interface{} {
	values := r.URL.Query()
	params := make(map[string]interface{})
	for name, value := range values {
		params[name] = value[0]
	}
	return params
}

//init get Handler
func (this *restServer) initGetHandler() {
	for k := range this.getMap {
		this.router.Get(k, func(w http.ResponseWriter, r *http.Request) {
			var resp map[string]interface{}
			url := this.getPath(r.URL.Path)
			if h, ok := this.getMap[url]; ok {
				req := this.getUrlParams(r)
				resp = h.handler(req)
				resp["action"] = h.name
			} else {
				resp = PackResponse(INVALID_METHOD)
				resp["action"] = h.name
			}
			this.response(w, resp)
		})
	}
}

//init post Handler
func (this *restServer) initPostHandler() {
	for k := range this.postMap {
		this.router.Post(k, func(w http.ResponseWriter, r *http.Request) {

			body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			var req = make(map[string]interface{})
			var resp map[string]interface{}

			url := this.getPath(r.URL.Path)
			if h, ok := this.postMap[url]; ok {
				if err := json.Unmarshal(body, &req); err == nil {
					resp = h.handler(req)
				} else {
					log.Error("unmarshal body error:", err)
					resp = PackResponse(ILLEGAL_DATAFORMAT)
				}
				resp["action"] = h.name
			} else {
				resp = PackResponse(INVALID_METHOD)
				resp["action"] = h.name
			}
			this.response(w, resp)
		})
	}
	//Options
	for k := range this.postMap {
		this.router.Options(k, func(w http.ResponseWriter, r *http.Request) {
			this.write(w, []byte{})
		})
	}

}
func (this *restServer) write(w http.ResponseWriter, data []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

//response
func (this *restServer) response(w http.ResponseWriter, resp map[string]interface{}) {
	resp["desc"] = ErrMap[resp["error"].(uint32)]
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("HTTP Handle - json.Marshal: %v", err)
		return
	}
	this.write(w, data)
}

//stop restful server
func (this *restServer) Stop() {
	if this.server != nil {
		this.server.Shutdown(context.Background())
		log.Error("Close restful ")
	}
}

//restart server
func (this *restServer) Restart(cmd map[string]interface{}) map[string]interface{} {
	go func() {
		time.Sleep(time.Second)
		this.Stop()
		time.Sleep(time.Second)
		go this.Start()
	}()

	var resp = PackResponse(SUCCESS)
	return resp
}
