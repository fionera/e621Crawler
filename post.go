package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strconv"
)

func CrawlPost(id int) error {
	code, body, err := fasthttp.Get([]byte{}, "https://e621.net/post/show/"+strconv.Itoa(id))
	if err != nil {
		return err
	}

	switch code {
	case 200:
		break
	case 404:
		return nil
	case 503:
		return errors.New("503 Bad Gateway")
	default:
		return errors.New(fmt.Sprintf("Unknown StatusCode - %d", code))
	}

	logrus.Infof("Found Post | %d", id)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		logrus.Error(err)
		return err
	}

	foundVideo := false
	foundImage := false

	videoMeta := doc.Find("meta[property=\"og:video\"]")
	if url, exists := videoMeta.Attr("content"); exists {
		foundVideo = true
		logrus.Debugf("Found Video | Post: %d - %s", id, url)
		go startDownload(id, url)
	}

	highRes := doc.Find("#highres")
	if url, exists := highRes.Attr("href"); exists && !foundVideo {
		foundImage = true
		logrus.Debugf("Found High Resolution | Post: %d - %s", id, url)
		go startDownload(id, url)
	}

	imageMeta := doc.Find("meta[property=\"og:image\"]")
	if url, exists := imageMeta.Attr("content"); exists && !foundImage {
		logrus.Debugf("Found Thumbnail | Post: %d - %s", id, url)
		go startDownload(id, url)
	}

	return err
}

func startDownload(id int, url string) {
	err := DownloadFile(id, url)
	if err != nil {
		logrus.Error(err)
	}
}
