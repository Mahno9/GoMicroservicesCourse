package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	newOrder := converter.ModelToRepositoryOrder(order)

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_uuids", "total_price", "payment_method", "order_status").
		Values(newOrder.UserUuid, newOrder.PartUuids, newOrder.TotalPrice, newOrder.PaymentMethod, newOrder.Status).
		Suffix("RETURNING *")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("❗ [Create] Failed to build query: %v\n", err)
		return nil, model.ErrQueryBuild
	}

	rows, err := r.dbConnPool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("❗ [Create] Failed to execute query: %v\n", err)
		return nil, model.ErrQueryExecution
	}
	defer rows.Close()

	repoOrderResult, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.Order])
	if err != nil {
		log.Printf("❗ [Create] Failed to scan row into struct: %v\n", err)
		return nil, model.ErrQueryResponseScanning
	}

	return converter.RepositoryToModelOrder(&repoOrderResult), nil
}
