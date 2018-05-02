package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// ActionType is type of business param "Action"
type ActionType = string

// FormatType is type of system param "Format"
type FormatType string

// SignatureMethod is type of system param "SignatureMethod"
type SignatureMethod string

// Timestamp is type of system param "Timestamp"
type Timestamp time.Time

// SignatureNonce is type of system param "SignatureNonce"
type SignatureNonce uuid.UUID

type reqHandlerOption struct {
	handler ReqHandler
}

// ReqHandlerOption is helper func to set ReqHandler Option
func ReqHandlerOption(handler ReqHandler) Option {
	return reqHandlerOption{handler: handler}
}

// DefaultSignatureVersion "1.0"
const DefaultSignatureVersion = "1.0"

// DefaultEndPoint "http://dysmsapi.aliyuncs.com/"
const DefaultEndPoint = "http://dysmsapi.aliyuncs.com/"

// DefaultVersion "2017-05-25"
const DefaultVersion = "2017-05-25"

// HTTPMethod "GET"
const HTTPMethod = "GET"

const (
	// SendSms is value of business param "Action"
	SendSms = "SendSms"

	// QuerySendDetails is value of business param "Action"
	QuerySendDetails = "QuerySendDetails"
)

const (
	// JSON is value of system param "Format"
	JSON FormatType = "JSON"

	// XML is value of system param "Format"
	XML FormatType = "XML"
)

const (
	// HmacSha1 is value of system param "SignatureMethod"
	HmacSha1 SignatureMethod = "HMAC-SHA1"
)

// String returns the Timestamp of type "2006-01-02T15:04:05Z07:00"
func (ts Timestamp) String() string {
	return time.Time(ts).Format(time.RFC3339)
}

// MarshalJSON implements the encoding.Marshaler interface.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	t := time.Time(ts)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Timestamp.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.RFC3339)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.RFC3339)
	b = append(b, '"')
	return b, nil
}

// ReqHandler for aliyun sms api request
type ReqHandler interface {
	DoReq(opts Options) ([]byte, error)
}

type defaultReqHandler struct{}

func (h defaultReqHandler) DoReq(opts Options) ([]byte, error) {
	resp, err := http.Get(opts.URL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type action interface {
	Client() Client
}

type baseAction struct {
	c              *Client
	businessParams interface{}
	responseType   reflect.Type
	reqHandler     ReqHandler
}

func (a *baseAction) Client() Client {
	return *a.c
}

func (a *baseAction) generateOpts(extOpts ...Option) (*options, error) {
	opts := options{}

	opts.systemParams.AccessKeyID = a.c.conf.AccessKeyID
	opts.accessSecret = a.c.conf.AccessSecret
	opts.endPoint = a.c.conf.Endpoint

	opts.systemParams.Format = JSON
	opts.systemParams.SignatureMethod = HmacSha1
	opts.systemParams.SignatureVersion = DefaultSignatureVersion

	u4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	opts.systemParams.SignatureNonce = SignatureNonce(u4)
	opts.systemParams.Timestamp = Timestamp(time.Now().UTC())

	opts.businessParams = a.businessParams
	opts.reqHandler = a.reqHandler

	for _, opt := range extOpts {
		opt.Apply(&opts)
	}

	opts.res = reflect.New(a.responseType).Interface()

	return &opts, nil
}

func (a *baseAction) doAction(extOpts ...Option) (*options, error) {
	opts, err := a.generateOpts(extOpts...)
	if err != nil {
		return nil, err
	}
	err = opts.generateURL()
	if err != nil {
		return nil, err
	}

	data, err := opts.reqHandler.DoReq(opts)

	err = opts.processResponse(data)
	if err != nil {
		return nil, err
	}
	return opts, nil
}

// Response represents api response of action
type Response struct {
	RequestID string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

// Config of client
type Config struct {
	AccessKeyID  string
	AccessSecret string
	Endpoint     string
}

// Client of aliyun sms
// it's concurrent safe
type Client struct {
	conf *Config
}

// NewClient ini a new sms client
func NewClient(conf Config) Client {
	sc := Client{}
	if ed := conf.Endpoint; ed == "" {
		conf.Endpoint = DefaultEndPoint
	}
	sc.conf = &conf

	return sc
}

type systemParams struct {
	AccessKeyID      string          `param:"AccessKeyId"`
	Timestamp        Timestamp       `param:"Timestamp"`
	Format           FormatType      `param:"Format,omitempty"`
	SignatureMethod  SignatureMethod `param:"SignatureMethod"`
	SignatureVersion string          `param:"SignatureVersion"`
	SignatureNonce   SignatureNonce  `param:"SignatureNonce"`
	Signature        string          `param:"Signature,omitempty"`
}

// Options represent every action's configurations
type Options interface {
	AccessKeyID() string
	Timestamp() Timestamp
	Format() FormatType
	SignatureMethod() SignatureMethod
	SignatureVersion() string
	SignatureNonce() SignatureNonce
	Signature() string

	URL() string
	EndPoint() string
	AccessSecret() string

	SetSignatureNonce(s SignatureNonce)
	SetFormatType(f FormatType)
	SetTimestamp(ts Timestamp)
	SetReqHandler(reqHandler ReqHandler)
}

type options struct {
	systemParams   systemParams
	businessParams interface{}
	accessSecret   string
	endPoint       string

	reqHandler ReqHandler
	res        interface{}
	url        string
}

func (opts *options) SetSignatureNonce(s SignatureNonce) {
	opts.systemParams.SignatureNonce = s
}

func (opts *options) SetFormatType(f FormatType) {
	opts.systemParams.Format = f
}

func (opts *options) SetTimestamp(ts Timestamp) {
	opts.systemParams.Timestamp = ts
}

func (opts *options) SetReqHandler(reqHandler ReqHandler) {
	opts.reqHandler = reqHandler
}

func (opts *options) URL() string {
	return opts.url
}

func (opts *options) EndPoint() string {
	return opts.endPoint
}

func (opts *options) AccessSecret() string {
	return opts.accessSecret
}

func (opts *options) AccessKeyID() string {
	return opts.systemParams.AccessKeyID
}

func (opts *options) Timestamp() Timestamp {
	return opts.systemParams.Timestamp
}

func (opts *options) Format() FormatType {
	return opts.systemParams.Format
}

func (opts *options) SignatureMethod() SignatureMethod {
	return opts.systemParams.SignatureMethod
}

func (opts *options) SignatureVersion() string {
	return opts.systemParams.SignatureVersion
}

func (opts *options) SignatureNonce() SignatureNonce {
	return opts.systemParams.SignatureNonce
}

func (opts *options) Signature() string {
	return opts.systemParams.Signature
}

// Option represents a single config in every action's configuration
type Option interface {
	// Apply this option to Options
	Apply(opts Options)
}

func (s SignatureNonce) String() string {
	return uuid.UUID(s).String()
}

// Apply option SignatureNonce
func (s SignatureNonce) Apply(opts Options) {
	opts.SetSignatureNonce(s)
}

// Apply option FormatType
func (f FormatType) Apply(opts Options) {
	opts.SetFormatType(f)
}

// Apply option Timestamp
func (ts Timestamp) Apply(opts Options) {
	opts.SetTimestamp(ts)
}

// Apply option ReqHandler
func (handlerOpt reqHandlerOption) Apply(opts Options) {
	opts.SetReqHandler(handlerOpt.handler)
}

func (opts *options) generateURL() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()

	sortedQueryString := opts.sortedQueryString()
	opts.sign(sortedQueryString)

	opts.url = opts.endPoint + "?Signature=" + opts.systemParams.Signature + "&" + sortedQueryString

	return nil
}

func (opts *options) sign(sortedQueryString string) {
	stringToSign := HTTPMethod + "&" + specialQueryEscape("/") + "&" + specialQueryEscape(sortedQueryString)

	// The signature method is supposed to be HmacSHA1
	// A switch case is required if there is other methods available
	mac := hmac.New(sha1.New, []byte(opts.accessSecret+"&"))
	mac.Write([]byte(stringToSign))
	signData := mac.Sum(nil)

	opts.systemParams.Signature = specialQueryEscape(base64.StdEncoding.EncodeToString(signData))
}

func (opts *options) sortedQueryString() string {
	data := url.Values{}

	prepareParameters(&data, opts.systemParams, opts.businessParams)

	// data.Encode() encodes the value sorted by key
	return specialURLEncode(data.Encode())
}

func (opts *options) processResponse(data []byte) error {
	var err error
	switch opts.systemParams.Format {
	case XML:
		err = xml.Unmarshal(data, opts.res)
	case JSON:
		err = json.Unmarshal(data, opts.res)
	}
	if err != nil {
		return err
	}
	return nil
}

func prepareParameters(data *url.Values, params ...interface{}) {
	for _, p := range params {
		v := reflect.ValueOf(p)

		if k := v.Kind(); k == reflect.Ptr || k == reflect.Interface {
			v = v.Elem()
		}

		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i)
			param := fieldInfo.Tag.Get("param")

			tag, tagOptions := parseTag(param)

			if tag == "" {
				if k := v.Field(i).Kind(); k == reflect.Ptr || k == reflect.Interface {
					prepareParameters(data, v.Field(i).Elem().Interface())
				}
				continue
			}

			if tagOptions.contains("omitempty") &&
				reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type()).Interface()) {
				continue
			}

			data.Set(tag, fmt.Sprintf("%v", v.Field(i)))
		}
	}
}

func specialQueryEscape(s string) string {
	return specialURLEncode(url.QueryEscape(s))
}

func specialURLEncode(s string) string {
	s = strings.Replace(s, "+", "%20", -1)
	s = strings.Replace(s, "*", "%2A", -1)
	s = strings.Replace(s, "%7E", "~", -1)
	return s
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
