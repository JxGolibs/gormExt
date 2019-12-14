package gormExt

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Config struct {
	Adapter  string `default:"sqlite3"`
	Alias    string `default:""`
	Host     string `default:"localhost"`
	Port     int    `default:"3306"`
	DataBase string `default:"cloutAPI"`
	Encoding string `default:"utf8"`
	User     string `default:"root"`
	Password string `required:"true" env:"DBPassword"`
}

var (
	dbList = make(map[string]*gorm.DB)
)

func New(cnf Config, debug bool, models ...interface{}) *gorm.DB {
	var err error
	var jdb *gorm.DB
	switch cnf.Adapter {
	case "sqlite3":
		jdb, err = gorm.Open("sqlite3", cnf.DataBase+".db")
	case "mysql":
		jdb, err = gorm.Open("mysql", fmt.Sprint("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.DataBase, cnf.Encoding))
	}
	if err != nil {
		panic("连接数据库失败")
	}
	if debug {
		jdb.Debug()
	}
	if len(models) > 0 {
		// 自动迁移模式
		jdb.AutoMigrate(models)
	}

	dbList[cnf.Alias] = jdb
	return jdb
}

//根据别名获取数据库连接
func DB(alias ...string) *gorm.DB {
	if len(alias) > 0 {
		return dbList[alias[0]]
	}
	return dbList[""]
}

//关闭所有连接
func CloseAll() {
	for _, v := range dbList {
		if v != nil {
			v.Close()
		}
	}
}
