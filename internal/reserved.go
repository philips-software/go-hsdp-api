package internal

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var ReservedStrings = []string{
	"roles", "permissions", "propositions", "applications", "services", "oauth2clients",
	"devices", "emailtemplates", "passwordpolicies", "accesspolicies", "delegations",
	"cn=", "dn=", "uid="}

func ValidateNoReservedStrings(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	for _, r := range ReservedStrings {
		if strings.Contains(val, r) {
			return false
		}
	}
	return true
}
