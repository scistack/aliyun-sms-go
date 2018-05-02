package sms

import (
	"reflect"
	"time"
)

const (
	// QueryMinPageSize is lower limit page size in api param
	QueryMinPageSize = 1

	// QueryMaxPageSize is upper limit page size in api param
	QueryMaxPageSize = 50
)

// Date is type of business param "Date"
type Date time.Time

// String returns the Date of type "20060102"
func (d Date) String() string {
	return time.Time(d).Format("20060102")
}

// DateStr is a helper func to transform string to type Date
// will panic if time.Parse("20060102", value) returns an err
func DateStr(value string) Date {
	d, err := time.Parse("20060102", value)
	if err != nil {
		panic(err)
	}
	return Date(d)
}

// QuerySendDetailsParams is business param of action "QuerySendDetails"
type QuerySendDetailsParams struct {
	PhoneNumber string `param:"PhoneNumber"`
	BizID       string `param:"BizId,omitempty"`
	SendDate    Date   `param:"SendDate"`
	PageSize    int    `param:"PageSize"`
	CurrentPage int    `param:"CurrentPage"`
	RegionID    string `param:"RegionId,omitempty"`
}

type querySendDetailsParams struct {
	Action  ActionType `param:"Action"`
	Version string     `param:"Version"`
	*QuerySendDetailsParams
}

// QuerySendDetailsOptions represent QuerySendDetailsAction's configurations
type QuerySendDetailsOptions interface {
	Options
	Action() ActionType
	Version() string
	PhoneNumber() string
	BizID() string
	SendDate() Date
	PageSize() int
	CurrentPage() int
	RegionID() string

	Response() *QuerySendDetailsResponse
}

type querySendDetailsOptions struct {
	*options
}

func (q *querySendDetailsOptions) Action() ActionType {
	return q.businessParams.(*querySendDetailsParams).Action
}

func (q *querySendDetailsOptions) Version() string {
	return q.businessParams.(*querySendDetailsParams).Version
}

func (q *querySendDetailsOptions) PhoneNumber() string {
	return q.businessParams.(*querySendDetailsParams).PhoneNumber
}

func (q *querySendDetailsOptions) BizID() string {
	return q.businessParams.(*querySendDetailsParams).BizID
}

func (q *querySendDetailsOptions) SendDate() Date {
	return q.businessParams.(*querySendDetailsParams).SendDate
}

func (q *querySendDetailsOptions) PageSize() int {
	return q.businessParams.(*querySendDetailsParams).PageSize
}

func (q *querySendDetailsOptions) CurrentPage() int {
	return q.businessParams.(*querySendDetailsParams).CurrentPage
}

func (q *querySendDetailsOptions) RegionID() string {
	return q.businessParams.(*querySendDetailsParams).RegionID
}

func (q *querySendDetailsOptions) Response() *QuerySendDetailsResponse {
	return q.res.(*QuerySendDetailsResponse)
}

// QuerySendDetailsAction is action "QuerySendDetails"
type QuerySendDetailsAction interface {
	action
	Do(extOpts ...Option) (QuerySendDetailsOptions, error)
}

type querySendDetailsAction struct {
	baseAction
}

// Do the send action
func (a *querySendDetailsAction) Do(extOpts ...Option) (QuerySendDetailsOptions, error) {
	opts, err := a.baseAction.doAction(extOpts...)
	if err != nil {
		return nil, err
	}
	return &querySendDetailsOptions{opts}, nil
}

func (p *QuerySendDetailsParams) cleanParams() {
	if p.CurrentPage == 0 {
		p.CurrentPage = 1
	}
	if p.PageSize < QueryMinPageSize || p.PageSize > QueryMaxPageSize {
		p.PageSize = QueryMaxPageSize
	}
}

// NewQuerySendDetailsAction init an action "QuerySendDetails"
// can be used concurrently
func NewQuerySendDetailsAction(c Client, params QuerySendDetailsParams) QuerySendDetailsAction {
	params.cleanParams()

	return &querySendDetailsAction{
		baseAction{
			&c,
			&querySendDetailsParams{
				Action:                 QuerySendDetails,
				Version:                DefaultVersion,
				QuerySendDetailsParams: &params,
			},
			reflect.TypeOf(QuerySendDetailsResponse{}),
			defaultReqHandler{},
		},
	}
}

// SendDetailDTO of every sms
type SendDetailDTO struct {
	PhoneNum     string `json:"PhoneNum" xml:"PhoneNum"`
	SendStatus   int    `json:"SendStatus" xml:"SendStatus"`
	ErrCode      string `json:"ErrCode" xml:"ErrCode"`
	TemplateCode string `json:"TemplateCode" xml:"TemplateCode"`
	Content      string `json:"Content" xml:"Content"`
	SendDate     string `json:"SendDate" xml:"SendDate"`
	ReceiveDate  string `json:"ReceiveDate" xml:"ReceiveDate"`
	OutID        string `json:"OutId" xml:"OutId"`
}

// SendDetailDTOs have list of SendDetailDTO
type SendDetailDTOs struct {
	SmsSendDetailDTO []SendDetailDTO `json:"SmsSendDetailDTO" xml:"SmsSendDetailDTO"`
}

// QuerySendDetailsResponse is Response of action "QuerySendDetails"
type QuerySendDetailsResponse struct {
	Response
	TotalCount        int            `json:"TotalCount" xml:"TotalCount"`
	TotalPage         int            `json:"TotalPage" xml:"TotalPage"`
	SmsSendDetailDTOs SendDetailDTOs `json:"SmsSendDetailDTOs" xml:"SmsSendDetailDTOs"`
}
