package sms

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

func ExampleSendAction() {
	c := NewClient(Config{AccessKeyId: "testId", AccessSecret: "testSecret"})
	tp := map[string]string{"version": "v1.0"}

	a := NewSendAction(c, SendSmsParams{
		RegionId:     "cn-hangzhou",
		PhoneNumbers: "15300000001",
		SignName:     "可乐贩售机", TemplateCode: "SMS_132940015", TemplateParam: tp, OutId: "123",
	})
	// Do the send action
	// default format type is JSON, we use XML here
	opts, err := a.Do(XML)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(opts.Action())
	fmt.Println(reflect.TypeOf(opts.Response()))

	// Output:
	// SendSms
	// *sms.SendSmsResponse
}

func ExampleNewQuerySendDetailsAction() {
	c := NewClient(Config{AccessKeyId: "testId", AccessSecret: "testSecret"})

	// RegionId is optional here
	// in official http api doc, RegionId
	// in action "QuerySendDetails" is not required,
	// however in official java sdk, RegionId is
	// considered as system params
	a := NewQuerySendDetailsAction(c, QuerySendDetailsParams{
		RegionId:    "cn-hangzhou",
		PhoneNumber: "15300000001",
		SendDate:    Date(time.Now()),
	})

	opts, err := a.Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(opts.Action())
	fmt.Println(reflect.TypeOf(opts.Response()))

	// Output:
	// QuerySendDetails
	// *sms.QuerySendDetailsResponse
}
