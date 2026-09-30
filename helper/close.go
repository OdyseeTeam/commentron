package helper

import (
	"io"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/sirupsen/logrus"
)

// CloseBody used to catch and log errors from closing a response body
func CloseBody(responseBody io.ReadCloser) {
	if err := responseBody.Close(); err != nil {
		closeError := errors.Prefix("closing body response error: ", errors.Err(err))
		logrus.Error(closeError)
	}
}

// CloseWriter used to catch and log errors from closing a writer
func CloseWriter(w io.WriteCloser) {
	if err := w.Close(); err != nil {
		closeError := errors.Prefix("closing writer response error: ", errors.Err(err))
		logrus.Error(closeError)
	}
}
