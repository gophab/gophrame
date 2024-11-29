package query

import "gorm.io/gorm"

type Query struct {
	*gorm.DB
	Pagination
}

func (p *Query) Page(pageable Pageable) *Query {
	result := p.Infinite(pageable)
	if result.Error != nil {
		return result
	}

	result.Count(&result.Pagination.Total)
	return result
}

func (p *Query) PageWithDefaultOrders(pageable Pageable, defaultOrders ...string) *Query {
	result := p.InfiniteWithDefaultOrder(pageable, defaultOrders...)
	if result.Error != nil {
		return p
	}

	result.Count(&result.Pagination.Total)
	return result
}

func (p *Query) PageWithDefaultSorts(pageable Pageable, defaultSorts ...Sort) *Query {
	result := p.InfiniteWithDefaultSort(pageable, defaultSorts...)
	if result.Error != nil {
		return p
	}

	result.Count(&result.Pagination.Total)
	return p
}

func (p *Query) Infinite(pageable Pageable) *Query {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			p.Statement.Order(sort.String())
		}
	}

	// Page/Size
	p.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())

	return &Query{
		DB: p.DB,
		Pagination: Pagination{
			Total: -1,
			Page:  pageable.GetPage(),
			Size:  pageable.GetSize(),
			Sort:  pageable.GetSort(),
		},
	}
}

func (p *Query) InfiniteWithDefaultSort(pageable Pageable, defaultSorts ...Sort) *Query {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			p.Statement.Order(sort.String())
		}
	} else if len(defaultSorts) > 0 {
		for _, sort := range defaultSorts {
			p.Statement.Order(sort.String())
		}
	}

	// Page/Size
	p.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())

	return &Query{
		DB: p.DB,
		Pagination: Pagination{
			Total: -1,
			Page:  pageable.GetPage(),
			Size:  pageable.GetSize(),
			Sort:  pageable.GetSort(),
		},
	}
}

func (p *Query) InfiniteWithDefaultOrder(pageable Pageable, defaultOrders ...string) *Query {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			p.Statement.Order(sort.String())
		}
	} else if len(defaultOrders) > 0 {
		for _, order := range defaultOrders {
			p.Statement.Order(order)
		}
	}

	// Page/Size
	p.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())
	return &Query{
		DB: p.DB,
		Pagination: Pagination{
			Total: -1,
			Page:  pageable.GetPage(),
			Size:  pageable.GetSize(),
			Sort:  pageable.GetSort(),
		},
	}
}

func Page(tx *gorm.DB, pageable Pageable) *gorm.DB {
	if pageable != nil {
		// Order
		if !pageable.NoSort() {
			for _, sort := range pageable.GetSort() {
				tx.Statement.Order(sort.String())
			}
		}

		// Page/Size
		tx.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())
	}
	return tx
}

func PageWithDefaultOrder(tx *gorm.DB, pageable Pageable, defaultOrders ...string) *gorm.DB {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			tx.Statement.Order(sort.String())
		}
	} else if len(defaultOrders) > 0 {
		for _, order := range defaultOrders {
			tx.Statement.Order(order)
		}
	}

	// Page/Size
	tx.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())

	return tx
}

func PageWithDefaultSort(tx *gorm.DB, pageable Pageable, defaultSorts ...Sort) *gorm.DB {
	// Order
	if !pageable.NoSort() {
		for _, sort := range pageable.GetSort() {
			tx.Statement.Order(sort.String())
		}
	} else if len(defaultSorts) > 0 {
		for _, sort := range defaultSorts {
			tx.Statement.Order(sort.String())
		}
	}

	// Page/Size
	tx.Statement.Offset(pageable.GetOffset()).Limit(pageable.GetLimit())

	return tx
}
