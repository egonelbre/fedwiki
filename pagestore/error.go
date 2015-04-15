package pagestore

import (
	"os"

	"github.com/egonelbre/fedwiki"
)

// ConvertOSError converts os.Errors to appropriate fedwiki error
func ConvertOSError(err error) error {
	if os.IsNotExist(err) {
		return fedwiki.ErrNotExist
	}
	if os.IsPermission(err) {
		return fedwiki.ErrPermission
	}
	return err
}
