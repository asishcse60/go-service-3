package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/asishcse60/service/business/data/store/user"
	"github.com/asishcse60/service/business/sys/database"
)

// Users retrieves all users from the database.
func Users(log *zap.SugaredLogger, cfg database.Config, pageNumber string, rowsPerPage string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()
	fmt.Println("User function call")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		return fmt.Errorf("converting page number: %w", err)
	}

	rows, err := strconv.Atoi(rowsPerPage)
	if err != nil {
		return fmt.Errorf("converting rows per page: %w", err)
	}

	store := user.NewStore(log, db)

	users, err := store.Query(ctx, page, rows)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(users)
}
