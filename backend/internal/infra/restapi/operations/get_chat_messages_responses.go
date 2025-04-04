// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/babadro/tutor/internal/models/swagger"
)

// GetChatMessagesOKCode is the HTTP code returned for type GetChatMessagesOK
const GetChatMessagesOKCode int = 200

/*
GetChatMessagesOK A list of chat messages

swagger:response getChatMessagesOK
*/
type GetChatMessagesOK struct {

	/*
	  In: Body
	*/
	Payload *GetChatMessagesOKBody `json:"body,omitempty"`
}

// NewGetChatMessagesOK creates GetChatMessagesOK with default headers values
func NewGetChatMessagesOK() *GetChatMessagesOK {

	return &GetChatMessagesOK{}
}

// WithPayload adds the payload to the get chat messages o k response
func (o *GetChatMessagesOK) WithPayload(payload *GetChatMessagesOKBody) *GetChatMessagesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get chat messages o k response
func (o *GetChatMessagesOK) SetPayload(payload *GetChatMessagesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChatMessagesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetChatMessagesBadRequestCode is the HTTP code returned for type GetChatMessagesBadRequest
const GetChatMessagesBadRequestCode int = 400

/*
GetChatMessagesBadRequest Bad request

swagger:response getChatMessagesBadRequest
*/
type GetChatMessagesBadRequest struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewGetChatMessagesBadRequest creates GetChatMessagesBadRequest with default headers values
func NewGetChatMessagesBadRequest() *GetChatMessagesBadRequest {

	return &GetChatMessagesBadRequest{}
}

// WithPayload adds the payload to the get chat messages bad request response
func (o *GetChatMessagesBadRequest) WithPayload(payload *swagger.Error) *GetChatMessagesBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get chat messages bad request response
func (o *GetChatMessagesBadRequest) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChatMessagesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetChatMessagesUnauthorizedCode is the HTTP code returned for type GetChatMessagesUnauthorized
const GetChatMessagesUnauthorizedCode int = 401

/*
GetChatMessagesUnauthorized unauthorized

swagger:response getChatMessagesUnauthorized
*/
type GetChatMessagesUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewGetChatMessagesUnauthorized creates GetChatMessagesUnauthorized with default headers values
func NewGetChatMessagesUnauthorized() *GetChatMessagesUnauthorized {

	return &GetChatMessagesUnauthorized{}
}

// WithPayload adds the payload to the get chat messages unauthorized response
func (o *GetChatMessagesUnauthorized) WithPayload(payload *swagger.Error) *GetChatMessagesUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get chat messages unauthorized response
func (o *GetChatMessagesUnauthorized) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChatMessagesUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetChatMessagesInternalServerErrorCode is the HTTP code returned for type GetChatMessagesInternalServerError
const GetChatMessagesInternalServerErrorCode int = 500

/*
GetChatMessagesInternalServerError Internal server error

swagger:response getChatMessagesInternalServerError
*/
type GetChatMessagesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewGetChatMessagesInternalServerError creates GetChatMessagesInternalServerError with default headers values
func NewGetChatMessagesInternalServerError() *GetChatMessagesInternalServerError {

	return &GetChatMessagesInternalServerError{}
}

// WithPayload adds the payload to the get chat messages internal server error response
func (o *GetChatMessagesInternalServerError) WithPayload(payload *swagger.Error) *GetChatMessagesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get chat messages internal server error response
func (o *GetChatMessagesInternalServerError) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChatMessagesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
GetChatMessagesDefault error

swagger:response getChatMessagesDefault
*/
type GetChatMessagesDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewGetChatMessagesDefault creates GetChatMessagesDefault with default headers values
func NewGetChatMessagesDefault(code int) *GetChatMessagesDefault {
	if code <= 0 {
		code = 500
	}

	return &GetChatMessagesDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get chat messages default response
func (o *GetChatMessagesDefault) WithStatusCode(code int) *GetChatMessagesDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get chat messages default response
func (o *GetChatMessagesDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get chat messages default response
func (o *GetChatMessagesDefault) WithPayload(payload *swagger.Error) *GetChatMessagesDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get chat messages default response
func (o *GetChatMessagesDefault) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetChatMessagesDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
