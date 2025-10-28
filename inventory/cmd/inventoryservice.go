package main

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mut   sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (service *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	service.mut.RLock()
	defer service.mut.RUnlock()

	part, ok := service.parts[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Unable to find part with UUID %s", req.Uuid)
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (service *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	service.mut.RLock()
	defer service.mut.RUnlock()

	filter := req.Filter
	parts := make([]*inventoryV1.Part, 0)

	for uuid, part := range service.parts {
		if len(filter.Uuids) > 0 && !contains(filter.Uuids, uuid) {
			continue
		}
		if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
			continue
		}
		if len(filter.Categories) > 0 && !contains(filter.Categories, part.Category) {
			continue
		}
		if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if len(filter.Tags) > 0 {
			for _, tag := range filter.Tags {
				if contains(part.Tags, tag) {
					parts = append(parts, part)
					continue
				}
			}
			continue
		}

		// empty filters case
		parts = append(parts, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func contains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
