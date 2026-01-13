package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotFound = ape.DeclareError("PLACE_NOT_FOUND")

	ErrorPlaceStatusSuspended        = ape.DeclareError("PLACE_STATUS_SUSPENDED")
	ErrorCannotSetPlaceStatusSuspend = ape.DeclareError("CANNOT_SET_PLACE_STATUS_SUSPEND")

	ErrorPlaceOutOfTerritory = ape.DeclareError("PLACE_OUT_OF_TERRITORY")
)
