// Code generated by go-swagger; DO NOT EDIT.

package cal_dav

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/models"
)

// SubscribeScheduleOKCode is the HTTP code returned for type SubscribeScheduleOK
const SubscribeScheduleOKCode int = 200

/*
SubscribeScheduleOK Subscription successful.

swagger:response subscribeScheduleOK
*/
type SubscribeScheduleOK struct {

	/*
	  In: Body
	*/
	Payload *models.SubscribeResponse `json:"body,omitempty"`
}

// NewSubscribeScheduleOK creates SubscribeScheduleOK with default headers values
func NewSubscribeScheduleOK() *SubscribeScheduleOK {

	return &SubscribeScheduleOK{}
}

// WithPayload adds the payload to the subscribe schedule o k response
func (o *SubscribeScheduleOK) WithPayload(payload *models.SubscribeResponse) *SubscribeScheduleOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the subscribe schedule o k response
func (o *SubscribeScheduleOK) SetPayload(payload *models.SubscribeResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SubscribeScheduleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SubscribeScheduleBadRequestCode is the HTTP code returned for type SubscribeScheduleBadRequest
const SubscribeScheduleBadRequestCode int = 400

/*
SubscribeScheduleBadRequest Bad request.

swagger:response subscribeScheduleBadRequest
*/
type SubscribeScheduleBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSubscribeScheduleBadRequest creates SubscribeScheduleBadRequest with default headers values
func NewSubscribeScheduleBadRequest() *SubscribeScheduleBadRequest {

	return &SubscribeScheduleBadRequest{}
}

// WithPayload adds the payload to the subscribe schedule bad request response
func (o *SubscribeScheduleBadRequest) WithPayload(payload *models.Error) *SubscribeScheduleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the subscribe schedule bad request response
func (o *SubscribeScheduleBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SubscribeScheduleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SubscribeScheduleInternalServerErrorCode is the HTTP code returned for type SubscribeScheduleInternalServerError
const SubscribeScheduleInternalServerErrorCode int = 500

/*
SubscribeScheduleInternalServerError Internal server error.

swagger:response subscribeScheduleInternalServerError
*/
type SubscribeScheduleInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSubscribeScheduleInternalServerError creates SubscribeScheduleInternalServerError with default headers values
func NewSubscribeScheduleInternalServerError() *SubscribeScheduleInternalServerError {

	return &SubscribeScheduleInternalServerError{}
}

// WithPayload adds the payload to the subscribe schedule internal server error response
func (o *SubscribeScheduleInternalServerError) WithPayload(payload *models.Error) *SubscribeScheduleInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the subscribe schedule internal server error response
func (o *SubscribeScheduleInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SubscribeScheduleInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
