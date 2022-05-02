package models

type User struct {
	ID       int64
	UserName string
	IpInfo   []IpInfo `gorm:"many2many:user_ip_infos;"`
	Admin    bool
}
