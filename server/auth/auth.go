package auth

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/coreos/go-oidc"
	"github.com/volatiletech/null/v8"
)

const oauthClientID = "commentron"
const oauthProviderURL = "https://sso.odysee.com/auth/realms/Users"

var verifier *oidc.IDTokenVerifier

func init() {
	// This executes a remote call to the OICD provider.
	// Setting the `SKIP_OICD` variable won't kill development without wifi.
	if os.Getenv("SKIP_OICD") == "" {
		provider, err := oidc.NewProvider(context.Background(), oauthProviderURL)
		if err != nil {
			panic(err)
		}
		verifier = provider.Verifier(&oidc.Config{ClientID: oauthClientID})
	}
}

// ErrNotOAuth missing oauth header
var ErrNotOAuth = errors.Base("request does not contain oauth header")

//ModAuthenticate authenticates a moderator
func ModAuthenticate(r *http.Request, modAuthorization *commentapi.ModAuthorization) (*model.Channel, *model.Channel, *UserInfo, error) {
	modChannel, ownerChannel, err := helper.GetModerator(modAuthorization.ModChannelID, modAuthorization.ModChannelName, modAuthorization.CreatorChannelID, modAuthorization.CreatorChannelName)
	if err != nil {
		return nil, nil, nil, err
	}
	var userInfo *UserInfo
	authorization := &commentapi.Authorization{ChannelName: modChannel.Name, ChannelID: modChannel.ClaimID, Signature: modAuthorization.Signature, SigningTS: modAuthorization.SigningTS}
	if modChannel, userInfo, err := oAuth(r, authorization); !errors.Is(err, ErrNotOAuth) {
		if err != nil {
			return nil, nil, nil, err
		}
		return modChannel, ownerChannel, userInfo, nil
	}

	err = lbry.ValidateSignatureAndTS(modChannel.ClaimID, modAuthorization.Signature, modAuthorization.SigningTS, modChannel.Name)
	if err != nil {
		return nil, nil, nil, err
	}
	return modChannel, ownerChannel, userInfo, nil
}

// Authenticate regular authentication
func Authenticate(r *http.Request, authorization *commentapi.Authorization) (*model.Channel, *UserInfo, error) {
	var userInfo *UserInfo
	if channel, userInfo, err := oAuth(r, authorization); !errors.Is(err, ErrNotOAuth) {
		if err != nil {
			return nil, nil, err
		}
		return channel, userInfo, nil
	}

	if authorization == nil {
		return nil, nil, errors.Err("if not authorizing channel for call, must use OAuth")
	}
	err := lbry.ValidateSignatureAndTS(authorization.ChannelID, authorization.Signature, authorization.SigningTS, authorization.ChannelName)
	if err != nil {
		return nil, nil, errors.Prefix("could not authenticate channel signature:", err)
	}
	channel, err := helper.FindOrCreateChannel(authorization.ChannelID, authorization.ChannelName)
	if err != nil {
		return nil, nil, errors.Err(err)
	}

	return channel, userInfo, nil
}

func oAuth(r *http.Request, authorization *commentapi.Authorization) (*model.Channel, *UserInfo, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, nil, ErrNotOAuth
	}
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, nil, errors.Err("token passed must be Bearer token")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	userInfo, err := extractUserInfo(tokenString)
	if err != nil {
		return nil, nil, err
	}

	err = checkAuthorization(userInfo)
	if err != nil {
		return nil, nil, err
	}
	var channel *model.Channel
	if authorization != nil {
		channel, err = model.Channels(model.ChannelWhere.ClaimID.EQ(authorization.ChannelID), model.ChannelWhere.Sub.EQ(null.StringFrom(userInfo.Sub))).One(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, errors.Err(err)
		}
		if channel == nil {
			return nil, nil, errors.Err("could not find verified channel %s for id %s. please verify channel first", authorization.ChannelID, userInfo.Sub)
		}
	}

	return channel, userInfo, nil
}

func checkAuthorization(info *UserInfo) error {
	audienceRespected := false
	for _, aud := range info.Aud {
		if aud == oauthClientID {
			audienceRespected = true
		}
	}
	if !audienceRespected {
		return errors.Err("this token is not meant for Odysee APIs")
	}

	/* If we could get the http request we could valid allowed sources
	allowedSource := false
	for _, url := range info.AllowedOrigins {
		if url == r.Host {
			allowedSource = true
		}
	}
	if !allowedSource {
		return errors.Err("this token cannot be used from %s", r.Host)
	}*/

	return nil
}

// UserInfo contains all claim information included in the access token.
type UserInfo struct {
	Acr               string              `mapstructure:"acr"`
	AllowedOrigins    []string            `json:"allowed-origins"`
	Aud               []string            `mapstructure:"aud"`
	Azp               string              `mapstructure:"azp"`
	Email             string              `mapstructure:"email"`
	EmailVerified     bool                `mapstructure:"email_verified"`
	Exp               int64               `mapstructure:"exp"`
	FamilyName        string              `mapstructure:"family_name"`
	GivenName         string              `mapstructure:"given_name"`
	Iat               int64               `mapstructure:"iat"`
	Iss               string              `mapstructure:"iss"`
	Jti               string              `mapstructure:"jti"`
	Name              string              `mapstructure:"name"`
	PreferredUsername string              `mapstructure:"preferred_username"`
	RealmAccess       map[string][]string `mapstructure:"realm_access"`
	ResourceAccess    struct {
		Commentron struct {
			Roles []string `mapstructure:"roles"`
		} `mapstructure:"commentron"`
	} `mapstructure:"resource_access"`
	Scope        string `mapstructure:"scope"`
	SessionState string `mapstructure:"session_state"`
	Sid          string `mapstructure:"sid"`
	Sub          string `mapstructure:"sub"`
	Typ          string `mapstructure:"typ"`
}

func extractUserInfo(tokenString string) (*UserInfo, error) {
	userInfo := &UserInfo{}

	t, err := verifier.Verify(context.Background(), tokenString)
	if err != nil {
		return nil, errors.Err(err)
	}
	err = t.Claims(userInfo)
	if err != nil {
		return nil, errors.Err(err)
	}

	return userInfo, nil
}
