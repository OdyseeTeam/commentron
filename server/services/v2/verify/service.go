package verify

import (
	"encoding/hex"
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// Service is the service struct defined for the comment package for rpc service "moderation.*"
type Service struct{}

// Signature returns a list of reactions for the comments requested
func (s Service) Signature(r *http.Request, args *commentapi.SignatureArgs, reply *commentapi.SignatureResponse) error {
	bytes, err := hex.DecodeString(args.DataHex)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignatureAndTS(args.ChannelID, args.Signature, args.SigningTS, string(bytes))
	if err != nil {
		return err
	}
	reply.IsValid = true
	return nil
}

// ClaimSignature validates a channel signed a particular claim id
func (s Service) ClaimSignature(r *http.Request, args *commentapi.SignatureArgs, reply *commentapi.SignatureResponse) error {
	err := lbry.ValidateSignatureAndTSForClaim(args.ChannelID, args.ClaimID, args.Signature, args.SigningTS, args.ClaimID)
	if err != nil {
		return err
	}
	reply.IsValid = true
	return nil
}
