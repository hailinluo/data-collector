package utils

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

type Document struct {
	*goquery.Document
	body io.ReadCloser
}

func GetDocument(url string) (*Document, error) {
	// 获取html
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("status code is not 200. status: %d", res.Status)
		res.Body.Close()
		return nil, errors.New(errMsg)
	}

	// 解析html
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return &Document{
		Document: doc,
		body:     res.Body,
	}, nil
}

func (doc *Document) Close() error {
	if doc.body != nil {
		doc.body.Close()
	}
	return nil
}
