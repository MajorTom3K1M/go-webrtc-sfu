package main

import (
	"flag"
	"html/template"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	addr          = flag.String("addr", ":8080", "http service address")
	indexTemplate = &template.Template{}
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	indexHTML, err := os.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	indexTemplate = template.Must(template.New("").Parse(string(indexHTML)))

	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.Static("static", "./static")
	r.GET("/", func(c *gin.Context) {
		if err := indexTemplate.Execute(c.Writer, "ws://"+c.Request.Host+"/websocket"); err != nil {
			log.Fatal(err)
		}
	})

	hub := newWebSocketHub()
	go hub.run()

	r.GET("/websocket", func(c *gin.Context) {
		handleWebsocket(c.Writer, c.Request, hub)
	})

	log.Fatal(r.Run(*addr))
}
