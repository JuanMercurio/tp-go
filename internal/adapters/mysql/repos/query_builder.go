package repos

import (
	"fmt"
	"strings"
)

type queryBuilder struct {
	Select  string
	From    string
	where   []string // los OR deberan ser manuales por ahora
	OrderBy []string
	limit   string
	offset  string
	joins   []string
}

func (s queryBuilder) toString() string {
	return s.Select + " " +
		s.From + " " +
		s.joinString() + " " +
		s.whereString() + " " +
		s.orderByString() + " " +
		s.limit + " " +
		s.offset
}

func (s *queryBuilder) AddWhere(condicion string) {
	s.where = append(s.where, condicion)
}

func (s *queryBuilder) AddLimit(limite int) {
	s.limit = fmt.Sprintf("LIMIT %d", limite)
}

func (s *queryBuilder) AddOffset(offset int) {
	s.offset = fmt.Sprintf("OFFSET %d", offset)
}

func (s *queryBuilder) AddOrderBy(orden string) {
	s.OrderBy = append(s.OrderBy, orden)
}

func (s queryBuilder) orderByString() string {
	if len(s.OrderBy) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(s.OrderBy, ",")
}

func (s queryBuilder) whereString() string {
	if len(s.where) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(s.where, " AND ")
}

func (s *queryBuilder) AddJoin(tabla, condicion string) {
	s.joins = append(s.joins, tabla+" ON "+condicion)
}

func (s queryBuilder) joinString() string {
	if len(s.joins) == 0 {
		return ""
	}
	return "JOIN " + strings.Join(s.joins, " JOIN ")
}
