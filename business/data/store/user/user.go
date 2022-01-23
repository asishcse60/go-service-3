// Package user contains user related CRUD functionality.
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/asishcse60/service/business/sys/auth"
	"github.com/asishcse60/service/business/sys/database"
	"github.com/asishcse60/service/business/sys/validate"
)

// Store manages the set of API's for user access.
type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

// NewStore constructs a user store for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	usr := User{
		ID:           validate.GenerateID(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := database.NamedExecContext(ctx, s.log, s.db, CreateUserQuery, usr); err != nil {
		return User{}, fmt.Errorf("inserting user: %w", err)
	}

	return usr, nil
}

// Update replaces a user document in the database.
func (s Store) Update(ctx context.Context, claims auth.Claims, userID string, uu UpdateUser, now time.Time) error {
	if err := validate.CheckID(userID); err != nil {
		return database.ErrInvalidID
	}
	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	usr, err := s.QueryByID(ctx, claims, userID)
	if err != nil {
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHash = pw
	}
	usr.DateUpdated = now

	if err := database.NamedExecContext(ctx, s.log, s.db, UpdateUserQuery, usr); err != nil {
		return fmt.Errorf("updating userID[%s]: %w", userID, err)
	}

	return nil
}

// Delete removes a user from the database.
func (s Store) Delete(ctx context.Context, claims auth.Claims, userID string) error {
	if err := validate.CheckID(userID); err != nil {
		return database.ErrInvalidID
	}

	// If you are not an admin and looking to delete someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return database.ErrForbidden
	}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	if err := database.NamedExecContext(ctx, s.log, s.db, DeleteUserQuery, data); err != nil {
		return fmt.Errorf("deleting userID[%s]: %w", userID, err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]User, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	var users []User
	if err := database.NamedQuerySlice(ctx, s.log, s.db, ListUserQuery, data, &users); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return users, nil
}

// QueryByID gets the specified user from the database.
func (s Store) QueryByID(ctx context.Context, claims auth.Claims, userID string) (User, error) {
	if err := validate.CheckID(userID); err != nil {
		return User{}, database.ErrInvalidID
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return User{}, database.ErrForbidden
	}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	var usr User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, IDUserQuery, data, &usr); err != nil {
		if err == database.ErrNotFound {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("selecting userID[%q]: %w", userID, err)
	}

	return usr, nil
}

// QueryByEmail gets the specified user from the database by email.
func (s Store) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (User, error) {

	// Add Email Validate function in validate
	// if err := validate.Email(email); err != nil {
	// 	return User{}, ErrInvalidEmail
	// }

	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	var usr User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, EmailUserQuery, data, &usr); err != nil {
		if err == database.ErrNotFound {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != usr.ID {
		return User{}, database.ErrForbidden
	}

	return usr, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (s Store) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	var usr User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, EmailUserQuery, data, &usr); err != nil {
		if err == database.ErrNotFound {
			return auth.Claims{}, database.ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("selecting user[%q]: %w", email, err)
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, database.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   usr.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	return claims, nil
}
