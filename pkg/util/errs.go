package util

import (
	"fmt"
	"strings"
)

type ErrorSet struct {
	format string
	Errs   []error
}

func (s *ErrorSet) Error() string {
	msgs := make([]string, 0)
	for _, v := range s.Errs {
		msgs = append(msgs, v.Error())
	}
	return fmt.Sprintf(s.format, strings.Join(msgs, "; "))
}

func MergeErrs(errs []error) error {
	return MergeErrsf("", errs)
}

func MergeErrsf(format string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if format == "" {
		format = "multiple errors: %v"
	}
	errSet := &ErrorSet{
		format: format,
		Errs:   make([]error, 0, len(errs)),
	}
	for _, e := range errs {
		if e != nil {
			errSet.Errs = append(errSet.Errs, e)
		}
	}
	if len(errSet.Errs) == 0 {
		return nil
	}
	return errSet
}
