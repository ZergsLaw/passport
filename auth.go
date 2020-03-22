package passport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Auth the object that performs the authorization.
type Auth struct {
	client      OauthClient
	ctx         context.Context
	maxBodySize int64
}

// Errors.
var (
	ErrUnknownOAuthClient = errors.New("unknown oauth client")
)

// Do calls the oauth server and returns the error response or unpacks the resulting json into dest.
func (a *Auth) Do(code string, dest interface{}) error {
	if a.client == nil {
		return ErrUnknownOAuthClient
	}

	status, body, err := a.client.Account(a.ctx, code)
	if err != nil {
		return fmt.Errorf("client account: %w", err)
	}
	defer body.Close()

	lr := io.LimitReader(body, a.maxBodySize)

	if status != http.StatusOK {
		myErr := &OAuthErr{}

		if err := json.NewDecoder(lr).Decode(&myErr.Value); err != nil {
			return fmt.Errorf("json decode error: %w", err)
		}

		return myErr
	}

	err = json.NewDecoder(lr).Decode(dest)
	if err != nil {
		return fmt.Errorf("json decode result: %w", err)
	}

	return nil
}

// OAuthErr error from OAuth server.
type OAuthErr struct {
	Value interface{}
}

func (err OAuthErr) Error() string {
	return fmt.Sprintf("%v", err.Value)
}
