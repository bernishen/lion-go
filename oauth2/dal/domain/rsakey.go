package domain

// RSAKey : Record a rsa information,
type RSAKey struct {
	BaseModel
	Scope      string `gorm:"type:varchar(20)"`
	Version    int    `gorm:"not_null"`
	PublicKey  string `gorm:"type:varchar(500)"`
	PrivateKey string `gorm:"type:varchar(1000)"`
}
