package status

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/meta"
)

// Service is the service for the server package "server.*"
type Service struct{}

// Args arguments for the server.Status rpc call
type Args struct {
}

// Response response for the server.Status rpc call
type Response struct {
	Version string
	Message string
	Running bool
	Commit  string
}

// Status shows the status of commentron
func (t *Service) Status(r *http.Request, args *Args, reply *Response) error {
	reply.Running = true
	reply.Message = meta.GetCommitMessage()
	reply.Commit = meta.GetVersionLong()
	return nil
}
