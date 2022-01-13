package docs

import (
	"flight_app/app/api/flightaware_api"
	"flight_app/app/store"
)

// swagger:route POST /flightaware_api/calculate Query idOfCalculateEndpoint
// Calculate returns calculated fee.
// responses:
//   200: calculateResponse

// Fee calculated on the basis of flight data.
// swagger:response calculateResponse
type calculateResponseWrapper struct {
	// in:body
	Body flightaware_api.CalculateFeeResponse
}

// swagger:parameters idOfCalculateEndpoint
type calculateParamsWrapper struct {
	// Flight information.
	// in:body
	Body flightaware_api.CalculateFeeRequest
}

// swagger:route POST /flightaware_api/contract/create Query idOfContractCreateEndpoint
// Contract create returns fee and contractID.
// responses:
//   200: createContractResponse

// Fee calculated on the basis of flight data and contract id for payment.
// swagger:response createContractResponse
type createContractResponseWrapper struct {
	// Contract information.
	// in:body
	Body store.CreateContractResponse
}

// swagger:parameters idOfContractCreateEndpoint
type createContractParamsWrapper struct {
	// in:body
	Body store.CreateContractRequest
}

// swagger:route POST /flightaware_api/contracts Query idOfGetContractEndpoint
// GetContracts returns contracts of specified user
// responses:
//   200: []ContractsInfo

// Contracts returns with current status and paid reward
// swagger:response createContractResponse
type getContractsResponseWrapper struct {
	// Contract information.
	// in:body
	Body []store.ContractsInfo
}

// swagger:parameters idOfGetContractEndpoint
type getContractParamsWrapper struct {
	// in:body
	Body store.GetContractsReq
}
