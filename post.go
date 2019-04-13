package main

import (
	"github.com/fionera/e621Crawler/api"
	"github.com/sirupsen/logrus"
)

func ProcessPost(post api.Post) {
	logrus.Debugf("Processing Post | %d", post.ID)

	if post.FileURL != "" {
		logrus.Debugf("Found File | Post: %d - %s", post.ID, post.FileURL)
		startDownload(post.ID, post.FileURL, "file")
	}

	if post.PreviewURL != "" {
		logrus.Debugf("Found Preview | Post: %d - %s", post.ID, post.PreviewURL)
		startDownload(post.ID, post.PreviewURL, "preview")

	}

	if post.SampleURL != "" {
		logrus.Debugf("Found Sample | Post: %d - %s", post.ID, post.SampleURL)
		startDownload(post.ID, post.SampleURL, "sample")

	}
}

func startDownload(id int, url, fileType string) {
	err := DownloadFile(id, url, fileType)
	if err != nil {
		logrus.Error(err)
	}
}
