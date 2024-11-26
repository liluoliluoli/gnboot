package gomysql

import (
	"github.com/liluoliluoli/gnboot/internal/utils/context_util"
	"gorm.io/gorm"
	"time"
)

func CreateUserOperator(db *gorm.DB) {
	createBy := ""
	updateBy := ""
	defer func() {
		if r := recover(); r != nil {
			updateBy = ""
			createBy = ""
		}
	}()
	createBy = context_util.GetOperatorUid(db.Statement.Context)
	updateBy = context_util.GetOperatorUid(db.Statement.Context)

	currentTime := time.Now()
	db.Statement.SetColumn("create_time", currentTime)
	db.Statement.SetColumn("update_time", currentTime)
	db.Statement.SetColumn("create_by", createBy)
	db.Statement.SetColumn("update_by", updateBy)
}

func UpdateUserOperator(db *gorm.DB) {
	updateBy := ""
	defer func() {
		if r := recover(); r != nil {
			updateBy = ""
		}
	}()

	updateBy = context_util.GetOperatorUid(db.Statement.Context)

	currentTime := time.Now()
	db.Statement.SetColumn("update_time", currentTime)
	db.Statement.SetColumn("update_by", updateBy)

}
