package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fenthope/gzip"
	"github.com/infinite-iroha/touka"
)

func main() {
	r := touka.Default()
	r.Use(
		func(c *touka.Context) {
			c.Writer.Header().Add("Vary", "Oritouka")
		},
		gzip.Gzip(
			gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{"/ping2"}),
		))

	r.GET("/ping", func(c *touka.Context) {
		c.String(http.StatusOK, "%s", "pong "+fmt.Sprint(time.Now().Unix()))
	})
	r.GET("/ping2", func(c *touka.Context) {
		c.String(http.StatusOK, "%s", "pong "+fmt.Sprint(time.Now().Unix()))
	})
	r.GET("/stream", func(c *touka.Context) {
		c.SetHeader("Content-Type", "text/event-stream")
		c.SetHeader("Connection", "keep-alive")
		for i := 0; i < 10; i++ {
			fmt.Fprintf(c.Writer, "id: %d\ndata: tick %d\n\n", i, time.Now().Unix())
			c.Writer.Flush()
			time.Sleep(1 * time.Second)
		}
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
