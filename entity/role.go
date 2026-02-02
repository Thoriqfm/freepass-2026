package entity

type Role struct {
	RoleID   int    `json:"role_id" gorm:"type:smallint;primaryKey;autoIncrement"`
	RoleName string `json:"role_name" gorm:"type:varchar(255);not null"`

	Users []User `json:"users" gorm:"foreignKey:RoleID"`
}
