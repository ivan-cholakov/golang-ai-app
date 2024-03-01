package types

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}

const UserContextKey = "user"
