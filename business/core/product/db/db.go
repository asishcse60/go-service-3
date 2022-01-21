// Package db contains product related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/asishcse60/service/business/sys/database"
)

// Store manages the set of APIs for user access.
type Store struct {
	log          *zap.SugaredLogger
	tr           database.Transactor
	db           sqlx.ExtContext
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		tr:  db,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		fn(s.db)
	}
	return database.WithinTran(ctx, s.log, s.tr, fn)
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// Create adds a Product to the database. It returns the created Product with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, prd Product) error {
	if err := database.NamedExecContext(ctx, s.log, s.db, ProductCreateQuery, prd); err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}
	return nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (s Store) Update(ctx context.Context, prd Product) error {
	if err := database.NamedExecContext(ctx, s.log, s.db, ProductUpdateQuery, prd); err != nil {
		return fmt.Errorf("updating product productID[%s]: %w", prd.ID, err)
	}

	return nil
}

// Delete removes the product identified by a given ID.
func (s Store) Delete(ctx context.Context, productID string) error {
	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID,
	}

	if err := database.NamedExecContext(ctx, s.log, s.db, ProductDeleteQuery, data); err != nil {
		return fmt.Errorf("deleting product productID[%s]: %w", productID, err)
	}

	return nil
}

// Query gets all Products from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Product, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	var prds []Product
	if err := database.NamedQuerySlice(ctx, s.log, s.db, ProductListQuery, data, &prds); err != nil {
		return nil, fmt.Errorf("selecting products: %w", err)
	}

	return prds, nil
}

// QueryByID finds the product identified by a given ID.
func (s Store) QueryByID(ctx context.Context, productID string) (Product, error) {
	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID,
	}

	var prd Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, ProductIDQuery, data, &prd); err != nil {
		return Product{}, fmt.Errorf("selecting product productID[%q]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID finds the product identified by a given User ID.
func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Product, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	var prds []Product
	if err := database.NamedQuerySlice(ctx, s.log, s.db, ProductUserIDQuery, data, &prds); err != nil {
		return nil, fmt.Errorf("selecting products userID[%s]: %w", userID, err)
	}

	return prds, nil
}