package repositories

import (
	"context"
	"fmt"
	"github.com/IlyaZayats/servord/internal/entities"
	"github.com/IlyaZayats/servord/internal/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(db *pgxpool.Pool) (interfaces.OrderRepository, error) {
	return &PostgresOrderRepository{
		db: db,
	}, nil
}

func (r *PostgresOrderRepository) GetOrders() ([]entities.Order, error) {
	var orders []entities.Order
	q := `select json_build_object(
               'order_uid', W.order_uid,
               'track_number', W.track_number,
               'entry', W.entry,
               'locale', W.locale, 
               'internal_signature', W.internal_signature, 
               'customer_id', W.customer_id, 
               'delivery_service', W.delivery_service, 
               'shard_key', W.shard_key,
               'sm_id', W.sm_id,
               'date_created', W.date_created,
               'oof_shard', W.oof_shard,
               'delivery', (select row_to_json(x) from (select D.name, D.phone, D.zip, D.city, D.address, D.region, D.email from Delivery as D where D.order_uid = W.order_uid) as x),
               'payment', (select row_to_json(x) from (select P.order_uid as transaction, P.request_id, P.currency, P.provider, P.amount, P.payment_dt, P.Bank, P.delivery_cost, P.goods_total, P.custom_fee from Payment as P where P.order_uid = W.order_uid) as x),
               'items', array(select row_to_json(x) from (select I.track_number, I.rid, I.name, I.size, I.brand, I.chrt_id, I.price, I.sale, I.total_price, I.nm_id, i.status from Item as I where I.track_number = W.track_number) as x)
       ) from WBOrder as W;`
	rows, err := r.db.Query(context.Background(), q)
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *PostgresOrderRepository) InsertOrder(order entities.Order) error {
	logrus.Println(order)
	q := `insert into WBOrder (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	if _, err := r.db.Exec(context.Background(), q, order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard); err != nil {
		return errors.Wrap(err, "insert order")
	}
	q = `insert into Delivery (order_uid, name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := r.db.Exec(context.Background(), q, order.OrderUid, order.Del.Name, order.Del.Phone, order.Del.Zip, order.Del.City, order.Del.Address, order.Del.Region, order.Del.Email); err != nil {
		return errors.Wrap(err, "insert delivery")
	}
	q = `insert into Payment (order_uid, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if _, err := r.db.Exec(context.Background(), q, order.OrderUid, order.Pay.RequestId, order.Pay.Currency, order.Pay.Provider, order.Pay.Amount, order.Pay.PaymentDt, order.Pay.Bank, order.Pay.DeliveryCost, order.Pay.GoodsTotal, order.Pay.CustomFee); err != nil {
		return errors.Wrap(err, "insert payment")
	}
	for i, item := range order.Items {
		q := `insert into Item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
		if _, err := r.db.Exec(context.Background(), q, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status); err != nil {
			return errors.Wrap(err, fmt.Sprintf("insert item[%v]", i))
		}
	}
	return nil
}
