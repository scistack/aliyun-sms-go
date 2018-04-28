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

type ActionType = string
type FormatType string
type SignatureMethod string
type Timestamp time.Time
type SignatureNonce uuid.UUID

const DefaultSignatureVersion = "1.0"
const DefaultEndPoint = "http://dysmsapi.aliyuncs.com/"
const DefaultVersion = "2017-05-25"
const HTTPMethod = "GET"

const (
	SendSms          = "SendSms"
	QuerySendDetails = "QuerySendDetails"
)

const (
	JSON FormatType = "JSON"
	XML  FormatType = "XML"
)

const (
	HmacSha1 SignatureMethod = "HMAC-SHA1"
)

func (ts Timestamp) String() string {
	return time.Time(ts).Format(time.RFC3339)
}

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

// handler for aliyun sms api request
type requestHandler interface {
	doReq(opts *options) ([]byte, error)
}

type defaultHandler struct{}

func (h defaultHandler) doReq(opts *options) ([]byte, error) {
	resp, err := http.Get(opts.url)
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
	reqHandler     requestHandler
}

func (a *baseAction) Client() Client {
	return *a.c
}

func (a *baseAction) generateOpts(extOpts ...option) (*options, error) {
	opts := options{}

	opts.systemParams.AccessKeyId = a.c.conf.AccessKeyId
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

	for _, opt := range extOpts {
		opt.apply(&opts)
	}

	opts.res = reflect.New(a.responseType).Interface()

	return &opts, nil
}

func (a *baseAction) doAction(extOpts ...option) (*options, error) {
	opts, err := a.generateOpts(extOpts...)
	if err != nil {
		return nil, err
	}
	err = opts.generateUrl()
	if err != nil {
		return nil, err
	}

	data, err := a.reqHandler.doReq(opts)

	err = opts.processResponse(data)
	if err != nil {
		return nil, err
	}
	return opts, nil
}

type Response struct {
	RequestID string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
}

type Config struct {
	AccessKeyId  string
	AccessSecret string
	Endpoint     string
}

type Client struct {
	conf *Config
}

func NewClient(conf Config) Client {
	sc := Client{}
	if ed := conf.Endpoint; ed == "" {
		conf.Endpoint = DefaultEndPoint
	}
	sc.conf = &conf

	return sc
}

type systemParams struct {
	AccessKeyId      string          `param:"AccessKeyId"`
	Timestamp        Timestamp       `param:"Timestamp"`
	Format           FormatType      `param:"Format,omitempty"`
	SignatureMethod  SignatureMethod `param:"SignatureMethod"`
	SignatureVersion string          `param:"SignatureVersion"`
	SignatureNonce   SignatureNonce  `param:"SignatureNonce"`
	Signature        string          `param:"Signature,omitempty"`
}

type Options interface {
	AccessKeyId() string
	Timestamp() Timestamp
	Format() FormatType
	SignatureMethod() SignatureMethod
	SignatureVersion() string
	SignatureNonce() SignatureNonce
	Signature() string

	Url() string
	EndPoint() string
	AccessSecret() string
}

type options struct {
	systemParams   systemParams
	businessParams interface{}
	accessSecret   string
	endPoint       string

	res interface{}
	url string
}

func (opts *options) Url() string {
	return opts.url
}

func (opts *options) EndPoint() string {
	return opts.endPoint
}

func (opts *options) AccessSecret() string {
	return opts.accessSecret
}

func (opts *options) AccessKeyId() string {
	return opts.systemParams.AccessKeyId
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

type option interface {
	apply(*options)
}

func (s SignatureNonce) String() string {
	return uuid.UUID(s).String()
}

func (s SignatureNonce) apply(opts *options) {
	opts.systemParams.SignatureNonce = s
}

func (f FormatType) apply(opts *options) {
	opts.systemParams.Format = f
}

func (ts Timestamp) apply(opts *options) {
	opts.systemParams.Timestamp = ts
}

func (opts *options) generateUrl() (err error) {
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
	return specialUrlEncode(data.Encode())
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
	return specialUrlEncode(url.QueryEscape(s))
}

func specialUrlEncode(s string) string {
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
