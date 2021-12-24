package server

import (
	"bytes"
	jsonmarshall "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lbryio/commentron/server/services/v2/appeals"

	"github.com/lbryio/commentron/config"

	"github.com/lbryio/commentron/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lbryio/commentron/server/websocket"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/server/services/v1/comments"
	rpcHack "github.com/lbryio/commentron/server/services/v1/rpc"
	jsonHack "github.com/lbryio/commentron/server/services/v1/rpc/json"
	"github.com/lbryio/commentron/server/services/v1/status"
	"github.com/lbryio/commentron/server/services/v2/blockedlists"
	"github.com/lbryio/commentron/server/services/v2/moderation"
	"github.com/lbryio/commentron/server/services/v2/reactions"
	"github.com/lbryio/commentron/server/services/v2/settings"
	"github.com/lbryio/commentron/server/services/v2/verify"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	json "github.com/gorilla/rpc/v2/json2"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// RPCHost specifies the host name the rpc server binds to
var RPCHost string

// RPCPort specifies the port the rpc server listens on
var RPCPort int

const promPath = "/metrics"

var corsHandler = cors.New(cors.Options{
	AllowedHeaders: []string{"Authorization", "*"},
	MaxAge:         1728000,
}).Handler

// Start starts the rpc server after any configuration
func Start() {
	logrus.SetOutput(os.Stdout)
	chain := alice.New(corsHandler)
	router := mux.NewRouter()
	router.Handle("/", state())
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.Handle("/api", v1RPCServer())
	router.Handle("/api/v1", v1RPCServer())
	router.Handle("/api/v2", chain.Then(v2RPCServer()))
	router.Handle("/api/v2/live-chat/subscribe", websocket.SubscribeLiveChat())
	router.Handle(promPath, promBasicAuthWrapper(promhttp.Handler()))

	mux := http.Handler(router)
	for _, middleware := range []func(h http.Handler) http.Handler{
		promRequestHandler,
	} {
		mux = middleware(mux)
	}

	logrus.Infof("Running RPC Server @ http://%s:%d/api", RPCHost, RPCPort)
	address := fmt.Sprintf("%s:%d", RPCHost, RPCPort)
	logrus.Fatal(http.ListenAndServe(address, mux))
}

func promRequestHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.URL.Path, "/")
		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		codecRequest := json.NewCodec().NewRequest(r)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if method, err := codecRequest.Method(); err == nil {
			version, service, method := getCallDetails(path, method)
			metrics.UserLoadOverall.Inc()
			defer metrics.UserLoadOverall.Dec()
			metrics.UserLoadByAPI.WithLabelValues(version, service, method).Inc()
			defer metrics.UserLoadByAPI.WithLabelValues(version, service, method).Dec()
			apiStart := time.Now()
			h.ServeHTTP(w, r)
			duration := time.Since(apiStart).Seconds()
			metrics.Durations.WithLabelValues(version, service, method).Observe(duration)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func promBasicAuthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authentication required", http.StatusBadRequest)
			return
		}
		if user == "prom" && pass == "prom-commentron-access" {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, "invalid username or password", http.StatusForbidden)
		}
	})
}

func state() http.Handler {
	var startUp = time.Now()
	type status struct {
		Text      string  `json:"text"`
		IsRunning bool    `json:"is_running"`
		Uptime    float64 `json:"up_time"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonBytes, err := jsonmarshall.Marshal(&status{
			Text:      "OK",
			IsRunning: true,
			Uptime:    time.Since(startUp).Seconds(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				logrus.Error(errors.FullTrace(err))
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(jsonBytes)
			if err != nil {
				logrus.Error(errors.FullTrace(err))
			}
		}
	})
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
		consoleText := info.Request.RemoteAddr + " [" + strconv.Itoa(info.StatusCode) + "]: " + info.Method
		if info.Error != nil {
			err, ok := info.Error.(api.StatusError)
			if ok {
				info.StatusCode = err.Status
			} else {
				message := info.Error.Error()
				if config.IsTestMode {
					message = errors.FullTrace(info.Error)
				}
				logrus.Error(color.RedString(consoleText + ": " + message))
			}
		} else if helper.Debugging {
			logrus.Debug(color.GreenString(consoleText))
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
	reactionService := new(reactions.Service)
	moderationService := new(moderation.Service)
	settingService := new(settings.Service)
	verifyService := new(verify.Service)
	blockedlistService := new(blockedlists.Service)
	appealsService := new(appeals.Service)

	err := rpcServer.RegisterService(commentService, "comment")
	if err != nil {
		logrus.Panicf("Error registering v2 comment service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(statusService, "server")
	if err != nil {
		logrus.Panicf("Error registering v2 status service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(reactionService, "reaction")
	if err != nil {
		logrus.Panicf("Error registering v2 reaction service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(moderationService, "moderation")
	if err != nil {
		logrus.Panicf("Error registering v2 moderation service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(settingService, "setting")
	if err != nil {
		logrus.Panicf("Error registering v2 setting service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(verifyService, "verify")
	if err != nil {
		logrus.Panicf("Error registering v2 verify service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(blockedlistService, "blockedlist")
	if err != nil {
		logrus.Panicf("Error registering v2 verify service: %s", errors.FullTrace(err))
	}
	err = rpcServer.RegisterService(appealsService, "appeals")
	if err != nil {
		logrus.Panicf("Error registering v2 verify service: %s", errors.FullTrace(err))
	}
	rpcServer.RegisterBeforeFunc(func(info *rpc.RequestInfo) {
		logrus.Debugf("M->%s: from %s, %d", info.Method, getIP(info.Request), info.StatusCode)
	})
	rpcServer.RegisterAfterFunc(func(info *rpc.RequestInfo) {
		consoleText := info.Request.RemoteAddr + " [" + strconv.Itoa(info.StatusCode) + "]: " + info.Method
		if info.Error != nil {
			statusErr, ok := info.Error.(api.StatusError)
			if ok {
				info.StatusCode = statusErr.Status
			}
			if info.StatusCode >= http.StatusInternalServerError {
				message := err.Error()
				if config.IsTestMode {
					message = errors.FullTrace(err)
				}
				logrus.Error(color.RedString(consoleText + ": " + message))
			} else {
				logrus.Debug(color.RedString(consoleText + ": " + info.Error.Error()))
			}
		} else if helper.Debugging {
			logrus.Debug(color.GreenString(consoleText))
		}
	})

	rpcServer.RegisterValidateRequestFunc(func(r *rpc.RequestInfo, i interface{}) error {
		v, ok := i.(commentapi.Validator)
		if ok {
			err := v.Validate()
			if err.Err != nil {
				return errors.Err(err)
			}
		}
		return nil
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

func getCallDetails(path, method string) (string, string, string) {
	version := strings.TrimPrefix(path, "api/")
	method = strings.ToLower(method)
	parts := strings.Split(method, ".")
	if len(parts) == 0 {
		return version, "", ""
	}
	if len(parts) < 2 {
		return version, parts[0], ""
	}
	return version, parts[0], parts[1]
}
