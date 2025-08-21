package pg

import "backend-poc/backoffice/extensions/pg"

type Repository struct {
	*pg.Pool
}

func NewRepository(pool *pg.Pool) *Repository {
	return &Repository{
		Pool: pool,
	}
}
