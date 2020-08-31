package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lbryio/commentron/server/services/v1/comments"
	rpcHack "github.com/lbryio/commentron/server/services/v1/rpc"
	jsonHack "github.com/lbryio/commentron/server/services/v1/rpc/json"
	"github.com/lbryio/commentron/server/services/v1/status"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	json "github.com/gorilla/rpc/v2/json2"
	"github.com/sirupsen/logrus"
)

// RPCHost specifies the host name the rpc server binds to
var RPCHost string

// RPCPort specifies the port the rpc server listens on
var RPCPort int

// Start starts the rpc server after any configuration
func Start() {
	logrus.SetOutput(os.Stdout)

	router := mux.NewRouter()
	router.Handle("/api", v1RPCServer())
	router.Handle("/api/v1", v1RPCServer())
	router.Handle("/api/v2", v2RPCServer())
	logrus.Infof("Running RPC Server @ http://%s:%d/api", RPCHost, RPCPort)
	address := fmt.Sprintf("%s:%d", RPCHost, RPCPort)
	logrus.Fatal(http.ListenAndServe(address, router))
}

func v1RPCServer() http.Handler {
	// This rpc server is a copy of github.com/gorilla/rpc/v2 because someone thought it was smart to ignore the
	// json 2.0 specification where ids must be sent if you want a response and instead thought it wise to allow
	// null ids to be sent while expecting a response. The specification requires null ids to be translated as a
	// notification which gets no response. So I had to copy an entire library to make that one change, @jessop
	// has at least fixed the sdk for the future. Hopefully at some point in the future we can do it right.
	rpcServer := rpcHack.NewServer()

	rpcServer.RegisterCodec(jsonHack.NewCodec(), "application/json")
	rpcServer.RegisterCodec(jsonHack.NewCodec(), "application/json;charset=UTF-8")

	commentService := new(comments.Service)
	statusService := new(status.Service)

	err := rpcServer.RegisterService(commentService, "comment")
	if err != nil {
		logrus.Panicf("Error registering comment service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(statusService, "server")
	if err != nil {
		logrus.Panicf("Error registering status service: %s", errors.FullTrace(err))
	}
	rpcServer.RegisterBeforeFunc(func(info *rpcHack.RequestInfo) {
		logrus.Debugf("M->%s: from %s, %d", info.Method, getIP(info.Request), info.StatusCode)
	})
	rpcServer.RegisterAfterFunc(func(info *rpcHack.RequestInfo) {
		if info.Error != nil {
			info.StatusCode = http.StatusInternalServerError
			logrus.Error(errors.FullTrace(info.Error))
		}
	})

	return rpcServer
}

func v2RPCServer() http.Handler {
	rpcServer := rpc.NewServer()

	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	commentService := new(comments.Service)
	statusService := new(status.Service)

	err := rpcServer.RegisterService(commentService, "comment")
	if err != nil {
		logrus.Panicf("Error registering v2 comment service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(statusService, "server")
	if err != nil {
		logrus.Panicf("Error registering v2 status service: %s", errors.FullTrace(err))
	}
	rpcServer.RegisterBeforeFunc(func(info *rpc.RequestInfo) {
		logrus.Debugf("M->%s: from %s, %d", info.Method, getIP(info.Request), info.StatusCode)
	})
	rpcServer.RegisterAfterFunc(func(info *rpc.RequestInfo) {
		if info.Error != nil {
			info.StatusCode = http.StatusInternalServerError
			logrus.Error(errors.FullTrace(info.Error))
		}
	})

	return rpcServer
}

// getIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
