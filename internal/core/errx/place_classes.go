package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceClassNotExists    = ape.DeclareError("CLASS_NOT_FOUND")
	ErrorPlaceClassIsDeprecated = ape.DeclareError("CLASS_IS_DEPRECATED")
	ErrorPlaceClassParentCycle  = ape.DeclareError("CLASS_PARENT_CYCLE")

	ErrorNoContentUploaded = ape.DeclareError("NO_CONTENT_UPLOADED")

	ErrorPlaceClassIconKeyIsInvalid        = ape.DeclareError("PLACE_CLASS_AVATAR_KEY_IS_INVALID")
	ErrorPlaceClassIconContentIsExceedsMax = ape.DeclareError("PLACE_CLASS_AVATAR_CONTENT_EXCEEDS_MAX")
	ErrorPlaceClassIconResolutionIsInvalid = ape.DeclareError("PLACE_CLASS_AVATAR_RESOLUTION_IS_INVALID")
	ErrorPlaceClassIconFormatIsNotAllowed  = ape.DeclareError("PLACE_CLASS_AVATAR_FORMAT_IS_NOT_ALLOWED")
)
