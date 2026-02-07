package errx

import "github.com/netbill/ape"

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var (
	ErrorNotEnoughRights = ape.DeclareError("NOT_ENOUGH_RIGHTS")
	ErrorInvalidInput    = ape.DeclareError("INVALID_INPUT")
)
