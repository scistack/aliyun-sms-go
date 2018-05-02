package sms

import (
	"encoding/json"
	"reflect"
)

// TemplateParam is type of business param "TemplateParam"
type TemplateParam map[string]string

// String returns the TemplateParam of type JSON string
func (tp TemplateParam) String() string {
	data, err := json.Marshal(tp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// SendSmsParams is business param of action "SendSms"
type SendSmsParams struct {
	RegionID      string        `param:"RegionId"`
	PhoneNumbers  string        `param:"PhoneNumbers"`
	SignName      string        `param:"SignName"`
	TemplateCode  string        `param:"TemplateCode"`
	TemplateParam TemplateParam `param:"TemplateParam,omitempty"`
	OutID         string        `param:"OutId,omitempty"`
}

type sendSmsParams struct {
	Action  ActionType `param:"Action"`
	Version string     `param:"Version"`
	*SendSmsParams
}

// SendSmsOptions represent SendSmsAction's configurations
type SendSmsOptions interface {
	Options
	Action() string
	Version() string
	RegionID() string
	PhoneNumbers() string
	SignName() string
	TemplateCode() string
	TemplateParam() TemplateParam
	OutID() string

	Response() *SendSmsResponse
}

type sendOptions struct {
	*options
}

func (s *sendOptions) Action() ActionType {
	return s.businessParams.(*sendSmsParams).Action
}

func (s *sendOptions) Version() string {
	return s.businessParams.(*sendSmsParams).Version
}

func (s *sendOptions) RegionID() string {
	return s.businessParams.(*sendSmsParams).RegionID
}

func (s *sendOptions) PhoneNumbers() string {
	return s.businessParams.(*sendSmsParams).PhoneNumbers
}

func (s *sendOptions) SignName() string {
	return s.businessParams.(*sendSmsParams).SignName
}

func (s *sendOptions) TemplateCode() string {
	return s.businessParams.(*sendSmsParams).TemplateCode
}

func (s *sendOptions) TemplateParam() TemplateParam {
	return s.businessParams.(*sendSmsParams).TemplateParam
}

func (s *sendOptions) OutID() string {
	return s.businessParams.(*sendSmsParams).OutID
}

func (s *sendOptions) Response() *SendSmsResponse {
	return s.res.(*SendSmsResponse)
}

// SendSmsAction is action "SendSms"
type SendSmsAction interface {
	action
	Do(extOpts ...Option) (SendSmsOptions, error)
}

type sendAction struct {
	baseAction
}

// Do the send action
func (a *sendAction) Do(extOpts ...Option) (SendSmsOptions, error) {
	opts, err := a.baseAction.doAction(extOpts...)
	if err != nil {
		return nil, err
	}
	return &sendOptions{opts}, nil
}

// NewSendAction init an action "SendSms"
// can be used concurrently
func NewSendAction(c Client, params SendSmsParams) SendSmsAction {
	return &sendAction{
		baseAction{
			&c,
			&sendSmsParams{
				Action:        SendSms,
				Version:       DefaultVersion,
				SendSmsParams: &params,
			},
			reflect.TypeOf(SendSmsResponse{}),
			defaultReqHandler{},
		},
	}
}

// SendSmsResponse is Response of action "SendSms"
type SendSmsResponse struct {
	Response
	BizID string `json:"BizId" xml:"BizId"`
}
