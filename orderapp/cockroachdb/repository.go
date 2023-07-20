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
}

type repository struct {
	weaver.Implements[Repository]
	weaver.WithConfig[config]
	db *sql.DB
}

func (repo *repository) Init(context.Context) error {
	cfg := repo.Config()
	if err := cfg.Validate(); err != nil {
		repo.Logger().Error("error", err)
	}
	db, err := sql.Open(cfg.Driver, cfg.Source)
	repo.Logger().Info("connected sql")

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
