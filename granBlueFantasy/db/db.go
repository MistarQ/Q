package db

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Mysql *gorm.DB

func InitDB() {

	// 连接表
	var err error
	Mysql, err = getDB(viper.GetString("mysql.Name"))
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	// 创建users表
	err = Mysql.Set("gorm:table_options", "Engine=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&User{})
	if err != nil {
		panic("建表失败, error=" + err.Error())
	}
}

func getDB(dbName string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		viper.GetString("mysql.User"),
		viper.GetString("mysql.Password"),
		viper.GetString("mysql.Host"),
		viper.GetInt("mysql.Port"),
		dbName,
		viper.GetString("mysql.TimeOut"),
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return db, err
}
