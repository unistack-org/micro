package auth

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/auth"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/util/id"
)

// Verify the auth credentials and refresh the auth token periodically
func Verify(a auth.Auth) error {
	// extract the account creds from options, these can be set by flags
	accID := a.Options().ID
	accSecret := a.Options().Secret

	// if no credentials were provided, self generate an account
	if len(accID) == 0 && len(accSecret) == 0 {
		opts := []auth.GenerateOption{
			auth.WithType("service"),
			auth.WithScopes("service"),
		}

		id, err := id.New()
		if err != nil {
			return err
		}
		acc, err := a.Generate(id, opts...)
		if err != nil {
			return err
		}
		if logger.V(logger.DebugLevel) {
			logger.Debug(context.TODO(), "Auth [%v] Generated an auth account: %s", a.String())
		}

		accID = acc.ID
		accSecret = acc.Secret
	}

	// generate the first token
	token, err := a.Token(
		auth.WithCredentials(accID, accSecret),
		auth.WithExpiry(time.Minute*10),
	)
	if err != nil {
		return err
	}

	// set the credentials and token in auth options
	_ = a.Init(
		auth.ClientToken(token),
		auth.Credentials(accID, accSecret),
	)

	// periodically check to see if the token needs refreshing
	go func() {
		timer := time.NewTicker(time.Second * 15)

		for {
			<-timer.C

			// don't refresh the token if it's not close to expiring
			tok := a.Options().Token
			if tok.Expiry.Unix() > time.Now().Add(time.Minute).Unix() {
				continue
			}

			// generate the first token
			tok, err := a.Token(
				auth.WithToken(tok.RefreshToken),
				auth.WithExpiry(time.Minute*10),
			)
			if err != nil {
				if logger.V(logger.WarnLevel) {
					logger.Warn(context.TODO(), "[Auth] Error refreshing token: %v", err)
				}
				continue
			}

			// set the token
			_ = a.Init(auth.ClientToken(tok))
		}
	}()

	return nil
}
