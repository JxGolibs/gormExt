package gormExt

import (
	"github.com/JxGolibs/responsePack"
	"github.com/jinzhu/gorm"
	"math"
)

//数据分页
// db = db.Model(&model.ParsedData{}).
// 	Where(where).
// 	Order("created_at DESC").
// 	Count(&page.Total).
// 	Limit(page.Limit). // Do limit and offset after count, but you now need Model before.
// 	Offset(page.Offset).
// 	Find(&parsedData)
func Paging(table interface{}, page ...*responsePack.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(page) == 0 || page[0] == nil {
			//无分页处理
			return db.Find(table)
		}

		pageOffset := ((page[0].PageNo - 1) * page[0].PageSize)

		//	page_no=1&page_size=10
		tmpModel := table
		_db := db.
			Model(tmpModel).
			Count(&page[0].TotalRecord).
			Limit(page[0].PageSize). // Do limit and offset after count, but you now need Model before.
			Offset(pageOffset).
			Find(table)
		page[0].TotalPages = int(math.Ceil(float64(page[0].TotalRecord) / float64(page[0].PageSize)))
		return _db
	}
}
