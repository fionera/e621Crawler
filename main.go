package main

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/fionera/e621Crawler/api"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var startTime = time.Now()
var totalBytes int64
var numDownloaded int64
var exitRequested int32
var worker sync.WaitGroup
var jobs chan api.Post
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

	jobs = make(chan api.Post, arguments.Concurrency)
	go listenCtrlC(cancel)
	go stats()

	worker.Add(arguments.Concurrency)
	for i := 0; i < arguments.Concurrency; i++ {
		go crawler()
	}

	currentPost = arguments.StartId
	for {
		if atomic.LoadInt32(&exitRequested) == 1 {
			break
		}

		var posts api.Posts
		err := backoff.Retry(func() error {
			logrus.Infof("Requesting next page before id: %d", currentPost)
			posts, err = api.List(320, currentPost, 0, "", false)
			if err != nil {
				logrus.WithError(err).
					Errorf("Failed crawling")
			}

			return err
		}, backoff.NewExponentialBackOff())

		if err != nil {
			logrus.Error(err)
		}

		for _, post := range posts {
			if atomic.LoadInt32(&exitRequested) == 1 {
				break
			}

			if currentPost == 0 || post.ID < currentPost {
				currentPost = post.ID
			} else if currentPost == post.ID {
				break
			}

			jobs <- post
			logrus.Debugf("Scheduled Id: %d", currentPost)
		}
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

func crawler() {
	defer worker.Done()

	for job := range jobs {
		if atomic.LoadInt32(&exitRequested) == 0 {
			ProcessPost(job)
		}
	}
}
