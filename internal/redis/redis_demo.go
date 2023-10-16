package redis

import (
	"context"
)

func GetCtx(c context.Context, key string) (string, error) {
	val, err := rdb.Do(c, "get", key).Text()
	if err != nil {
		return "", err
	}
	return val, nil
}
