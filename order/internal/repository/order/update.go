package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order *model.Order) error {
	repoOrder := converter.ModelToRepositoryOrder(order)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("user_uuid", repoOrder.UserUuid).
		Set("part_uuids", repoOrder.PartUuids).
		Set("transaction_uuid", repoOrder.TransactionUuid).
		Set("payment_method", repoOrder.PaymentMethod).
		Set("order_status", repoOrder.Status).
		Where(sq.Eq{"order_uuid": repoOrder.OrderUuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("❗ [Update] Failed to build query: %v\n", err)
		return model.ErrQueryBuild
	}

	_, err = r.dbConnPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("❗ [Update] Failed to execute query: %v\nQUERY: [%s}\nARGS: %+v\n", err, query, args)
		return model.ErrQueryExecution
	}

	return nil
}
