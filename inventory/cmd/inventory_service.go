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

	filters := covnertToFilterSets(req.Filter)

	filteredParts := make([]*inventoryV1.Part, 0)

	for uuid, part := range service.parts {
		if len(filters.Uuids) > 0 && !mapContains(filters.Uuids, uuid) {
			continue
		}
		if len(filters.Names) > 0 && !mapContains(filters.Names, part.Name) {
			continue
		}
		if len(filters.Categories) > 0 && !mapContains(filters.Categories, part.Category) {
			continue
		}
		if len(filters.ManufacturerCountries) > 0 && !mapContains(filters.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if len(filters.Tags) > 0 {
			for tag, _ := range filters.Tags {
				if sliceContains(part.Tags, tag) {
					filteredParts = append(filteredParts, part)
					continue
				}
			}
			continue
		}

		// empty filters case
		filteredParts = append(filteredParts, part)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

type filterSets struct {
	Uuids                 map[string]any
	Names                 map[string]any
	Categories            map[inventoryV1.Category]any
	ManufacturerCountries map[string]any
	Tags                  map[string]any
}

func covnertToFilterSets(filter *inventoryV1.PartsFilter) *filterSets {
	filterSets := &filterSets{}

	filterSets.Uuids = make(map[string]any)
	filterSets.Names = make(map[string]any)
	filterSets.Categories = make(map[inventoryV1.Category]any)
	filterSets.ManufacturerCountries = make(map[string]any)
	filterSets.Tags = make(map[string]any)

	for _, uuid := range filter.Uuids {
		filterSets.Uuids[uuid] = struct{}{}
	}
	for _, name := range filter.Names {
		filterSets.Names[name] = struct{}{}
	}
	for _, category := range filter.Categories {
		filterSets.Categories[category] = struct{}{}
	}
	for _, country := range filter.ManufacturerCountries {
		filterSets.ManufacturerCountries[country] = struct{}{}
	}
	for _, tag := range filter.Tags {
		filterSets.Tags[tag] = struct{}{}
	}

	return filterSets
}

func mapContains[T comparable](set map[T]any, val T) bool {
	for item, _ := range set {
		if item == val {
			return true
		}
	}
	return false
}

func sliceContains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
