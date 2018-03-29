package templateFuncs

import (
	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/controllers/pagination"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/nicksnyder/go-i18n/i18n"
	"reflect"
	"time"
)

var ActionNameDict = map[int]string{
	models.AddLevel:   "Stock",
	models.MinusLevel: "Withdraw",
}

func LastIndex(x int, a interface{}) bool {
	return x == reflect.ValueOf(a).Len()-1
}

func DateTime24(dt time.Time) string {
	if dt.IsZero() {
		return ""
	}
	return dt.Format("2006-01-02 15:04:05")
}

func Date(dt time.Time) string {
	if dt.IsZero() {
		return ""
	}
	return dt.Format("2006-01-02")
}

func ActionName(action int) string {
	return "t" + ActionNameDict[action]
}

func I18nTranslate(key, lang string) string {
	defaultLang := constants.DefaultLang
	T, _ := i18n.Tfunc(lang, defaultLang)

	return T(key)
}

func GetPageLi(pg pagination.Pagination) []int {
	if pg.Total == 0 {
		return nil
	}
	var startNum, endNum int
	if pg.Current-5 < 1 {
		startNum = 1
	} else {
		startNum = pg.Current - 5
	}
	if startNum+9 > pg.Total {
		endNum = pg.Total
	} else {
		endNum = startNum + 9
	}
	numArr := make([]int, 0)
	for i := startNum; i <= endNum; i++ {
		numArr = append(numArr, i)
	}
	return numArr
}
