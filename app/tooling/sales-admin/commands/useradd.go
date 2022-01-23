package commands

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/asishcse60/service/business/data/store/user"
	"github.com/asishcse60/service/business/sys/auth"
	"github.com/asishcse60/service/business/sys/database"
)

// UserAdd adds new users into the database.
func UserAdd(log *zap.SugaredLogger, cfg database.Config, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store := user.NewStore(log, db)

	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	usr, err := store.Create(ctx, nu, time.Now())
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
