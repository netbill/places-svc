package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotExists = ape.DeclareError("PLACE_NOT_FOUND")

	ErrorCannotSetStatusSuspend = ape.DeclareError("PLACE_CANNOT_SET_SUSPEND")
	ErrorPlaceStatusIsInvalid   = ape.DeclareError("PLACE_STATUS_IS_INVALID")

	ErrorPlaceOutOfTerritory = ape.DeclareError("PLACE_OUT_OF_TERRITORY")

	ErrorPlaceIconKeyIsInvalid        = ape.DeclareError("PLACE_ICON_KEY_IS_INVALID")
	ErrorPlaceIconContentIsExceedsMax = ape.DeclareError("PLACE_ICON_CONTENT_EXCEEDS_MAX")
	ErrorPlaceIconResolutionIsInvalid = ape.DeclareError("PLACE_ICON_RESOLUTION_IS_INVALID")
	ErrorPlaceIconFormatIsNotAllowed  = ape.DeclareError("PLACE_ICON_FORMAT_IS_NOT_ALLOWED")

	ErrorPlaceBannerKeyIsInvalid        = ape.DeclareError("PLACE_BANNER_KEY_IS_INVALID")
	ErrorPlaceBannerContentIsExceedsMax = ape.DeclareError("PLACE_BANNER_CONTENT_EXCEEDS_MAX")
	ErrorPlaceBannerResolutionIsInvalid = ape.DeclareError("PLACE_BANNER_RESOLUTION_IS_INVALID")
	ErrorPlaceBannerFormatIsNotAllowed  = ape.DeclareError("PLACE_BANNER_FORMAT_IS_NOT_ALLOWED")
)
