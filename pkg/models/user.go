package models

type User struct {
	ID       uint
	UserName string
	IpInfo   []IpInfo `gorm:"many2many:user_ip_infos;"`
	Admin    bool
}
