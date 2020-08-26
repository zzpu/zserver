package http

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zzpu/kratos/pkg/log"
	bm "github.com/zzpu/kratos/pkg/net/http/gin"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"zserver/util"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// RootHandler - http handler for handling / path
func RootHandler(c *bm.Context) {
	t := template.New("index").Delims("<<", ">>")
	t, err := t.ParseFiles("/opt/code/zserver/www/index.html")
	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}

	var fileList = make(map[string]interface{})

	fileList["FileList"] = util.Conf.Dir
	//fileList[csrf.TemplateTag] = csrf.Token(r)
	//fileList["token"] = csrf.Token(r)
	t.Execute(c.Writer, fileList)
}

// WSHandler - Websocket handler
func TailHandler(c *bm.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(c.Writer, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	file := c.Query("file")

	log.Info("file=%v", file)
	filenameB, _ := base64.StdEncoding.DecodeString(file)
	filename := string(filenameB)
	// sanitize the file if it is present in the index or not.
	filename = filepath.Clean(filename)
	ok := false
	for _, wFile := range util.Conf.Dir {
		if filename == wFile.Path {
			ok = true
			break
		}
	}

	// If the file is found, only then start tailing the file.
	// This is to prevent arbitrary file access. Otherwise send a 403 status
	// This should take care of stacking of filenames as it would first
	// be searched as a string in the index, if not found then rejected.
	if ok {
		go util.TailFile(conn, filename)
	}
	c.Writer.WriteHeader(http.StatusUnauthorized)
}

// WSHandler - Websocket handler
func FileList(c *bm.Context) {
	c.JSON(util.Conf.Dir, nil)
}
