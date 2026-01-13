package models

import (
	"fmt"

	"github.com/netbill/places-svc/internal/core/errx"
)

//TODO maybe this must be lib

const (
	PlacePossibilityMenu           = "menu"
	PlacePossibilityReserveTable   = "reserve.table"
	PlacePossibilityProductCatalog = "product.catalog"
)

var placePossibilityCodes = []string{
	PlacePossibilityMenu,
	PlacePossibilityReserveTable,
	PlacePossibilityProductCatalog,
}

func GetAllPossibilityCodes() []string {
	return placePossibilityCodes
}

func ValidateCodePossibilityCode(s string) error {
	for _, code := range placePossibilityCodes {
		if code == s {
			return nil
		}
	}

	return errx.ErrorPossibilityCodeIsInvalid.Raise(
		fmt.Errorf("possibility code %s is not valid", s),
	)
}
