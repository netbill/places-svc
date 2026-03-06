package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotExists = ape.DeclareError("PLACE_NOT_FOUND")
	ErrorPlaceDeleted   = ape.DeclareError("PLACE_DELETED")

	ErrorPlaceStatusIsInvalid = ape.DeclareError("PLACE_STATUS_IS_INVALID")

	ErrorPlaceOutOfTerritory = ape.DeclareError("PLACE_OUT_OF_TERRITORY")

	ErrorPlaceIconIsInvalid   = ape.DeclareError("PLACE_ICON_KEY_IS_INVALID")
	ErrorPlaceBannerIsInvalid = ape.DeclareError("PLACE_BANNER_KEY_IS_INVALID")
)
