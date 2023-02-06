package authorize

type User struct {
	Verified bool `json:"verified"`
	ID       uint `gorm:"primarykey"`
}
