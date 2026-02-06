package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceNotFound = ape.DeclareError("PLACE_NOT_FOUND")

	ErrorPlaceOutOfTerritory = ape.DeclareError("PLACE_OUT_OF_TERRITORY")

	ErrorPlaceIconTooLarge   = ape.DeclareError("PLACE_ICON_TOO_LARGE")
	ErrorPlaceBannerTooLarge = ape.DeclareError("PLACE_BANNER_TOO_LARGE")

	ErrorPlaceIconContentTypeNotAllowed   = ape.DeclareError("PLACE_ICON_CONTENT_TYPE_NOT_ALLOWED")
	ErrorPlaceBannerContentTypeNotAllowed = ape.DeclareError("PLACE_BANNER_CONTENT_TYPE_NOT_ALLOWED")

	ErrorPlaceIconContentFormatNotAllowed   = ape.DeclareError("PLACE_ICON_CONTENT_FORMAT_NOT_ALLOWED")
	ErrorPlaceBannerContentFormatNotAllowed = ape.DeclareError("PLACE_BANNER_CONTENT_FORMAT_NOT_ALLOWED")

	ErrorPlaceIconResolutionNotAllowed   = ape.DeclareError("PLACE_ICON_RESOLUTION_NOT_ALLOWED")
	ErrorPlaceBannerResolutionNotAllowed = ape.DeclareError("PLACE_BANNER_RESOLUTION_NOT_ALLOWED")
)
