package query

import (
	"strconv"
	"strings"

	"github.com/gophab/gophrame/core/form"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Sort struct {
	By        string
	Direction string
}

func (s *Sort) String() string {
	return s.By + " " + s.Direction
}

func (s Sort) For(v string) Sort {
	segs := strings.Split(v, " ")
	by := segs[0]
	dir := "ASC"
	if len(segs) > 1 {
		dir = segs[1]
	}
	return Sort{By: by, Direction: dir}
}

type Pageable interface {
	GetPage() int
	GetSize() int
	GetSort() []Sort
	GetLimit() int
	GetOffset() int
	NoSort() bool
}

type Pagination struct {
	Page int    `json:"page"`
	Size int    `json:"size"`
	Sort []Sort `json:"sort"`
}

func (p *Pagination) GetPage() int {
	if p.Page <= 1 {
		return 1
	} else {
		return p.Page
	}
}

func (p *Pagination) GetSize() int {
	if p.Size <= 0 {
		return 20
	} else {
		return p.Size
	}
}

func (p *Pagination) GetSort() []Sort {
	return p.Sort
}

func (p *Pagination) NoSort() bool {
	return len(p.Sort) <= 0
}

func (p *Pagination) GetLimit() int {
	return p.GetSize()
}

func (p *Pagination) GetOffset() (offset int) {
	if offset = (p.GetPage() - 1) * p.GetSize(); offset < 0 {
		offset = 0
	}
	return offset
}

func GetPageable(c *gin.Context) Pageable {
	// 1. From Query
	result := &Pagination{
		Page: GetPage(c),
		Size: GetSize(c),
		Sort: GetSort(c),
	}

	return result
}

func GetLimit(c *gin.Context) int {
	result := 0

	size := GetSize(c)
	page := GetPage(c)

	if page > 0 && size > 0 {
		result = (page - 1) * size
	}

	return result
}

func GetPage(c *gin.Context) int {
	result := 0

	page := c.GetInt("page")
	if page <= 0 {
		page, _ = strconv.Atoi(c.DefaultQuery("page", "0"))
		if page <= 0 {
			var p Pagination
			if err := form.ShouldBind(c, &p); err == nil {
				return p.Page
			}
		}
	}
	if page > 0 {
		result = page
	}

	return result
}

func GetSize(c *gin.Context) int {
	result := 20

	size := c.GetInt("size")
	if size <= 0 {
		size, _ = strconv.Atoi(c.DefaultQuery("size", "0"))
		if size <= 0 {
			var p Pagination
			if err := form.ShouldBind(c, &p); err == nil && p.Size > 0 {
				return p.Size
			}
		}
	}

	if size > 0 {
		result = size
	}

	return result
}

func GetSort(c *gin.Context) []Sort {

	sorts := c.GetStringSlice("sort")
	if len(sorts) <= 0 {
		sorts := c.QueryArray("sort")
		if len(sorts) <= 0 {
			var p Pagination
			if err := form.ShouldBind(c, &p); err == nil {
				return p.Sort
			}
		}
	}

	result := []Sort{}
	if len(sorts) > 0 {
		for _, v := range sorts {
			seg := strings.Split(v, ",")
			sort := &Sort{
				By:        seg[0],
				Direction: "ASC",
			}

			if len(seg) > 1 {
				sort.Direction = seg[1]
			}

			result = append(result, *sort)
		}
	}

	return result
}

func Page(tx *gorm.DB, pageable Pageable) *gorm.DB {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			tx.Statement.Order(sort.String())
		}
	}

	// Page/Size
	tx.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())

	return tx
}

func PageWithDefaultSort(tx *gorm.DB, pageable Pageable, defaultSort ...Sort) *gorm.DB {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			tx.Statement.Order(sort.String())
		}
	} else {
		for _, sort := range defaultSort {
			tx.Statement.Order(sort.String())
		}
	}

	// Page/Size
	tx.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())

	return tx
}

func PageWithDefaultOrder(tx *gorm.DB, pageable Pageable, defaultOrder ...string) *gorm.DB {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			tx.Statement.Order(sort.String())
		}
	} else {
		for _, order := range defaultOrder {
			tx.Statement.Order(order)
		}
	}

	// Page/Size
	tx.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())
	return tx
}
