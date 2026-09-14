//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var authErr *kvapi.AuthError
	var apiErr *kvapi.APIError
	var conflictErr *kvapi.ConflictError

	runGet := true
	var err error
	var getRequest *kvapi.GetRequest
	var getResponse *kvapi.GetResponse

	var responseValue *string
	var responseUUID uuid.UUID
	for {
		if runGet {
			getRequest = &kvapi.GetRequest{Key: key}
			getResponse, err = c.Get(getRequest)
			runGet = false

			switch {
			case err == nil:
				responseValue = &getResponse.Value
				responseUUID = getResponse.Version
			case errors.As(err, &authErr):
				return err
			case errors.As(err, &apiErr):
				unwrapped := apiErr.Unwrap()
				if unwrapped == kvapi.ErrKeyNotFound {
					responseValue = nil
					responseUUID = uuid.Nil
				} else {
					runGet = true
					continue
				}
			default:
				runGet = true
				continue
			}
		}

		newValue, err := updateFn(responseValue)
		if err != nil {
			return err
		}

		setRequest := kvapi.SetRequest{Key: key, Value: newValue, OldVersion: responseUUID}
		_, err = c.Set(&setRequest)

		if err == nil {
			return nil
		}

		switch {
		case errors.As(err, &conflictErr):
			if conflictErr.ExpectedVersion == uuid.Nil {
				return nil
			}
			runGet = true
		case errors.As(err, &apiErr):
			unwrapped := apiErr.Unwrap()
			switch {
			case unwrapped == kvapi.ErrKeyNotFound:
				responseValue = nil
				responseUUID = uuid.Nil
			case errors.As(unwrapped, &authErr):
				return err
			}
		default:
			return err
		}
	}
}
