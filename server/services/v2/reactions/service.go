package reactions

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
)

// Service is the service struct defined for the comment package for rpc service "reaction.*"
type Service struct{}

// List returns a list of reactions for the comments requested
func (s Service) List(r *http.Request, args *commentapi.ReactionListArgs, reply *commentapi.ReactionListResponse) error {
	return list(r, args, reply)
}

// React creates reactions for comments.
func (s Service) React(r *http.Request, args *commentapi.ReactArgs, reply *commentapi.ReactResponse) error {
	return react(r, args, reply)
}
