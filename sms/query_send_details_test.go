package sms

import (
	"reflect"
	"testing"
	"time"
)

type testQuerySendDetailsHandler struct {
}

var rightQuerySendDetailsRes = QuerySendDetailsResponse{
	response{
		"0F8F57E7-B72B-492A-853F-F0F8A78D4DEE",
		"OK",
		"OK",
	},
	1,
	0,
	SendDetailDTOs{[]SendDetailDTO{
		{PhoneNum: "15300000001", SendStatus: 3, ErrCode: "DELIVRD",
			TemplateCode: "SMS_132940015", Content: "【可乐贩售机】正在使用Go SDK，版本号：v1.0。",
			SendDate: "2018-04-27 14:19:30", ReceiveDate: "2018-04-27 14:19:35", OutID: "123"},
	}},
}

func (h testQuerySendDetailsHandler) DoReq(opts Options) ([]byte, error) {
	var body []byte
	switch opts.Format() {
	case JSON:
		body = []byte(`{"TotalCount":1,"Message":"OK","RequestId":"0F8F57E7-B72B-492A-853F-F0F8A78D4DEE","SmsSendDetailDTOs":{"SmsSendDetailDTO":[{"OutId":"123","SendDate":"2018-04-27 14:19:30","SendStatus":3,"ReceiveDate":"2018-04-27 14:19:35","ErrCode":"DELIVRD","TemplateCode":"SMS_132940015","Content":"【可乐贩售机】正在使用Go SDK，版本号：v1.0。","PhoneNum":"15300000001"}]},"Code":"OK"}`)
	case XML:
		body = []byte(`<?xml version='1.0' encoding='UTF-8'?><QuerySendDetailsResponse><TotalCount>1</TotalCount><Message>OK</Message><RequestId>0F8F57E7-B72B-492A-853F-F0F8A78D4DEE</RequestId><SmsSendDetailDTOs><SmsSendDetailDTO><OutId>123</OutId><SendDate>2018-04-27 14:19:30</SendDate><SendStatus>3</SendStatus><ReceiveDate>2018-04-27 14:19:35</ReceiveDate><ErrCode>DELIVRD</ErrCode><TemplateCode>SMS_132940015</TemplateCode><Content>【可乐贩售机】正在使用Go SDK，版本号：v1.0。</Content><PhoneNum>15300000001</PhoneNum></SmsSendDetailDTO></SmsSendDetailDTOs><Code>OK</Code></QuerySendDetailsResponse>`)
	}
	return body, nil
}

func testQuerySendDetailsActionDo(t *testing.T, rightURL string, extOpts ...Option) {
	extOpts = append(extOpts, SignatureNonce(u4), Timestamp(ts), ReqHandlerOption(testQuerySendDetailsHandler{}))

	a := NewQuerySendDetailsAction(c, QuerySendDetailsParams{
		RegionID:    "cn-hangzhou",
		PhoneNumber: "15300000001",
		SendDate:    Date(ts),
	})

	opts, err := a.Do(extOpts...)
	if err != nil {
		t.Errorf("Do \"QuerySendDetails\" action err: %v", err)
	}
	if opts.URL() != rightURL {
		t.Errorf("URL: %s != %s", opts.URL(), rightURL)
	}

	res := *opts.Response()
	if !reflect.DeepEqual(res, rightQuerySendDetailsRes) {
		t.Errorf("Response: %v != %v", res, rightQuerySendDetailsRes)
	}
}

func TestQuerySendDetailsAction_Do(t *testing.T) {
	// JSON
	testQuerySendDetailsActionDo(t, "http://dysmsapi.aliyuncs.com/?Signature=IHO%2FUSQcgVW7sWYWoSvCr9%2FoQlI%3D&AccessKeyId=testId&Action=QuerySendDetails&CurrentPage=1&Format=JSON&PageSize=50&PhoneNumber=15300000001&RegionId=cn-hangzhou&SendDate=20180409&SignatureMethod=HMAC-SHA1&SignatureNonce=57d1303b-0068-4892-994d-c2d70d4c37c6&SignatureVersion=1.0&Timestamp=2018-04-09T15%3A27%3A02Z&Version=2017-05-25")

	// XML
	testQuerySendDetailsActionDo(t, "http://dysmsapi.aliyuncs.com/?Signature=fOzO5rT5V8qIY6Td4EMlwm2AtkE%3D&AccessKeyId=testId&Action=QuerySendDetails&CurrentPage=1&Format=XML&PageSize=50&PhoneNumber=15300000001&RegionId=cn-hangzhou&SendDate=20180409&SignatureMethod=HMAC-SHA1&SignatureNonce=57d1303b-0068-4892-994d-c2d70d4c37c6&SignatureVersion=1.0&Timestamp=2018-04-09T15%3A27%3A02Z&Version=2017-05-25",
		XML)
}

func TestDate_String(t *testing.T) {
	d, _ := time.Parse("20060102", "20180502")
	if ds := Date(d).String(); ds != "20180502" {
		t.Errorf("Date string: %s != %s", ds, "20180502")
	}
}

// Use test request Handler, no network latency
func BenchmarkQuerySendDetailsAction_Do(b *testing.B) {
	a := NewQuerySendDetailsAction(c, QuerySendDetailsParams{
		RegionID:    "cn-hangzhou",
		PhoneNumber: "15300000001",
		SendDate:    Date(ts),
	})

	for i := 0; i < b.N; i++ {
		_, err := a.Do(ReqHandlerOption(testQuerySendDetailsHandler{}))
		if err != nil {
			b.Fatal(err)
		}
	}
}
