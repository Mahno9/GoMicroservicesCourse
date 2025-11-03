package part

import (
	"context"

	"github.com/samber/lo"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filters *domainModel.PartsFilter) ([]*domainModel.Part, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	parts := make([]*domainModel.Part, 0)
	repoFilters := converter.DomainToRepoFilter(filters)

	for uuid, part := range r.parts {
		if len(repoFilters.Uuids) > 0 && !lo.Contains(repoFilters.Uuids, uuid) {
			continue
		}
		if len(repoFilters.Names) > 0 && !lo.Contains(repoFilters.Names, part.Name) {
			continue
		}
		if len(repoFilters.Categories) > 0 && !lo.Contains(repoFilters.Categories, part.Category) {
			continue
		}
		if len(repoFilters.ManufacturerCountries) > 0 && !lo.Contains(repoFilters.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if len(repoFilters.Tags) > 0 {
			for _, tag := range repoFilters.Tags {
				if lo.Contains(part.Tags, tag) {
					parts = append(parts, converter.RepoToDomainPart(part))
					continue
				}
			}
			continue
		}

		// empty filters case
		parts = append(parts, converter.RepoToDomainPart(part))
	}

	return parts, nil
}
