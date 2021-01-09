package User

type Profile struct {
	ID          uint   `gorm:"primary_key" db:"id"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	IsSuperuser bool   `gorm:"default:'0'" db:"is_superuser"`
}
