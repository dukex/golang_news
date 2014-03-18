package gorm

import (
	"regexp"
	"strconv"
)

type search struct {
	db           *DB
	whereClause  []map[string]interface{}
	orClause     []map[string]interface{}
	notClause    []map[string]interface{}
	initAttrs    []interface{}
	assignAttrs  []interface{}
	havingClause map[string]interface{}
	orders       []string
	joinsStr     string
	selectStr    string
	offsetStr    string
	limitStr     string
	groupStr     string
	tableName    string
	unscope      bool
	raw          bool
}

func (s *search) clone() *search {
	return &search{
		whereClause:  s.whereClause,
		orClause:     s.orClause,
		notClause:    s.notClause,
		initAttrs:    s.initAttrs,
		assignAttrs:  s.assignAttrs,
		havingClause: s.havingClause,
		orders:       s.orders,
		selectStr:    s.selectStr,
		offsetStr:    s.offsetStr,
		limitStr:     s.limitStr,
		unscope:      s.unscope,
		groupStr:     s.groupStr,
		joinsStr:     s.joinsStr,
		tableName:    s.tableName,
		raw:          s.raw,
	}
}

func (s *search) where(query interface{}, values ...interface{}) *search {
	s.whereClause = append(s.whereClause, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *search) not(query interface{}, values ...interface{}) *search {
	s.notClause = append(s.notClause, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *search) or(query interface{}, values ...interface{}) *search {
	s.orClause = append(s.orClause, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *search) attrs(attrs ...interface{}) *search {
	s.initAttrs = append(s.initAttrs, toSearchableMap(attrs...))
	return s
}

func (s *search) assign(attrs ...interface{}) *search {
	s.assignAttrs = append(s.assignAttrs, toSearchableMap(attrs...))
	return s
}

func (s *search) order(value string, reorder ...bool) *search {
	if len(reorder) > 0 && reorder[0] {
		s.orders = []string{value}
	} else {
		s.orders = append(s.orders, value)
	}
	return s
}

func (s *search) selects(value interface{}) *search {
	s.selectStr = s.getInterfaceAsSql(value)
	return s
}

func (s *search) limit(value interface{}) *search {
	s.limitStr = s.getInterfaceAsSql(value)
	return s
}

func (s *search) offset(value interface{}) *search {
	s.offsetStr = s.getInterfaceAsSql(value)
	return s
}

func (s *search) group(query string) *search {
	s.groupStr = s.getInterfaceAsSql(query)
	return s
}

func (s *search) having(query string, values ...interface{}) *search {
	s.havingClause = map[string]interface{}{"query": query, "args": values}
	return s
}

func (s *search) includes(value interface{}) *search {
	return s
}

func (s *search) joins(query string) *search {
	s.joinsStr = query
	return s
}

func (s *search) setraw(b bool) *search {
	s.raw = b
	return s
}

func (s *search) unscoped() *search {
	s.unscope = true
	return s
}

func (s *search) table(name string) *search {
	s.tableName = name
	return s
}

func (s *search) getInterfaceAsSql(value interface{}) (str string) {
	switch value := value.(type) {
	case string:
		str = value
	case int:
		if value < 0 {
			str = ""
		} else {
			str = strconv.Itoa(value)
		}
	default:
		s.db.err(InvalidSql)
	}

	if !regexp.MustCompile("^\\s*[\\w\\s,.*()]*\\s*$").MatchString(str) {
		s.db.err(InvalidSql)
	}
	return
}
