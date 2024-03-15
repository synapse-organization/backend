package models

import (
	"github.com/spf13/cast"
)

type Postgres struct {
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
}

func (p *Postgres) GetPostgresURL() string {
	return "postgres://" + p.UserName + ":" + p.Password + "@" + p.Host + ":" + cast.ToString(p.Port) + "/" + p.DbName
}
