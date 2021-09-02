package docs

import (
	"flight_app/app/api"
)

// swagger:route POST /api/calculate Query idOfCalculateEndpoint
// Calculate returns calculated fee.
// responses:
//   200: calculateResponse

// Fee calculated on the basis of flight data.
// swagger:response calculateResponse
type calculateResponseWrapper struct {
	// in:body
	Body api.CalculateFeeResponse
}

// swagger:parameters idOfCalculateEndpoint
type calculateParamsWrapper struct {
	// Flight information.
	// in:body
	Body api.CalculateFeeRequest
}

// swagger:route POST /api/contract/create Query idOfContractCreateEndpoint
// Contract create returns
