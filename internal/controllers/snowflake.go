package controllers

import (
	"context"

	flog "github.com/everywan/foundation-go/log"
	"github.com/everywan/identifier"
	"github.com/everywan/identifier/pb"
)

// SnowflakeController is ...
type SnowflakeController struct {
	logger *flog.Logger
	sfSvc  identifier.SnowflakeService
}

// NewSnowflakeController is ...
func NewSnowflakeController(logger *flog.Logger, sfSvc identifier.SnowflakeService) *SnowflakeController {
	return &SnowflakeController{
		logger: logger,
		sfSvc:  sfSvc,
	}
}

var _ pb.SnowflakeServer = &SnowflakeController{}

// Generate is ...
func (s *SnowflakeController) Generate(ctx context.Context, request *pb.Request) (response *pb.Response, err error) {
	response = new(pb.Response)
	id, err := s.sfSvc.Generate(ctx)
	if err != nil {
		return response, err
	}
	response.Uniqid = id
	return response, nil
}
