package sms

import (
	"reflect"
	"testing"
)

type testSendHandler struct {
}

var rightSendSmsRes = SendSmsResponse{
	Response{
		"6EE2B27D-6833-4D5F-9B9B-CE7FA0A85CC7",
		"OK",
		"OK"},
	"199303724724900469^0",
}

func (h testSendHandler) DoReq(opts Options) ([]byte, error) {
	var body []byte
	switch opts.Format() {
	case JSON:
		body = []byte(`{"Message":"OK","RequestId":"6EE2B27D-6833-4D5F-9B9B-CE7FA0A85CC7","BizId":"199303724724900469^0","Code":"OK"}`)
	case XML:
		body = []byte(`<?xml version='1.0' encoding='UTF-8'?><SendSmsResponse><Message>OK</Message><RequestId>6EE2B27D-6833-4D5F-9B9B-CE7FA0A85CC7</RequestId><BizId>199303724724900469^0</BizId><Code>OK</Code></SendSmsResponse>`)
	}
	return body, nil
}

func testSendActionDo(t *testing.T, rightUrl string, templateParam TemplateParam, outId string, extOpts ...Option) {
	extOpts = append(extOpts, SignatureNonce(u4), Timestamp(ts), ReqHandlerOption(testSendHandler{}))

	a := NewSendAction(c, SendSmsParams{
		"cn-hangzhou",
		"15300000001",
		"阿里云短信测试专用",
		"SMS_71390007",
		templateParam,
		outId})
	opts, err := a.Do(extOpts...)
	if err != nil {
		t.Errorf("Do \"SendSms\" action err: %v", err)
	}
	if opts.Url() != rightUrl {
		t.Errorf("Url: %s != %s", opts.Url(), rightUrl)
	}

	res := *opts.Response()
	if !reflect.DeepEqual(res, rightSendSmsRes) {
		t.Errorf("Response: %v != %v", res, rightSendSmsRes)
	}
}

func TestSendAction_Do(t *testing.T) {
	// all params exist
	// JSON
	testSendActionDo(t,
		"http://dysmsapi.aliyuncs.com/?Signature=gr6VTI2L7pboVdzhg6m96zGfofw%3D&AccessKeyId=testId&Action=SendSms&Format=JSON&OutId=123&PhoneNumbers=15300000001&RegionId=cn-hangzhou&SignName=%E9%98%BF%E9%87%8C%E4%BA%91%E7%9F%AD%E4%BF%A1%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8&SignatureMethod=HMAC-SHA1&SignatureNonce=57d1303b-0068-4892-994d-c2d70d4c37c6&SignatureVersion=1.0&TemplateCode=SMS_71390007&TemplateParam=%7B%22customer%22%3A%22test%22%7D&Timestamp=2018-04-09T15%3A27%3A02Z&Version=2017-05-25",
		templateParam, outId)

	// XML
	testSendActionDo(t,
		"http://dysmsapi.aliyuncs.com/?Signature=IjPuuQwDI864Lsn2ccnzcyOvKEs%3D&AccessKeyId=testId&Action=SendSms&Format=XML&OutId=123&PhoneNumbers=15300000001&RegionId=cn-hangzhou&SignName=%E9%98%BF%E9%87%8C%E4%BA%91%E7%9F%AD%E4%BF%A1%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8&SignatureMethod=HMAC-SHA1&SignatureNonce=57d1303b-0068-4892-994d-c2d70d4c37c6&SignatureVersion=1.0&TemplateCode=SMS_71390007&TemplateParam=%7B%22customer%22%3A%22test%22%7D&Timestamp=2018-04-09T15%3A27%3A02Z&Version=2017-05-25",
		templateParam, outId, XML)

	// omit optional params
	// JSON
	testSendActionDo(t,
		"http://dysmsapi.aliyuncs.com/?Signature=HwBmFIGbv22re%2F3vqdvAxYFqSp0%3D&AccessKeyId=testId&Action=SendSms&Format=JSON&PhoneNumbers=15300000001&RegionId=cn-hangzhou&SignName=%E9%98%BF%E9%87%8C%E4%BA%91%E7%9F%AD%E4%BF%A1%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8&SignatureMethod=HMAC-SHA1&SignatureNonce=57d1303b-0068-4892-994d-c2d70d4c37c6&SignatureVersion=1.0&TemplateCode=SMS_71390007&Timestamp=2018-04-09T15%3A27%3A02Z&Version=2017-05-25",
		nil, "")

	// XML
	testSendActionDo(t,
		"http://dysmsapi.aliyuncs.com/?Signature=gw%2BvEFcdCGYFwxPh7qGab6IoY64%3D&AccessKeyId=testId&Action=SendSms&Format=XML&PhoneNumbers=15300000001&RegionId=cn-hangzhou&SignName=%E9%98%BF%E9%87%8C%E4%BA%91%E7%9F%AD%E4%BF%A1%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8&SignatureMethod=HMAC-SHA1&SignatureNonce=57d1303b-0068-4892-994d-c2d70d4c37c6&SignatureVersion=1.0&TemplateCode=SMS_71390007&Timestamp=2018-04-09T15%3A27%3A02Z&Version=2017-05-25",
		nil, "", XML)
}

func TestTemplateParam_String(t *testing.T) {
	data := TemplateParam{"version": "v1.0"}
	if ds := data.String(); ds != `{"version":"v1.0"}` {
		t.Errorf("TemplateParam string: %s != %s", data.String(), `{"version":"v1.0"}`)
	}
}

// Use test request Handler, no network latency
func BenchmarkSendAction_Do(b *testing.B) {
	a := NewSendAction(c, SendSmsParams{
		"cn-hangzhou",
		"15300000001",
		"阿里云短信测试专用",
		"SMS_71390007",
		templateParam,
		outId})

	for i := 0; i < b.N; i++ {
		_, err := a.Do(ReqHandlerOption(testSendHandler{}))
		if err != nil {
			b.Fatal(err)
		}
	}
}
