package validate

import (
	"context"
	"errors"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidationInterceptor struct {
	validator *protovalidate.Validator
}

func NewValidationInterceptor() (*ValidationInterceptor, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &ValidationInterceptor{validator: validator}, nil
}

func (v *ValidationInterceptor) ValidateInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		message, ok := req.(proto.Message)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid request type: not a proto.Message")
		}
		if err := v.validator.Validate(message); err != nil {
			violation := ErrorResponses(err)
			return nil, invalidArgumentError(violation)
		}

		return handler(ctx, req)
	}
}

func invalidArgumentError(violation []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violation}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}
func ErrorResponses(err error) []*errdetails.BadRequest_FieldViolation {
	var details []*errdetails.BadRequest_FieldViolation
	var ve *protovalidate.ValidationError
	if errors.As(err, &ve) {
		for _, violation := range ve.Violations {
			details = append(details, &errdetails.BadRequest_FieldViolation{
				Field:       *violation.FieldPath,
				Description: *violation.Message,
			})
		}

	}
	return details
}
