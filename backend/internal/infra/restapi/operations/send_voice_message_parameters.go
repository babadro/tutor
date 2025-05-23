// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewSendVoiceMessageParams creates a new SendVoiceMessageParams object
// no default values defined in spec.
func NewSendVoiceMessageParams() SendVoiceMessageParams {

	return SendVoiceMessageParams{}
}

// SendVoiceMessageParams contains all the bound params for the send voice message operation
// typically these are obtained from a http.Request
//
// swagger:parameters SendVoiceMessage
type SendVoiceMessageParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The chat ID.
	  In: formData
	*/
	ChatID *string
	/*The audio file to be sent.
	  Required: true
	  In: formData
	*/
	File io.ReadCloser
	/*The timestamp of the message.
	  Required: true
	  In: formData
	*/
	Timestamp int64
	/*
	  In: formData
	*/
	Typ *int32
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSendVoiceMessageParams() beforehand.
func (o *SendVoiceMessageParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}
	fds := runtime.Values(r.Form)

	fdChatID, fdhkChatID, _ := fds.GetOK("chatId")
	if err := o.bindChatID(fdChatID, fdhkChatID, route.Formats); err != nil {
		res = append(res, err)
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "file", err))
	} else if err := o.bindFile(file, fileHeader); err != nil {
		// Required: true
		res = append(res, err)
	} else {
		o.File = &runtime.File{Data: file, Header: fileHeader}
	}

	fdTimestamp, fdhkTimestamp, _ := fds.GetOK("timestamp")
	if err := o.bindTimestamp(fdTimestamp, fdhkTimestamp, route.Formats); err != nil {
		res = append(res, err)
	}

	fdTyp, fdhkTyp, _ := fds.GetOK("typ")
	if err := o.bindTyp(fdTyp, fdhkTyp, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindChatID binds and validates parameter ChatID from formData.
func (o *SendVoiceMessageParams) bindChatID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.ChatID = &raw

	return nil
}

// bindFile binds file parameter File.
//
// The only supported validations on files are MinLength and MaxLength
func (o *SendVoiceMessageParams) bindFile(file multipart.File, header *multipart.FileHeader) error {
	return nil
}

// bindTimestamp binds and validates parameter Timestamp from formData.
func (o *SendVoiceMessageParams) bindTimestamp(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("timestamp", "formData", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("timestamp", "formData", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("timestamp", "formData", "int64", raw)
	}
	o.Timestamp = value

	return nil
}

// bindTyp binds and validates parameter Typ from formData.
func (o *SendVoiceMessageParams) bindTyp(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("typ", "formData", "int32", raw)
	}
	o.Typ = &value

	if err := o.validateTyp(formats); err != nil {
		return err
	}

	return nil
}

// validateTyp carries on validations for parameter Typ
func (o *SendVoiceMessageParams) validateTyp(formats strfmt.Registry) error {

	if err := validate.EnumCase("typ", "formData", *o.Typ, []interface{}{1, 2}, true); err != nil {
		return err
	}

	return nil
}
