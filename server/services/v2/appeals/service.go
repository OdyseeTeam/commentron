package appeals

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"
)

// Service is the service struct defined for the comment package for rpc service "appeal.*"
type Service struct{}

// ListBlocks lists the blocks for a particular user that could be potentially appealed.
func (s Service) ListBlocks(r *http.Request, args *commentapi.AppealBlockListArgs, reply *commentapi.AppealBlockListResponse) error {
	return listBlocks(r, args, reply)
}

// ListAppeals lists appeals requested by users against blocks by creators in the shared blocked lists.
func (s Service) ListAppeals(r *http.Request, args *commentapi.AppealListArgs, reply *commentapi.AppealListResponse) error {
	return listAppeals(r, args, reply)
}

// File files an appeal for a particular block
func (s Service) File(r *http.Request, args *commentapi.AppealFileArgs, reply *commentapi.AppealRequest) error {
	return file(r, args, reply)
}

// Close closes an appeal by accepting it or rejecting it.
func (s Service) Close(r *http.Request, args *commentapi.SharedBlockedListGetArgs, reply *commentapi.SharedBlockedListGetResponse) error {
	return nil //get(r, args, reply)
}
