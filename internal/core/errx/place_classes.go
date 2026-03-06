package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceClassNotExists    = ape.DeclareError("CLASS_NOT_FOUND")
	ErrorPlaceClassIsDeprecated = ape.DeclareError("CLASS_IS_DEPRECATED")
	ErrorPlaceClassParentCycle  = ape.DeclareError("CLASS_PARENT_CYCLE")

	ErrorPlaceClassIconIsInvalid = ape.DeclareError("CLASS_ICON_KEY_IS_INVALID")
)
