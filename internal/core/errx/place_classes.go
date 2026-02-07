package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceClassNotExists = ape.DeclareError("CLASS_NOT_FOUND")

	ErrorPlaceClassParentCycle = ape.DeclareError("CLASS_PARENT_CYCLE")
	ErrorPlaceClassCodeExists  = ape.DeclareError("CLASS_CODE_EXISTS")

	ErrorPlaceClassHaveChildren   = ape.DeclareError("CLASS_HAVE_CHILDREN")
	ErrorPlacesExitsWithThisClass = ape.DeclareError("CLASS_EXIT_WITH_THIS_CLASS")

	ErrorPlaceClassIconTooLarge                = ape.DeclareError("CLASS_ICON_TOO_LARGE")
	ErrorPlaceClassIconContentTypeNotAllowed   = ape.DeclareError("CLASS_ICON_CONTENT_TYPE_NOT_ALLOWED")
	ErrorPlaceClassIconContentFormatNotAllowed = ape.DeclareError("CLASS_ICON_CONTENT_FORMAT_NOT_ALLOWED")
	ErrorPlaceClassIconResolutionNotAllowed    = ape.DeclareError("CLASS_ICON_RESOLUTION_NOT_ALLOWED")
)
