package commentclassification

import (
	"net/http"
	"strings"
	"sync"

	"github.com/OdyseeTeam/commentron/env"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/sirupsen/logrus"
)

var (
	commentClassificationInProgress bool
	commentClassificationMutex      sync.Mutex
	defaultBatchSize                = 10
	inferenceServiceURI             = "http://localhost:2345/api/v1/wtf"
	authUser                        = ""
	authPass                        = ""
)

// Init gets and sets the comment classifier api url
func Init(conf *env.Config) {
	inferenceServiceURI = conf.CommentClassifierAPIURL
	if conf.IsTestMode {
		logrus.Info("Comment classification server: ", inferenceServiceURI)
	}

	// Do not bring down commentron if this fails. Only make it obvious in the api call.
	parts := strings.Split(conf.CommentClassificationEndpointAuth, ":")
	if len(parts) != 2 {
		logrus.Error("COMMENT_CLASSIFICATION_ENDPOINT_AUTH is invalid")
	} else {
		authUser, authPass = parts[0], parts[1]
		if len(parts) != 2 {
			logrus.Error("COMMENT_CLASSIFICATION_ENDPOINT_AUTH not set")
		}
	}
}

// IsAuthenticated checks if the request is authenticated for comment classification calls
//
// This is a temporary solution.
func IsAuthenticated(r *http.Request) error {
	// WARNING: Do not take down commentron if the password is not set. Just fail
	// when doing the action api calls in a loud way.
	if authUser == "" || authPass == "" {
		return errors.Err("COMMENT_CLASSIFICATION_ENDPOINT_AUTH not set")
	}

	gotUser, gotPass, ok := r.BasicAuth()
	if !ok {
		return errors.Err("authentication required")
	}

	if authUser != gotUser && authPass != gotPass {
		return errors.Err("invalid credentials")
	}

	return nil
}
