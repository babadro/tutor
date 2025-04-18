// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
)

// SendVoiceMessageHandlerFunc turns a function with the right signature into a send voice message handler
type SendVoiceMessageHandlerFunc func(SendVoiceMessageParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn SendVoiceMessageHandlerFunc) Handle(params SendVoiceMessageParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// SendVoiceMessageHandler interface for that can handle valid send voice message params
type SendVoiceMessageHandler interface {
	Handle(SendVoiceMessageParams, *models.Principal) middleware.Responder
}

// NewSendVoiceMessage creates a new http.Handler for the send voice message operation
func NewSendVoiceMessage(ctx *middleware.Context, handler SendVoiceMessageHandler) *SendVoiceMessage {
	return &SendVoiceMessage{Context: ctx, Handler: handler}
}

/*
SendVoiceMessage swagger:route POST /chat_voice_messages sendVoiceMessage

Sends an audio message to the AI and receives a response.

This endpoint receives a user's audio message and returns the AI's response.
*/
type SendVoiceMessage struct {
	Context *middleware.Context
	Handler SendVoiceMessageHandler
}

func (o *SendVoiceMessage) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSendVoiceMessageParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// SendVoiceMessageOKBody send voice message o k body
//
// swagger:model SendVoiceMessageOKBody
type SendVoiceMessageOKBody struct {

	// If the chatId is not provided in the request, the new chat will be created and this chat will be returned in the response.
	Chat *swagger.Chat `json:"chat,omitempty"`

	// Url to the AI's audio response.
	ReplyAudio string `json:"replyAudio,omitempty"`

	// The timestamp of the AI's response.
	ReplyTime int64 `json:"replyTime,omitempty"`

	// AI's text response to the user's message.
	ReplyTxt string `json:"replyTxt,omitempty"`

	// Url to the user's audio message.
	UsrAudio string `json:"usrAudio,omitempty"`

	// The timestamp of the user's message.
	UsrTime int64 `json:"usrTime,omitempty"`

	// The user's message in text format.
	UsrTxt string `json:"usrTxt,omitempty"`
}

// Validate validates this send voice message o k body
func (o *SendVoiceMessageOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateChat(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SendVoiceMessageOKBody) validateChat(formats strfmt.Registry) error {
	if swag.IsZero(o.Chat) { // not required
		return nil
	}

	if o.Chat != nil {
		if err := o.Chat.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sendVoiceMessageOK" + "." + "chat")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this send voice message o k body based on the context it is used
func (o *SendVoiceMessageOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateChat(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SendVoiceMessageOKBody) contextValidateChat(ctx context.Context, formats strfmt.Registry) error {

	if o.Chat != nil {
		if err := o.Chat.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sendVoiceMessageOK" + "." + "chat")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SendVoiceMessageOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SendVoiceMessageOKBody) UnmarshalBinary(b []byte) error {
	var res SendVoiceMessageOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
