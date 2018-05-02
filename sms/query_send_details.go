package sms

import (
	"reflect"
	"time"
)

const (
	QueryMinPageSize = 1
	QueryMaxPageSize = 50
)

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("20060102")
}

func DateStr(value string) Date {
	d, err := time.Parse("20060102", value)
	if err != nil {
		panic(err)
	}
	return Date(d)
}

type QuerySendDetailsParams struct {
	PhoneNumber string `param:"PhoneNumber"`
	BizId       string `param:"BizId,omitempty"`
	SendDate    Date   `param:"SendDate"`
	PageSize    int    `param:"PageSize"`
	CurrentPage int    `param:"CurrentPage"`
	RegionId    string `param:"RegionId,omitempty"`
}

type querySendDetailsParams struct {
	Action  ActionType `param:"Action"`
	Version string     `param:"Version"`
	*QuerySendDetailsParams
}

type QuerySendDetailsOptions interface {
	Options
	Action() ActionType
	Version() string
	PhoneNumber() string
	BizId() string
	SendDate() Date
	PageSize() int
	CurrentPage() int
	RegionId() string

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

func (q *querySendDetailsOptions) BizId() string {
	return q.businessParams.(*querySendDetailsParams).BizId
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

func (q *querySendDetailsOptions) RegionId() string {
	return q.businessParams.(*querySendDetailsParams).RegionId
}

func (q *querySendDetailsOptions) Response() *QuerySendDetailsResponse {
	return q.res.(*QuerySendDetailsResponse)
}

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

// NewQuerySendDetailsAction init a "QuerySendDetails" action
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

type SendDetailDTO struct {
	PhoneNum     string `json:"PhoneNum" xml:"PhoneNum"`
	SendStatus   int    `json:"SendStatus" xml:"SendStatus"`
	ErrCode      string `json:"ErrCode" xml:"ErrCode"`
	TemplateCode string `json:"TemplateCode" xml:"TemplateCode"`
	Content      string `json:"Content" xml:"Content"`
	SendDate     string `json:"SendDate" xml:"SendDate"`
	ReceiveDate  string `json:"ReceiveDate" xml:"ReceiveDate"`
	OutId        string `json:"OutId" xml:"OutId"`
}

type SendDetailDTOs struct {
	SmsSendDetailDTO []SendDetailDTO `json:"SmsSendDetailDTO" xml:"SmsSendDetailDTO"`
}

type QuerySendDetailsResponse struct {
	Response
	TotalCount        int            `json:"TotalCount" xml:"TotalCount"`
	TotalPage         int            `json:"TotalPage" xml:"TotalPage"`
	SmsSendDetailDTOs SendDetailDTOs `json:"SmsSendDetailDTOs" xml:"SmsSendDetailDTOs"`
}
