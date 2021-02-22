package validator

import (
	"regexp"

	v "github.com/lbryio/ozzo-validation"
)

var (
	TransactionIDRegex = `^[A-Za-z0-9]{64}$`
	TransactionID      = v.NewStringRule(func(str string) bool {
		return matchesRegex(TransactionIDRegex, str)
	}, "Invalid transaction id")
	ClaimIDRegex = `^[A-Za-z0-9]{40}$`
	ClaimID      = v.NewStringRule(func(str string) bool {
		return matchesRegex(ClaimIDRegex, str)
	}, "Invalid claim id")
)

func matchesRegex(regex, str string) bool {
	matched, err := regexp.Match(regex, []byte(str))
	return matched && err == nil
}
