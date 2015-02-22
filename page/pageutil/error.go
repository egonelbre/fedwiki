package pageutil

import (
	"os"

	"github.com/egonelbre/fedwiki/page"
)

func ConvertOSError(err error) error {
	if os.IsNotExist(err) {
		return page.ErrNotExist
	}
	if os.IsPermission(err) {
		return page.ErrPermission
	}
	return err
}
