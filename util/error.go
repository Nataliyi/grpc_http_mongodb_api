package util

import "github.com/pkg/errors"

func JoinErrors(errs []error) (err error) {
	if len(errs) == 0 {
		return
	}
	err = errs[0]
	for _, newErr := range errs[1:] {
		errors.Wrap(err, newErr.Error())
	}
	return
}
