package service

import (
	"database/sql"
	"fmt"
	"net/url"
	"seconda/cmd/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBDecorator struct {
	GormInterface   *gorm.DB
	NativeInterface *sql.DB
}

func InitORM(config config.DataBaseConfigInterface) *DBDecorator {
	dsn := dsn(config)

	ORM, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connect failed: %w", err))
	}

	nativeDB, err := ORM.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get native DB interface: %w", err))
	}

	dbd := DBDecorator{
		GormInterface:   ORM,
		NativeInterface: nativeDB,
	}

	return &dbd
}

func dsn(config config.DataBaseConfigInterface) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User(),
		config.Password(),
		config.Host(),
		config.Port(),
		config.Name(),
	)

	if tz := config.Timezone(); tz != "" {
		dsn += "&loc=" + url.QueryEscape(tz)
	}

	return dsn
}

func (dbd DBDecorator) CloseDB() {
	err := dbd.NativeInterface.Close()
	if err != nil {
		panic(fmt.Errorf("db close error: %w", err))
	}
}

func (dbd DBDecorator) GDB() *gorm.DB {
	return dbd.GormInterface
}

func (dbd DBDecorator) NativeDB() *sql.DB {
	return dbd.NativeInterface
}
