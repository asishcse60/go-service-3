// Package product provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package product

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/asishcse60/service/business/sys/auth"
	"github.com/asishcse60/service/business/sys/database"
	"github.com/asishcse60/service/business/sys/validate"
)

// Store manages the set of API's for product access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs a product store for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create adds a Product to the database. It returns the created Product with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, claims auth.Claims, np NewProduct, now time.Time) (Product, error) {
	if err := validate.Check(np); err != nil {
		return Product{}, fmt.Errorf("validating data: %w", err)
	}

	prd := Product{
		ID:          validate.GenerateID(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      claims.Subject,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := database.NamedExecContext(ctx, s.log, s.db, CreateProductQuery, prd); err != nil {
		return Product{}, fmt.Errorf("inserting product: %w", err)
	}

	return prd, nil
}

// Update modifies data about a Product. It will error if the specified ID is
// invalid or does not reference an existing Product.
func (s Store) Update(ctx context.Context, claims auth.Claims, productID string, up UpdateProduct, now time.Time) error {
	if err := validate.CheckID(productID); err != nil {
		return database.ErrInvalidID
	}
	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	prd, err := s.QueryByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("updating product productID[%s]: %w", productID, err)
	}

	// If you are not an admin and looking to retrieve someone elses product.
	if !claims.Authorized(auth.RoleAdmin) && prd.UserID != claims.Subject {
		return database.ErrForbidden
	}

	if up.Name != nil {
		prd.Name = *up.Name
	}
	if up.Cost != nil {
		prd.Cost = *up.Cost
	}
	if up.Quantity != nil {
		prd.Quantity = *up.Quantity
	}
	prd.DateUpdated = now

	if err := database.NamedExecContext(ctx, s.log, s.db, UpdateProductQuery, prd); err != nil {
		return fmt.Errorf("updating product productID[%s]: %w", productID, err)
	}

	return nil
}

// Delete removes the product identified by a given ID.
func (s Store) Delete(ctx context.Context, claims auth.Claims, productID string) error {
	if err := validate.CheckID(productID); err != nil {
		return database.ErrInvalidID
	}

	// If you are not an admin.
	if !claims.Authorized(auth.RoleAdmin) {
		return database.ErrForbidden
	}

	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID,
	}

	if err := database.NamedExecContext(ctx, s.log, s.db, DeleteProductQuery, data); err != nil {
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

	var products []Product
	if err := database.NamedQuerySlice(ctx, s.log, s.db, ListProductQuery, data, &products); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting products: %w", err)
	}

	return products, nil
}

// QueryByID finds the product identified by a given ID.
func (s Store) QueryByID(ctx context.Context, productID string) (Product, error) {
	if err := validate.CheckID(productID); err != nil {
		return Product{}, database.ErrInvalidID
	}

	data := struct {
		ProductID string `db:"product_id"`
	}{
		ProductID: productID,
	}

	var prd Product
	if err := database.NamedQueryStruct(ctx, s.log, s.db, IDProductQuery, data, &prd); err != nil {
		if err == database.ErrNotFound {
			return Product{}, database.ErrNotFound
		}
		return Product{}, fmt.Errorf("selecting product productID[%q]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID finds the product identified by a given User ID.
func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Product, error) {
	if err := validate.CheckID(userID); err != nil {
		return nil, database.ErrInvalidID
	}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	var products []Product
	if err := database.NamedQuerySlice(ctx, s.log, s.db, UpdateProductQuery, data, &products); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting products userID[%s]: %w", userID, err)
	}

	return products, nil
}
