package moderation

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"
)

// Service is the service struct defined for the comment package for rpc service "reaction.*"
type Service struct{}

// Block returns a list of reactions for the comments requested
func (s Service) Block(r *http.Request, args *commentapi.BlockArgs, reply *commentapi.BlockResponse) error {
	return block(r, args, reply)
}
