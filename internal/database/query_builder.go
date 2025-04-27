package database

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	operation string // SELECT, INSERT, UPDATE, DELETE

	table     string   // TABLE
	columns   []string // COLUMNS
	numValues int      // NUMBER OF VALUES (e.g. VALUES (?, ?, ?))

	whereClauses []string
}

func NewQueryBuilder(operation string) *QueryBuilder {
	return &QueryBuilder{operation: operation}
}

func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.table = table
	return qb
}

func (qb *QueryBuilder) Columns(columns ...string) *QueryBuilder {
	qb.columns = columns
	return qb
}

func (qb *QueryBuilder) Values(numValues int) *QueryBuilder {
	qb.numValues = numValues
	return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.whereClauses = append(qb.whereClauses, condition)
	return qb
}

func (qb *QueryBuilder) Build() string {
	var query string

	switch qb.operation {
	case "SELECT":
		cols := "*"
		if len(qb.columns) > 0 {
			cols = strings.Join(qb.columns, ", ")
		}
		query = fmt.Sprintf("SELECT %s FROM %s", cols, qb.table)
		if len(qb.whereClauses) > 0 {
			query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
		}

	case "INSERT":
		cols := strings.Join(qb.columns, ", ")
		placeholders := strings.Repeat("?, ", qb.numValues)
		placeholders = strings.TrimSuffix(placeholders, ", ")
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", qb.table, cols, placeholders)

	case "UPDATE":
		setClauses := []string{}
		for _, col := range qb.columns {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
		}
		setString := strings.Join(setClauses, ", ")
		query = fmt.Sprintf("UPDATE %s SET %s", qb.table, setString)
		if len(qb.whereClauses) > 0 {
			query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
		}

	case "DELETE":
		query = fmt.Sprintf("DELETE FROM %s", qb.table)
		if len(qb.whereClauses) > 0 {
			query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
		}
	}

	return query
}
