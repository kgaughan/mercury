package utils

import (
	"time"

	"github.com/pquerna/cachecontrol/cacheobject"
)

// ParseCacheControlExpiration parses the Cache-Control header and sets the
// expiration time based on its directives.
func ParseCacheControlExpiration(cc string, expires *time.Time) error {
	resDir, err := cacheobject.ParseResponseCacheControl(cc)
	if err != nil {
		return err //nolint:wrapcheck
	}
	*expires = time.Now().UTC()
	if !resDir.NoCachePresent && resDir.MaxAge > 0 {
		*expires = expires.Add(time.Second * time.Duration(resDir.MaxAge))
	}
	return nil
}
