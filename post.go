package main

import (
	"github.com/fionera/e621Crawler/api"
	"github.com/sirupsen/logrus"
)

func ProcessPost(post api.Post) {
	logrus.Debugf("Processing Post | %d", post.ID)

	if post.File.URL != "" {
		logrus.Debugf("Found File | Post: %d - %s", post.ID, post.File.URL)
		startDownload(post.ID, post.File.URL, "file")
	}

	if post.Preview.URL != "" {
		logrus.Debugf("Found Preview | Post: %d - %s", post.ID, post.Preview.URL)
		startDownload(post.ID, post.Preview.URL, "preview")

	}

	if post.Sample.URL != "" {
		logrus.Debugf("Found Sample | Post: %d - %s", post.ID, post.Sample.URL)
		startDownload(post.ID, post.Sample.URL, "sample")

	}
}

func startDownload(id int, url, fileType string) {
	err := DownloadFile(id, url, fileType)
	if err != nil {
		logrus.Error(err)
	}
}
