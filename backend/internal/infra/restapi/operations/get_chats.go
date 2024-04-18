// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/babadro/tutor/internal/models"
	"github.com/babadro/tutor/internal/models/swagger"
)

// GetChatsHandlerFunc turns a function with the right signature into a get chats handler
type GetChatsHandlerFunc func(GetChatsParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetChatsHandlerFunc) Handle(params GetChatsParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetChatsHandler interface for that can handle valid get chats params
type GetChatsHandler interface {
	Handle(GetChatsParams, *models.Principal) middleware.Responder
}

// NewGetChats creates a new http.Handler for the get chats operation
func NewGetChats(ctx *middleware.Context, handler GetChatsHandler) *GetChats {
	return &GetChats{Context: ctx, Handler: handler}
}

/*
GetChats swagger:route GET /chats getChats

Get chats
*/
type GetChats struct {
	Context *middleware.Context
	Handler GetChatsHandler
}

func (o *GetChats) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetChatsParams()

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

// GetChatsOKBody get chats o k body
//
// swagger:model GetChatsOKBody
type GetChatsOKBody struct {

	// chats
	Chats []*swagger.Chat `json:"chats"`
}

// Validate validates this get chats o k body
func (o *GetChatsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateChats(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetChatsOKBody) validateChats(formats strfmt.Registry) error {
	if swag.IsZero(o.Chats) { // not required
		return nil
	}

	for i := 0; i < len(o.Chats); i++ {
		if swag.IsZero(o.Chats[i]) { // not required
			continue
		}

		if o.Chats[i] != nil {
			if err := o.Chats[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getChatsOK" + "." + "chats" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get chats o k body based on the context it is used
func (o *GetChatsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateChats(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetChatsOKBody) contextValidateChats(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Chats); i++ {

		if o.Chats[i] != nil {
			if err := o.Chats[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getChatsOK" + "." + "chats" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetChatsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetChatsOKBody) UnmarshalBinary(b []byte) error {
	var res GetChatsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
