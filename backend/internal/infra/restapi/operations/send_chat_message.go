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
	"github.com/go-openapi/validate"

	"github.com/babadro/tutor/internal/models"
)

// SendChatMessageHandlerFunc turns a function with the right signature into a send chat message handler
type SendChatMessageHandlerFunc func(SendChatMessageParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn SendChatMessageHandlerFunc) Handle(params SendChatMessageParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// SendChatMessageHandler interface for that can handle valid send chat message params
type SendChatMessageHandler interface {
	Handle(SendChatMessageParams, *models.Principal) middleware.Responder
}

// NewSendChatMessage creates a new http.Handler for the send chat message operation
func NewSendChatMessage(ctx *middleware.Context, handler SendChatMessageHandler) *SendChatMessage {
	return &SendChatMessage{Context: ctx, Handler: handler}
}

/*
SendChatMessage swagger:route POST /chat_messages sendChatMessage

Sends a message to the AI and receives a response.

This endpoint receives a user's message and returns the AI's response.
*/
type SendChatMessage struct {
	Context *middleware.Context
	Handler SendChatMessageHandler
}

func (o *SendChatMessage) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSendChatMessageParams()

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

// SendChatMessageBody send chat message body
//
// swagger:model SendChatMessageBody
type SendChatMessageBody struct {

	// The chat ID.
	// Required: true
	ChatID *string `json:"chatId"`

	// The message text sent by the user.
	// Required: true
	Text *string `json:"text"`

	// The timestamp of the message.
	// Required: true
	Timestamp *int64 `json:"timestamp"`
}

// Validate validates this send chat message body
func (o *SendChatMessageBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateChatID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateText(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SendChatMessageBody) validateChatID(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"chatId", "body", o.ChatID); err != nil {
		return err
	}

	return nil
}

func (o *SendChatMessageBody) validateText(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"text", "body", o.Text); err != nil {
		return err
	}

	return nil
}

func (o *SendChatMessageBody) validateTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"timestamp", "body", o.Timestamp); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this send chat message body based on context it is used
func (o *SendChatMessageBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *SendChatMessageBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SendChatMessageBody) UnmarshalBinary(b []byte) error {
	var res SendChatMessageBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SendChatMessageOKBody send chat message o k body
//
// swagger:model SendChatMessageOKBody
type SendChatMessageOKBody struct {

	// AI's response to the user's message.
	Reply string `json:"reply,omitempty"`

	// timestamp
	Timestamp int64 `json:"timestamp,omitempty"`
}

// Validate validates this send chat message o k body
func (o *SendChatMessageOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this send chat message o k body based on context it is used
func (o *SendChatMessageOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *SendChatMessageOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SendChatMessageOKBody) UnmarshalBinary(b []byte) error {
	var res SendChatMessageOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
