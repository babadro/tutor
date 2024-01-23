// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/babadro/tutor/internal/models/swagger"
)

// SendChatMessageOKCode is the HTTP code returned for type SendChatMessageOK
const SendChatMessageOKCode int = 200

/*
SendChatMessageOK Successful response

swagger:response sendChatMessageOK
*/
type SendChatMessageOK struct {

	/*
	  In: Body
	*/
	Payload *SendChatMessageOKBody `json:"body,omitempty"`
}

// NewSendChatMessageOK creates SendChatMessageOK with default headers values
func NewSendChatMessageOK() *SendChatMessageOK {

	return &SendChatMessageOK{}
}

// WithPayload adds the payload to the send chat message o k response
func (o *SendChatMessageOK) WithPayload(payload *SendChatMessageOKBody) *SendChatMessageOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send chat message o k response
func (o *SendChatMessageOK) SetPayload(payload *SendChatMessageOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendChatMessageOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SendChatMessageBadRequestCode is the HTTP code returned for type SendChatMessageBadRequest
const SendChatMessageBadRequestCode int = 400

/*
SendChatMessageBadRequest Bad request

swagger:response sendChatMessageBadRequest
*/
type SendChatMessageBadRequest struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewSendChatMessageBadRequest creates SendChatMessageBadRequest with default headers values
func NewSendChatMessageBadRequest() *SendChatMessageBadRequest {

	return &SendChatMessageBadRequest{}
}

// WithPayload adds the payload to the send chat message bad request response
func (o *SendChatMessageBadRequest) WithPayload(payload *swagger.Error) *SendChatMessageBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send chat message bad request response
func (o *SendChatMessageBadRequest) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendChatMessageBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SendChatMessageUnauthorizedCode is the HTTP code returned for type SendChatMessageUnauthorized
const SendChatMessageUnauthorizedCode int = 401

/*
SendChatMessageUnauthorized unauthorized

swagger:response sendChatMessageUnauthorized
*/
type SendChatMessageUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewSendChatMessageUnauthorized creates SendChatMessageUnauthorized with default headers values
func NewSendChatMessageUnauthorized() *SendChatMessageUnauthorized {

	return &SendChatMessageUnauthorized{}
}

// WithPayload adds the payload to the send chat message unauthorized response
func (o *SendChatMessageUnauthorized) WithPayload(payload *swagger.Error) *SendChatMessageUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send chat message unauthorized response
func (o *SendChatMessageUnauthorized) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendChatMessageUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SendChatMessageInternalServerErrorCode is the HTTP code returned for type SendChatMessageInternalServerError
const SendChatMessageInternalServerErrorCode int = 500

/*
SendChatMessageInternalServerError Internal server error

swagger:response sendChatMessageInternalServerError
*/
type SendChatMessageInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewSendChatMessageInternalServerError creates SendChatMessageInternalServerError with default headers values
func NewSendChatMessageInternalServerError() *SendChatMessageInternalServerError {

	return &SendChatMessageInternalServerError{}
}

// WithPayload adds the payload to the send chat message internal server error response
func (o *SendChatMessageInternalServerError) WithPayload(payload *swagger.Error) *SendChatMessageInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send chat message internal server error response
func (o *SendChatMessageInternalServerError) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendChatMessageInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
SendChatMessageDefault error

swagger:response sendChatMessageDefault
*/
type SendChatMessageDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewSendChatMessageDefault creates SendChatMessageDefault with default headers values
func NewSendChatMessageDefault(code int) *SendChatMessageDefault {
	if code <= 0 {
		code = 500
	}

	return &SendChatMessageDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the send chat message default response
func (o *SendChatMessageDefault) WithStatusCode(code int) *SendChatMessageDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the send chat message default response
func (o *SendChatMessageDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the send chat message default response
func (o *SendChatMessageDefault) WithPayload(payload *swagger.Error) *SendChatMessageDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send chat message default response
func (o *SendChatMessageDefault) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendChatMessageDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
