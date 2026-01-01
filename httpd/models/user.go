package models

type User struct {
	BaseModel
	Name  string `gorm:"column:name;type:varchar(64)" json:"name" form:"name" binding:"required"`
	Phone string `gorm:"column:phone;type:varchar(32)" json:"phone" form:"phone" binding:""`
	Email string `gorm:"column:email;type:varchar(64)" json:"email" form:"email" binding:"required,email"`
}
