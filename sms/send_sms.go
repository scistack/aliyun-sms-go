package sms

import (
	"encoding/json"
	"reflect"
)

type TemplateParam map[string]string

func (tp TemplateParam) String() string {
	data, err := json.Marshal(tp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type SendSmsParams struct {
	RegionId      string        `param:"RegionId"`
	PhoneNumbers  string        `param:"PhoneNumbers"`
	SignName      string        `param:"SignName"`
	TemplateCode  string        `param:"TemplateCode"`
	TemplateParam TemplateParam `param:"TemplateParam,omitempty"`
	OutId         string        `param:"OutId,omitempty"`
}

type sendSmsParams struct {
	Action  ActionType `param:"Action"`
	Version string     `param:"Version"`
	*SendSmsParams
}

type SendOptions interface {
	Options
	Action() string
	Version() string
	RegionId() string
	PhoneNumbers() string
	SignName() string
	TemplateCode() string
	TemplateParam() TemplateParam
	OutId() string

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

func (s *sendOptions) RegionId() string {
	return s.businessParams.(*sendSmsParams).RegionId
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

func (s *sendOptions) OutId() string {
	return s.businessParams.(*sendSmsParams).OutId
}

func (s *sendOptions) Response() *SendSmsResponse {
	return s.res.(*SendSmsResponse)
}

type SendAction interface {
	action
	Do(extOpts ...option) (SendOptions, error)
}

type sendAction struct {
	baseAction
}

// Do the send action
func (a *sendAction) Do(extOpts ...option) (SendOptions, error) {
	opts, err := a.baseAction.doAction(extOpts...)
	if err != nil {
		return nil, err
	}
	return &sendOptions{opts}, nil
}

func NewSendAction(c Client, params SendSmsParams) SendAction {
	return &sendAction{
		baseAction{
			&c,
			&sendSmsParams{
				Action:        SendSms,
				Version:       DefaultVersion,
				SendSmsParams: &params,
			},
			reflect.TypeOf(SendSmsResponse{}),
			defaultHandler{},
		},
	}
}

type SendSmsResponse struct {
	Response
	BizID string `json:"BizId" xml:"BizId"`
}
