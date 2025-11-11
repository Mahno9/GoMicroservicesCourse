package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	dbConnPool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		dbConnPool: pool,
	}
}
