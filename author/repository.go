package author

import (
	"context"

	"github.com/williamchand/my-wallet/models"
)

// Repository represent the author's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.Author, error)
}
