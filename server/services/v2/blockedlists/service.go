package blockedlists

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"
)

// Service is the service struct defined for the comment package for rpc service "blockedlist.*"
type Service struct{}

// Update updates a users shared block list, and returns it adjusted
func (s Service) Update(r *http.Request, args *commentapi.SharedBlockedListUpdateArgs, reply *commentapi.SharedBlockedList) error {
	return update(r, args, reply)
}

// Invite invites a user to contribute to a shared blocked list.
func (s Service) Invite(r *http.Request, args *commentapi.SharedBlockedListInviteArgs, reply *commentapi.SharedBlockedListInviteResponse) error {
	return invite(r, args, reply)
}

// Accept accepts the invite and merges the users blocked entries into the shared blocked list.
func (s Service) Accept(r *http.Request, args *commentapi.SharedBlockedListInviteAcceptArgs, reply *commentapi.SharedBlockedListInviteResponse) error {
	return accept(r, args, reply)
}

// Get accepts the invite and merges the users blocked entries into the shared blocked list.
func (s Service) Get(r *http.Request, args *commentapi.SharedBlockedListGetArgs, reply *commentapi.SharedBlockedListGetResponse) error {
	return get(r, args, reply)
}
