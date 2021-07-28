package domain

type SignInUser struct {
	ClientID string
	UserID   string
	Password string
}

type SignUpUser struct {
	ClientID   string
	UserID     string
	Password   string
	UserName   string
	UserCardID string
	Email      string
	Phone      string
}
