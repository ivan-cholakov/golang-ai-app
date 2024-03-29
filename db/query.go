package db

import (
	"context"
	"dreampicai/types"

	"github.com/google/uuid"
)

func CreateAcount(account *types.Account) error {
	_, err := Bun.NewInsert().Model(account).Exec(context.Background())
	return err
}

func GetAccountByUserID(userID uuid.UUID) (types.Account, error) {
	var account types.Account
	err := Bun.NewSelect().Model(&account).Where("user_id = ?", userID).Scan(context.Background())
	return account, err
}
