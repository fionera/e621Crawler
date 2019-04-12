package main

import (
	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"sync/atomic"
)

func crawler() {
	defer worker.Done()

	for job := range jobs {
		if atomic.LoadInt32(&exitRequested) == 0 {
			err := backoff.Retry(func() error {
				err := CrawlPost(job)
				if err != nil {
					logrus.WithError(err).
						Errorf("Failed crawling")
				}

				return err
			}, backoff.NewExponentialBackOff())

			if err != nil {
				logrus.Fatal(err)
			}
		}
	}
}
