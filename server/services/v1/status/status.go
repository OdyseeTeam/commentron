package status

import (
	"net/http"

	"github.com/lbryio/commentron/meta"
)

type ServerService struct{}

type ServerArgs struct {
}

type Response struct {
	Version string
	Message string
	Running bool
	Commit  string
}

func (t *ServerService) Status(r *http.Request, args *ServerArgs, reply *Response) error {
	reply.Running = true
	reply.Message = meta.GetCommitMessage()
	reply.Commit = meta.GetVersionLong()
	return nil
}
