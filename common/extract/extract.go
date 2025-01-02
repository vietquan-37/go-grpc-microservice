package extract

import (
	"common/interceptor"
	"common/pb"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
)

func UsersMetadata(ctx context.Context) (response *pb.ValidateRsp, err error) {
	response, ok := ctx.Value(interceptor.UserContextKey).(*pb.ValidateRsp)
	if !ok {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response is nil")
	}
	log.Print(response.User)
	return response, nil
}
