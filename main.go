package main

import (
	"embed"
	"log"
	"strings"

	"io/fs"
	"main/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static/dist
var FS embed.FS

func main() {
	gin.SetMode(gin.DebugMode)
	app := gin.Default()
	staticFiles, _ := fs.Sub(FS, "static/dist")
	app.StaticFS("/static", http.FS(staticFiles))
	// app.StaticFile("/static", "./frontend/dist/index.html")
	app.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("404.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			ctx.DataFromReader(200, stat.Size(), "text/html", reader, nil)
		}

	})
	api := app.Group("/api")
	{
		api.GET("index", router.GetIndex)
	}
	app.Run(":4999")
}
