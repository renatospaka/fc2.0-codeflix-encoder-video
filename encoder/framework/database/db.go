package database

import (
	"encoder/domain"
	"log"

	"github.com/go-gorm/gorm"
	_ "github.com/go-gorm/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *						gorm.DB
	Dsn							string
	DsnTest					string
	DbType					string
	DbTypeTest			string
	Debug						bool
	AutoMigrateDb		bool
	Env							string
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	if d.Env != "test" {
		d.DB, err = gorm.Open(d.DbType, d.Dsn)
	} else {
		d.DB, err = gorm.Open(d.DbTypeTest, d.DsnTest)
	}

	if err != nil {
		return nil, err
	} 

	if d.Debug {
		d.DB.LogMode(true)
	}

	if d.DB.AutoMigrateDb {
		d.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.DB.Model(domain.Job{}).AddForeignKey(field: "video_id", dest: "videos (id)", onDelete: "CASCADE", onUpdate: "CASCADE")
	}
	return &d.DB, nil
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "Test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()
	if err != nil {
		log.Fatalf("Test db error: %w", err)
	} 
	return connection
}