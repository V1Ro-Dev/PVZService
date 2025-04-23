package server

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	Pvz "pvz/internal/grpc/pvz"
	"pvz/internal/usecase"
)

type PvzManager struct {
	Pvz.UnimplementedPVZServiceServer
	PvzService *usecase.PvzService
}

func NewPvzManager(pvzService *usecase.PvzService) *PvzManager {
	return &PvzManager{
		PvzService: pvzService,
	}
}

func (pm *PvzManager) GetPVZList(ctx context.Context, _ *Pvz.GetPVZListRequest) (*Pvz.GetPVZListResponse, error) {
	pvzList, err := pm.PvzService.GetPvzList(ctx)
	if err != nil {
		return nil, err
	}

	var resp []*Pvz.PVZ
	for _, pvz := range pvzList {
		resp = append(resp, &Pvz.PVZ{
			Id:               pvz.Id.String(),
			City:             pvz.City,
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
		})
	}

	return &Pvz.GetPVZListResponse{
		Pvzs: resp,
	}, nil
}
