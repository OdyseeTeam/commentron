package moderation

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"
)

func addDelegate(r *http.Request, args *commentapi.AddDelegateArgs, reply *commentapi.ListDelegateResponse) error {
	return nil
}

func removeDelegate(r *http.Request, args *commentapi.RemoveDelegateArgs, reply *commentapi.ListDelegateResponse) error {
	return nil
}

func listDelegates(r *http.Request, args *commentapi.ListDelegatesArgs, reply *commentapi.ListDelegateResponse) error {
	return nil
}
