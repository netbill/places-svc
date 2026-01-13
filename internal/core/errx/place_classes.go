package errx

import "github.com/netbill/ape"

var (
	ErrorPlaceClassNotFound = ape.DeclareError("CLASS_NOT_FOUND")

	ErrorPlaceClassParentCycle = ape.DeclareError("CLASS_PARENT_CYCLE")

	ErrorClassHaveChildren        = ape.DeclareError("CLASS_HAVE_CHILDREN")
	ErrorPlacesExitsWithThisClass = ape.DeclareError("CLASS_EXIT_WITH_THIS_CLASS")
)
