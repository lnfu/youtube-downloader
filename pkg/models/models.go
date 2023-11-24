package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Origin struct {
	Id         string
	Title      string
	Duration   int
	InfoStatus string
}

type Media struct {
	Id                 int
	OriginId           string
	Type               string
	MediaStatus        string
	AccessKey          string
	RecentlyAccessTime time.Time
}
