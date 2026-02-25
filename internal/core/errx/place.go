package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotExists = ape.DeclareError("PLACE_NOT_FOUND")

	ErrorCannotSetStatusSuspend = ape.DeclareError("PLACE_CANNOT_SET_SUSPEND")
	ErrorPlaceStatusIsInvalid   = ape.DeclareError("PLACE_STATUS_IS_INVALID")

	ErrorPlaceOutOfTerritory = ape.DeclareError("PLACE_OUT_OF_TERRITORY")

	ErrorPlaceIconTooLarge                = ape.DeclareError("PLACE_ICON_TOO_LARGE")
	ErrorPlaceIconContentTypeNotAllowed   = ape.DeclareError("PLACE_ICON_CONTENT_TYPE_NOT_ALLOWED")
	ErrorPlaceIconContentFormatNotAllowed = ape.DeclareError("PLACE_ICON_CONTENT_FORMAT_NOT_ALLOWED")
	ErrorPlaceIconResolutionNotAllowed    = ape.DeclareError("PLACE_ICON_RESOLUTION_NOT_ALLOWED")

	ErrorPlaceBannerTooLarge                = ape.DeclareError("PLACE_BANNER_TOO_LARGE")
	ErrorPlaceBannerContentTypeNotAllowed   = ape.DeclareError("PLACE_BANNER_CONTENT_TYPE_NOT_ALLOWED")
	ErrorPlaceBannerContentFormatNotAllowed = ape.DeclareError("PLACE_BANNER_CONTENT_FORMAT_NOT_ALLOWED")
	ErrorPlaceBannerResolutionNotAllowed    = ape.DeclareError("PLACE_BANNER_RESOLUTION_NOT_ALLOWED")
)
