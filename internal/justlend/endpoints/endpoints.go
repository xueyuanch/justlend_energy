package endpoints

// Response collects the returns values of a common response which
// contains a result to hold all response fields.
//
// For some submit requests we just return whether there is an
// error, the data will be nil in this case.
type Response struct {
	// Result hold all response fields.
	Result interface{} `json:"data,omitempty"`
	// should be intercepted by Failed/errorEncoder
	Err error `json:"-"`
}

// EmptyBody indicates an empty body of a response.
type EmptyBody struct{}

// Failed implements endpoint.Failer.
func (r Response) Failed() error { return r.Err }

func (r Response) Data() interface{} { return r.Result }

// NewResponse is a factory method to construct a new response.
func NewResponse(r interface{}, e error) Response {
	return Response{Result: r, Err: e}
}

// NewErrResponse is a factory method to construct a new response
// that returns whether there is an error in a submit operation, it
// signifies the operation is successful if returns a nil error and
// the data field is left nil.
func NewErrResponse(e error) Response { return Response{Err: e} }

// ListResponse collects the returns values from a find request which
// response contains the result lines of the requested page and a
// count of total matching the find filter.
type ListResponse struct {
	// Lines contains the result lines of the requested page.
	Lines interface{} `json:"lines,omitempty"`
	// Total is a count of total matching the find filter.
	Total int `json:"total"`
	// should be intercepted by Failed/errorEncoder
	Err error `json:"-"`
}

// Failed implements endpoint.Failer.
func (r ListResponse) Failed() error { return r.Err }
func (r ListResponse) Data() interface{} {
	return map[string]interface{}{"lines": r.Lines, "total": r.Total}
}

// NewListResponse is a factory method to construct a new find response.
func NewListResponse(l interface{}, t int, e error) ListResponse { return ListResponse{l, t, e} }
