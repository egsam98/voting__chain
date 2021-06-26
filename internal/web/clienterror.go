package web

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (c *ClientError) Error() string {
	return fmt.Sprintf("%s (%d)", c.Err, c.Code)
}

func WrapWithError(clientErr *ClientError, err error) error {
	return errors.Wrap(clientErr, err.Error())
}
