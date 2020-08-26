package http

import (
	"github.com/gorilla/websocket"
	"github.com/zzpu/kratos/pkg/conf/paladin"
	"github.com/zzpu/kratos/pkg/log"
	bm "github.com/zzpu/kratos/pkg/net/http/gin"
	"net"
	"net/http"
	pb "zserver/api"
	"zserver/internal/model"
	"zserver/util"
)

var svc pb.DemoServer

// New new a bm server.
func New(s pb.DemoServer) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	engine = bm.DefaultServer(&cfg)

	var acfg struct {
		Path []string
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&acfg); err != nil {
		return
	}

	util.IndexFiles(acfg.Path)
	pb.RegisterDemoBMServer(engine, s)
	initRouter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	//e.GET("/", RootHandler)
	g := e.Group("/zserver")
	{
		g.GET("/start", howToStart)
		g.GET("/log/list", FileList)
		g.GET("/log/ws", TailHandler)

		g.POST("/login", login)
		g.POST("/logout", logout)
		g.GET("/get_info", getUserInfo)
		g.GET("/ssh", ssh)

	}
}

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func ssh(ctx *bm.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Writer.Header())
	if err != nil {
		log.Error("upgrader.Upgrade,err=%v", err)
		return
	}
	defer conn.Close()

	addr := ctx.Query("addr")

	log.Info("addr=%v", addr)

	tcp, err := net.Dial("tcp", string(addr))
	if err != nil {
		log.Error("net.Dial,err=%v", err)
		return
	}
	defer tcp.Close()
	go func() {
		buf := make([]byte, 1024)
		for {
			len, err := tcp.Read(buf)
			if err != nil {
				log.Error("tcp.Read,err=%v", err)
				conn.Close()
				tcp.Close()
				break
			}
			conn.WriteMessage(websocket.BinaryMessage, buf[0:len])
		}
	}()
	for {
		msgType, buf, err := conn.ReadMessage()
		if err != nil {
			log.Error("conn.ReadMessage,err=%v", err)
			conn.Close()
			tcp.Close()
			break
		}
		if msgType != websocket.BinaryMessage {
			log.Error("unknown msgType")
		}
		tcp.Write(buf)
	}

}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}

	c.JSON(k, nil)
}

// example for http request handler.
func login(c *bm.Context) {
	k := map[string]interface{}{
		"name":    "admin",
		"user_id": 2,
		"access":  "['admin']",
		"token":   "admin",
		"avatar":  "https://avatars0.githubusercontent.com/u/20942571?s=460&v=4",
	}
	c.JSON(k, nil)
}

// example for http request handler.
func getUserInfo(c *bm.Context) {
	k := map[string]interface{}{
		"name":    "admin",
		"user_id": 2,
		"access":  "['admin']",
		"token":   "admin",
		"avatar":  "https://avatars0.githubusercontent.com/u/20942571?s=460&v=4",
	}
	c.JSON(k, nil)
}

func logout(c *bm.Context) {
	c.JSON("logout", nil)
}
