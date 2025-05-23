// Code generated by go-swagger; DO NOT EDIT.

package cal_dav

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// SubscribeScheduleHandlerFunc turns a function with the right signature into a subscribe schedule handler
type SubscribeScheduleHandlerFunc func(SubscribeScheduleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SubscribeScheduleHandlerFunc) Handle(params SubscribeScheduleParams) middleware.Responder {
	return fn(params)
}

// SubscribeScheduleHandler interface for that can handle valid subscribe schedule params
type SubscribeScheduleHandler interface {
	Handle(SubscribeScheduleParams) middleware.Responder
}

// NewSubscribeSchedule creates a new http.Handler for the subscribe schedule operation
func NewSubscribeSchedule(ctx *middleware.Context, handler SubscribeScheduleHandler) *SubscribeSchedule {
	return &SubscribeSchedule{Context: ctx, Handler: handler}
}

/*
	SubscribeSchedule swagger:route POST /subscribe CalDav subscribeSchedule

Subscribe and generate iCal for user.

Subscribes user by ISU and password, generates and stores iCal file.
*/
type SubscribeSchedule struct {
	Context *middleware.Context
	Handler SubscribeScheduleHandler
}

func (o *SubscribeSchedule) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewSubscribeScheduleParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
