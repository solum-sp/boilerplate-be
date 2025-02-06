package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// "proposal-template/model"
	"proposal-template/pkg/utils"

	sq "github.com/Masterminds/squirrel"
)

type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p *Paging) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
}


type GenericDAO[T any] struct {
	db *sql.DB
	tableName string
}

func NewGenericDAO[T any](db *sql.DB, tableName string) *GenericDAO[T] {
	return &GenericDAO[T]{
		db: db,
		tableName: tableName,
	}
}

func (dao *GenericDAO[T]) GetByColumn(ctx context.Context, column string, value interface{}) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 LIMIT 1", dao.tableName, column)
	row := dao.db.QueryRowContext(ctx, query, value)

	var obj T
	err := row.Scan(&obj)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no record is found
		}
		return nil, fmt.Errorf("error retrieving data: %w", err)
	}
	return &obj, nil
}

func (dao *GenericDAO[T]) List(ctx context.Context, paging Paging, query string, args ...interface{}) ([]T, error) {
	paging.Validate()

	offset := (paging.Page - 1) * paging.Limit
	queryWithPagination := fmt.Sprintf("%s ORDER BY id DESC LIMIT %d OFFSET %d", query, paging.Limit, offset)

	rows, err := dao.db.QueryContext(ctx, queryWithPagination, args...)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %w", err)
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var obj T
		if err := rows.Scan(&obj); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, obj)
	}

	return results, nil
}

func (dao *GenericDAO[T]) Create(ctx context.Context, model T) (int, error) {
	// Convert struct to map[string]interface{}, filtering out zero-values
	valuesMap, err := utils.StructToMap(model)
	if err != nil {
		return 0, fmt.Errorf("error converting model to map: %w", err)
	}
	valuesMap["CreatedAt"] = time.Now().UTC()
	valuesMap["UpdatedAt"] = time.Now().UTC()

	// Build the INSERT query using only non-zero fields
	query, args, err := sq.Insert(dao.tableName).
		SetMap(valuesMap).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building insert query: %w", err)
	}

	// Execute the query and scan the returning ID
	var id int
	err = dao.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting data: %w", err)
	}
	return id, nil
}

