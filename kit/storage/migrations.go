package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gitlab.medzdrav.ru/prototype/kit/log"
)

type Migration interface {
	Up() error
}

type migImpl struct {
	db     *sql.DB
	source string
	logger log.CLoggerFunc
}

func NewMigration(db *sql.DB, source string, logger log.CLoggerFunc) Migration {
	return &migImpl{
		db:     db,
		source: source,
		logger: logger,
	}
}

func (m *migImpl) Up() error {
	l := m.logger().Cmp("db-migration").Mth("up").InfF("applying from %s ...", m.source)
	err := goose.Up(m.db, m.source)
	if err != nil {
		l.E(err).St().Err("applying migrations")
		return err
	}
	version, err := goose.GetDBVersion(m.db)
	if err != nil {
		l.E(err).St().Err("getting version")
		return err
	}
	l.InfF("ok, version: %d", version)
	return nil
}
