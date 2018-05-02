package sms

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

func ExampleSendAction() {
	c := NewClient(Config{AccessKeyID: "testId", AccessSecret: "testSecret"})
	tp := map[string]string{"version": "v1.0"}

	a := NewSendAction(c, SendSmsParams{
		RegionID:     "cn-hangzhou",
		PhoneNumbers: "15300000001",
		SignName:     "可乐贩售机", TemplateCode: "SMS_132940015", TemplateParam: tp, OutID: "123",
	})
	// Do the send action
	// default format type is JSON, we use XML here
	opts, err := a.Do(XML)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(opts.Action())
	fmt.Println(reflect.TypeOf(opts.Response()))
	fmt.Println(opts.Version())
	fmt.Println(opts.RegionID())
	fmt.Println(opts.PhoneNumbers())
	fmt.Println(opts.SignName())
	fmt.Println(opts.TemplateCode())
	fmt.Println(opts.TemplateParam())
	fmt.Println(opts.OutID())

	// Output:
	// SendSms
	// *sms.SendSmsResponse
	// 2017-05-25
	// cn-hangzhou
	// 15300000001
	// 可乐贩售机
	// SMS_132940015
	// {"version":"v1.0"}
	// 123
}

func ExampleNewQuerySendDetailsAction() {
	c := NewClient(Config{AccessKeyID: "testId", AccessSecret: "testSecret"})

	// RegionId is optional here
	// in official http api doc, RegionId
	// in action "QuerySendDetails" is not required,
	// however in official java sdk, RegionId is
	// considered as system params
	// id CurrentPage is not specified, default is 1
	// if PageSize is not specified, default is max value 50
	a := NewQuerySendDetailsAction(c, QuerySendDetailsParams{
		RegionID:    "cn-hangzhou",
		PhoneNumber: "15300000001",
		SendDate:    Date(time.Now()),
		CurrentPage: 2,
	})

	opts, err := a.Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(opts.Action())
	fmt.Println(reflect.TypeOf(opts.Response()))
	fmt.Println(opts.Version())
	fmt.Println(opts.RegionID())
	fmt.Println(opts.PhoneNumber())
	fmt.Println(opts.SendDate())
	fmt.Println(opts.CurrentPage())
	fmt.Println(opts.PageSize())

	// Output:
	// QuerySendDetails
	// *sms.QuerySendDetailsResponse
	// 2017-05-25
	// cn-hangzhou
	// 15300000001
	// 20180502
	// 2
	// 50
}

// Helper func for SendDate
// will panic if value in time.Parse("20060102", value) return an non-nil err
func ExampleDateStr() {
	NewQuerySendDetailsAction(c, QuerySendDetailsParams{
		RegionID:    "cn-hangzhou",
		PhoneNumber: "15300000001",
		SendDate:    DateStr("20180502"),
	})
}
