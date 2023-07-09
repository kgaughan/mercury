package utils

import (
	"time"

	"github.com/pquerna/cachecontrol/cacheobject"
)

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
