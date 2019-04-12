package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
)

func DownloadFile(id int, url string) error {

	code, body, err := fasthttp.Get([]byte{}, url)
	if err != nil {
		return nil
	}

	if code != 200 {
		return fmt.Errorf("HTTP status %d", code)
	}

	folder := filepath.Join(arguments.Output, strconv.Itoa(id))

	err = os.MkdirAll(folder, 0755)
	if err != nil {
		return nil
	}

	fileName := filepath.Join(folder, filepath.Base(url))

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		logrus.Infof("Skipping File " + fileName)
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil
	}

	size, err := file.Write(body)
	if err != nil {
		return nil
	}

	atomic.AddInt64(&totalBytes, int64(size))
	atomic.AddInt64(&numDownloaded, 1)

	return file.Close()
}
