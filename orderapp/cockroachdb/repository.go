package cockroachdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"github.com/cockroachdb/cockroach-go/crdb"
	_ "github.com/lib/pq"

	"github.com/shijuvar/service-weaver/orderapp/model"
)

type Repository interface {
	CreateOrder(context.Context, model.Order) error
	GetOrderByID(context.Context, string) (model.Order, error)
}

type repository struct {
	weaver.Implements[Repository]
	weaver.WithConfig[config]
	db *sql.DB
}

func (repo *repository) Init(context.Context) error {
	cfg := repo.Config()
	if err := cfg.Validate(); err != nil {
		repo.Logger().Error("error:", err)
	}
	db, err := sql.Open(cfg.Driver, cfg.Source)
	repo.Logger().Info("connected to DB")

	if err != nil {
		return err
	}
	repo.db = db
	return nil
}

type config struct {
	Driver string //`toml:"Driver"` -> Name of the database driver.
	Source string //`toml:"Source"` -> Database server source URI.
}

func (cfg *config) Validate() error {
	if len(cfg.Driver) == 0 {
		return errors.New("DB driver is not provided")
	}
	if len(cfg.Source) == 0 {
		return errors.New("DB source is not provided")
	}
	return nil
}

// CreateOrder persist Order data into the query model
func (repo *repository) CreateOrder(ctx context.Context, order model.Order) error {

	// Run a transaction to sync the query model.
	err := crdb.ExecuteTx(ctx, repo.db, nil, func(tx *sql.Tx) error {
		return createOrder(tx, order)
	})
	if err != nil {
		return err
	}
	return nil
}

// GetOrderByID query the Orders by given id
func (repo *repository) GetOrderByID(ctx context.Context, id string) (model.Order, error) {
	var orderRow = model.Order{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT id, customerid, status, createdon, restaurantid FROM orders WHERE id = $1",
		id).
		Scan(
			&orderRow.ID, &orderRow.CustomerID, &orderRow.Status, &orderRow.CreatedOn, &orderRow.RestaurantId,
		); err != nil {
		return orderRow, err
	}
	return orderRow, nil
}

// GetOrderItems query the order items by given order id
func (repo *repository) GetOrderItems(ctx context.Context, id string) ([]model.OrderItem, error) {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT code, name, unitprice, quantity FROM orderitems WHERE orderid = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// An OrderItem slice to hold data from returned rows.
	var oitems []model.OrderItem

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ProductCode, &item.Name, &item.UnitPrice,
			&item.Quantity); err != nil {
			return oitems, err
		}
		oitems = append(oitems, item)
	}
	if err = rows.Err(); err != nil {
		return oitems, err
	}
	return oitems, nil
}
func createOrder(tx *sql.Tx, order model.Order) error {

	// Insert into the "orders" table.
	sql := `
			INSERT INTO orders (id, customerid, status, createdon, restaurantid, amount)
			VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := tx.Exec(sql, order.ID, order.CustomerID, order.Status, order.CreatedOn, order.RestaurantId, order.Amount)
	if err != nil {
		return err
	}
	// Insert items into the "orderitems" table.
	// Because it's store for read model, we can insert denormalized data
	for _, v := range order.OrderItems {
		sql = `
			INSERT INTO orderitems (orderid, code, name, unitprice, quantity)
			VALUES ($1,$2,$3,$4,$5)`

		_, err := tx.Exec(sql, order.ID, v.ProductCode, v.Name, v.UnitPrice, v.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}
