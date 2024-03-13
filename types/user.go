package types

import "github.com/google/uuid"

type AuthenticatedUser struct {
	ID       uuid.UUID
	Email    string
	LoggedIn bool
	Account
}

const UserContextKey = "user"
