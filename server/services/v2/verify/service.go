package verify

import (
	"encoding/hex"
	"net/http"

	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/lbryio/commentron/commentapi"
)

// Service is the service struct defined for the comment package for rpc service "moderation.*"
type Service struct{}

// Block returns a list of reactions for the comments requested
func (s Service) Signature(r *http.Request, args *commentapi.SignatureArgs, reply *commentapi.SignatureResponse) error {
	bytes, err := hex.DecodeString(args.DataHex)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(args.ChannelID, args.Signature, args.SigningTS, string(bytes))
	if err != nil {
		return err
	}
	reply.IsValid = true
	return nil
}
