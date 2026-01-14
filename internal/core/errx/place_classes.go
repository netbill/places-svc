package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceClassNotFound = ape.DeclareError("CLASS_NOT_FOUND")

	ErrorPlaceClassParentCycle = ape.DeclareError("CLASS_PARENT_CYCLE")
	ErrorPlaceClassCodeExists  = ape.DeclareError("CLASS_CODE_EXISTS")

	ErrorPlaceClassHaveChildren   = ape.DeclareError("CLASS_HAVE_CHILDREN")
	ErrorPlacesExitsWithThisClass = ape.DeclareError("CLASS_EXIT_WITH_THIS_CLASS")
)
