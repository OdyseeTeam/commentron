package moderation

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"
)

// Service is the service struct defined for the comment package for rpc service "moderation.*"
type Service struct{}

// Block returns a list of reactions for the comments requested
func (s Service) Block(r *http.Request, args *commentapi.BlockArgs, reply *commentapi.BlockResponse) error {
	return block(r, args, reply)
}

// AmI return whether or not the users is a moderator and the type. Also the channels they moderate
func (s Service) AmI(r *http.Request, args *commentapi.AmIArgs, reply *commentapi.AmIResponse) error {
	return amI(r, args, reply)
}

// UnBlock return whether or not the users is a moderator and the type. Also the channels they moderate
func (s Service) UnBlock(r *http.Request, args *commentapi.UnBlockArgs, reply *commentapi.UnBlockResponse) error {
	return unBlock(r, args, reply)
}

// BlockedList return the list of blocked channels for a moderator
func (s Service) BlockedList(r *http.Request, args *commentapi.BlockedListArgs, reply *commentapi.BlockedListResponse) error {
	return blockedList(r, args, reply)
}