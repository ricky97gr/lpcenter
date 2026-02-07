package pagination

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ricky97gr/lpcenter/server/utils"
)

type PageQuery struct {
	Page       int         `json:"page" form:"page"`
	PageSize   int         `json:"pageSize" form:"pageSize"`
	Conditions []Condition `json:"conditions" form:"conditions"`
	Sorts      []Sort      `json:"sorts" form:"sorts"`
}

type Condition struct {
	Field     string      `json:"field" form:"field"`
	Value     interface{} `json:"value" form:"value"`
	Operation int         `json:"operation" form:"operation"`
}

type Sort struct {
	Field   string `json:"field" form:"field"`
	OrderBy int    `json:"orderBy" form:"orderBy"`
}

const (
	Equal        = iota + 1
	NotEqual
	GreaterThan
	GreaterEqual
	LessThan
	LessEqual
	Like
	In
	NotIn
)

const (
	Asc  = 1
	Desc = -1
)

func GetPageQuery(ctx *gin.Context) (PageQuery, error) {
	var page PageQuery
	var s []Sort
	var conditions []Condition
	var pageNumber int = 1
	var pageSize int = 10

	if orderStr, ok := ctx.GetQuery("sorts"); ok {
		err := json.Unmarshal([]byte(orderStr), &s)
		if err != nil {
			utils.Logger.Errorw("failed to unmarshal sorts", "error", err)
		}
	}

	if conStr, ok := ctx.GetQuery("conditions"); ok {
		err := json.Unmarshal([]byte(conStr), &conditions)
		if err != nil {
			utils.Logger.Errorw("failed to unmarshal conditions", "error", err)
		}
	}

	if pageNumberStr, ok := ctx.GetQuery("page"); ok {
		n, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			utils.Logger.Errorw("failed to get page number", "error", err)
			n = 1
		}
		pageNumber = n
	}

	if pageSizeStr, ok := ctx.GetQuery("pageSize"); ok {
		n, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			utils.Logger.Errorw("failed to get page size", "error", err)
			n = 10
		}
		pageSize = n
	}

	page.Page = pageNumber
	page.PageSize = pageSize
	page.Sorts = s
	page.Conditions = conditions

	return page, nil
}

func (p *PageQuery) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PageQuery) GetLimit() int {
	return p.PageSize
}

func Order(sorts []Sort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		resultDB := db
		for _, sort := range sorts {
			if sort.OrderBy == Asc {
				resultDB = resultDB.Order(sort.Field + " asc")
			}
			if sort.OrderBy == Desc {
				resultDB = resultDB.Order(sort.Field + " desc")
			}
		}
		return resultDB
	}
}

func (q PageQuery) GetCondition(field string) (Condition, bool) {
	for _, cond := range q.Conditions {
		if cond.Field == field {
			return cond, true
		}
	}
	return Condition{}, false
}

func ParseQuery(q PageQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		resultDB := db
		for _, v := range q.Conditions {
			resultDB = resultDB.Scopes(QueryFilter(v.Field, v.Value, v.Operation))
		}
		resultDB = resultDB.Scopes(Order(q.Sorts))
		resultDB = resultDB.Scopes(QueryLimitShip(q.Page, q.PageSize))
		return resultDB
	}
}

func QueryFilter(field string, value interface{}, operation int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch operation {
		case Equal:
			return db.Where(field+" = ?", value)
		case NotEqual:
			return db.Where(field+" != ?", value)
		case GreaterThan:
			return db.Where(field+" > ?", value)
		case GreaterEqual:
			return db.Where(field+" >= ?", value)
		case LessThan:
			return db.Where(field+" < ?", value)
		case LessEqual:
			return db.Where(field+" <= ?", value)
		case Like:
			return db.Where(field+" like ?", "%"+value.(string)+"%")
		case In:
			return db.Where(field+" in ?", value)
		case NotIn:
			return db.Where(field+" not in ?", value)
		default:
			return db
		}
	}
}

func QueryLimitShip(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
