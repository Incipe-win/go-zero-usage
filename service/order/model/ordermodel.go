package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
		FindOneByOrderId(ctx context.Context, orderId uint64) (*Order, error)
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

func (c *customOrderModel) FindOneByOrderId(ctx context.Context, orderId uint64) (*Order, error) {
	db3UserUserIdKey := fmt.Sprintf("%s%v", cacheSqlTestOrderIdPrefix, orderId)
	var resp Order
	err := c.QueryRowIndexCtx(ctx, &resp, db3UserUserIdKey, c.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `order_id` = ? limit 1", orderRows, c.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, orderId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, c.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c, opts...),
	}
}
