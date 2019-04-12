package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var startTime = time.Now()
var totalBytes int64
var numDownloaded int64
var exitRequested int32
var worker sync.WaitGroup
var jobs chan int
var currentPost int

func main() {
	var err error

	// Parse arguments
	parseArgs(os.Args)

	if arguments.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting e621 Crawler")
	logrus.Info("  https://github.com/fionera/e621Crawler/")

	_, cancel := context.WithCancel(context.Background())

	jobs = make(chan int, arguments.Concurrency)
	go listenCtrlC(cancel)
	go stats()

	worker.Add(arguments.Concurrency)
	for i := 0; i < arguments.Concurrency; i++ {
		go crawler()
	}

	_, body, err := fasthttp.Get([]byte{}, "https://e621.net/post/index")
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("Visited Index")

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		logrus.Fatal(err)
	}

	var newestId int
	imageMeta := doc.Find("#post-list > div.content-post > div > .thumb:nth-child(1)")
	if postId, exists := imageMeta.Attr("id"); exists {
		postId = postId[1:]
		parsed, _ := strconv.ParseInt(postId, 10, 64)
		logrus.Debugf("Found Newest Image | Post: %d", parsed)

		newestId = int(parsed)
	}

	for currentPost = arguments.StartId; currentPost <= newestId; currentPost++ {
		if atomic.LoadInt32(&exitRequested) == 1 {
			break
		}

		jobs <- currentPost
		logrus.Debugf("Scheduled Id: %d", currentPost)
	}

	close(jobs)
	worker.Wait()

	logrus.Infof("Last scheduled id was %d", currentPost)
}

func listenCtrlC(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	atomic.StoreInt32(&exitRequested, 1)
	cancel()
	_, _ = fmt.Fprintln(os.Stderr, "\nWaiting for downloads to finish...")
	_, _ = fmt.Fprintln(os.Stderr, "Press ^C again to exit instantly.")
	<-c
	_, _ = fmt.Fprintln(os.Stderr, "\nKilled!")
	os.Exit(255)
}

func stats() {
	for range time.NewTicker(time.Second).C {
		total := atomic.LoadInt64(&totalBytes)
		dur := time.Since(startTime).Seconds()

		logrus.WithFields(logrus.Fields{
			"downloads":    numDownloaded,
			"current_post": currentPost,
			"total_bytes":  totalBytes,
			"avg_rate":     fmt.Sprintf("%.0f", float64(total)/dur),
		}).Info("Stats")
	}
}
