package mysqlclient

import (
	"fmt"
	"go-demo/internal/config"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ormDB *gorm.DB
	once  sync.Once
)

func InitMysql(conf *config.MysqlConfig) error {
	var err error
	// 使用 once 确保只初始化一次
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.DbName)
		ormDB, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{})
		if err != nil {
			return
		}
		// 设置连接池
		sqlDB, err := ormDB.DB()
		if err != nil {
			return
		}
		sqlDB.SetMaxIdleConns(conf.MaxConn) //设置最大连接数
		sqlDB.SetMaxOpenConns(conf.MaxOpen) //设置最大的空闲连接数
	})
	return err
}

func GetMysqlClient() *gorm.DB {
	return ormDB
}
