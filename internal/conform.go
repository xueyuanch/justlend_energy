package internal

import "context"

// Conformer defines a method `Conform` that used to verify that the data
// passed into the service conforms to our expected contract, returns a
// corresponding human-readable failure error if it does and not conform
// to the contract. The contract validations include but not limited to
// field length checks, regex checks etc.
// Usually this verification should be in each service's method of receiving the
// metadata, but for the sake of simplicity, we put it in a middleware of the API
// endpoint to make a unified call, so that we don't need to call this method in
// each service method for verification operations, just make sure that the incoming
// metadata implements this interface.
type Conformer interface{ Conform(context.Context) error }
