package common

import (
	"context"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc/metadata"
)

func GetKeyFromContext(ctx context.Context, key string) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errs.New("no metadata found in context")
	}

	// 获取传递的值
	values := md.Get(key)
	if len(values) == 0 {
		return "", errs.New("no value found for key:", key)
	}
	return values[0], nil
}
