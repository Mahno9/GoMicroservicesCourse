package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUuid uuid.UUID) (*model.Order, error) {
	builderSelect := sq.Select("*").
		From("orders").
		Where(sq.Eq{"order_uuid": orderUuid.String()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("❗ [Get] Failed to build query: %v\n", err)
		return nil, model.ErrQueryBuild
	}

	rows, err := r.dbConnPool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("❗ [Get] Failed to execute query: %v\n", err)
		return nil, model.ErrQueryExecution
	}
	defer rows.Close()

	repoOrder, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.Order])
	if err != nil {
		log.Printf("❗ [Get] Failed to scan row into struct: %v\n", err)
		return nil, model.ErrQueryResponseScanning
	}

	return converter.RepositoryToModelOrder(&repoOrder), nil
}
