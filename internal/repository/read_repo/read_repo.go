package read_repo

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

type ReadRepository struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewReadRepository(db *sql.DB) *ReadRepository {
	return &ReadRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
