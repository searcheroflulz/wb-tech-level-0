package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wb-tech-level-0/internal/model"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) AddOrder(ctx context.Context, order model.Order) error {
	conn, err := p.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.CustomerID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return err
	}

	err = p.addDelivery(ctx, order.Delivery, order.OrderUID)
	if err != nil {
		return err
	}

	err = p.addPayment(ctx, order.Payment, order.OrderUID)
	if err != nil {
		return err
	}

	for _, item := range order.OrderItems {
		err = p.addOrderItem(ctx, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) addDelivery(ctx context.Context, delivery model.Delivery, orderUid string) error {
	conn, err := p.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO NOTHING`,
		orderUid,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) addPayment(ctx context.Context, payment model.Payment, orderUid string) error {
	conn, err := p.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`,
		orderUid,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDate,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) addOrderItem(ctx context.Context, item model.OrderItem) error {
	conn, err := p.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx,
		`INSERT INTO order_items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.RID,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NMID,
		item.Brand,
		item.Status,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Orders(ctx context.Context) ([]model.Order, error) {
	conn, err := p.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	orders := make([]model.Order, 1)

	if err := conn.SelectContext(ctx, &orders, `SELECT * FROM orders`); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, nil
	}

	for i := range orders {
		err = p.delivery(ctx, &orders[i].Delivery, conn, orders[i].OrderUID)
		if err != nil {
			return nil, err
		}
		err = p.payment(ctx, &orders[i].Payment, conn, orders[i].OrderUID)
		if err != nil {
			return nil, err
		}
		err = p.items(ctx, &orders[i].OrderItems, conn, orders[i].TrackNumber)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (p *Postgres) delivery(ctx context.Context, delivery *model.Delivery, conn *sqlx.Conn, id string) error {
	if err := conn.GetContext(ctx, delivery, `SELECT * FROM delivery WHERE order_uid = $1`, id); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) payment(ctx context.Context, payment *model.Payment, conn *sqlx.Conn, id string) error {
	if err := conn.GetContext(ctx, payment, `SELECT * FROM payment WHERE order_uid = $1`, id); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) items(ctx context.Context, items *[]model.OrderItem, conn *sqlx.Conn, id string) error {
	if err := conn.SelectContext(ctx, items, `SELECT * FROM order_items WHERE track_number = $1`, id); err != nil {
		return err
	}

	return nil
}
