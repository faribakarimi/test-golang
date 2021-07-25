package models

type UserItems struct {
	ID		int	`gorm:"primary_key;auto_increment" json:"id"`
	UserID	int	`gorm:"not null" json:"user_id"`
	ItemID	int	`gorm:"not null" json:"item_id"`
}