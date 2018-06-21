package goes

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func MigrateEventsTable() error {
	var err error

	DB.DropTable(&EventDB{})
	if err = DB.AutoMigrate(&EventDB{}).Error; err != nil {
		return err
	}
	return nil
}

func InitDB(dbConn string, logMode bool) error {
	var err error

	DB, err = gorm.Open("postgres", dbConn)
	if err != nil {
		return err
	}
	DB.LogMode(logMode)
	return nil
}
