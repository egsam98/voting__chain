package web

import (
	"github.com/pkg/errors"
)

type ClientError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (c *ClientError) Error() string {
	return c.Err
}

func WrapWithError(clientErr *ClientError, err error) error {
	return errors.Wrap(clientErr, err.Error())
}
