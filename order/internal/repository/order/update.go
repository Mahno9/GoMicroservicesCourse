package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order *model.Order) error {
	modelOrder := converter.ModelToRepositoryOrder(order)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("order_uuid", modelOrder.OrderUuid).
		Set("user_uuid", modelOrder.UserUuid).
		Set("part_uuids", modelOrder.PartUuids).
		Set("transaction_uuid", modelOrder.TransactionUuid).
		Set("payment_method", modelOrder.PaymentMethod).
		Set("order_status", modelOrder.Status).
		Where(sq.Eq{"order_uuid": modelOrder.OrderUuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("❗ [Update] Failed to build query: %v\n", err)
		return model.ErrQueryBuild
	}

	_, err = r.dbConnPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("❗ [Update] Failed to execute query: %v\n", err)
		return model.ErrQueryExecution
	}

	return nil
}
