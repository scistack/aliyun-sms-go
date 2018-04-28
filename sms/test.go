package sms

import (
	"github.com/satori/go.uuid"
	"time"
)

var c = NewClient(Config{AccessKeyId: "testId", AccessSecret: "testSecret"})
var templateParam = map[string]string{"customer": "test"}
var u4, _ = uuid.FromString("57d1303b-0068-4892-994d-c2d70d4c37c6")
var ts, _ = time.Parse(time.RFC3339, "2018-04-09T15:27:02Z")
var outId = "123"
