package models

import (
	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

func Transactional(ts ...func(db2 *gorm.DB) error) error {
	if ts == nil || len(ts) == 0 {
		return constants.ErrNoTransaction
	}
	tx := db.DB.Begin()

	for _, t := range ts {
		if err := t(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func CusDelete(db *gorm.DB, model interface{}, id interface{}, user string) *gorm.DB {
	if CheckUserField(model) {
		return db.Model(model).Where("id = ?", id).Updates(map[string]interface{}{
			"user":       user,
			"deleted_at": time.Now(),
		})
	}

	return db.Delete(&model, "id = ?", id)
}

func CheckUserField(data interface{}) bool {
	fields := structs.Names(data)
	for _, field := range fields {
		if string(field) == constants.UserFieldStr {
			return true
		}
	}
	return false
}

func CreateNormalHistoryByModel(model interface{}, action, user, from string) (NormalHistory, error) {
	nh := CreateNormalHistory(action, user, from)

	err := nh.SetDataFromModel(model)

	return nh, err
}

func DeleteById(model, id interface{}) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		return tx.Delete(model, "id = ?", id).Error
	}
}

func SqlLikeString(str string) string {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return ""
	}
	return "%" + str + "%"
}

func GetFilterDateValue(startDate, endDate time.Time) (time.Time, time.Time) {
	if startDate.IsZero() {
		startDate = endDate
	}
	if endDate.IsZero() {
		endDate = startDate
	}
	endDate = endDate.AddDate(0, 0, 1)
	return startDate, endDate
}
