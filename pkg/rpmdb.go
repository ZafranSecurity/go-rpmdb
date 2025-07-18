package rpmdb

import (
	"github.com/ZafranSecurity/go-rpmdb/pkg/bdb"
	dbi "github.com/ZafranSecurity/go-rpmdb/pkg/db"
	"github.com/ZafranSecurity/go-rpmdb/pkg/ndb"
	"github.com/ZafranSecurity/go-rpmdb/pkg/sqlite3"
	"github.com/samber/lo"
	"golang.org/x/xerrors"
)

type RpmDB struct {
	Db dbi.RpmDBInterface
}

func Open(path string) (*RpmDB, error) {
	// SQLite3 Open() returns nil, nil in case of DB format other than SQLite3
	sqldb, err := sqlite3.Open(path)
	if err != nil && !xerrors.Is(err, sqlite3.ErrorInvalidSQLite3) {
		return nil, err
	}
	if sqldb != nil {
		return &RpmDB{Db: sqldb}, nil
	}

	// NDB Open() returns nil, nil in case of DB format other than NDB
	ndbh, err := ndb.Open(path)
	if err != nil && !xerrors.Is(err, ndb.ErrorInvalidNDB) {
		return nil, err
	}
	if ndbh != nil {
		return &RpmDB{Db: ndbh}, nil
	}

	odb, err := bdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &RpmDB{
		Db: odb,
	}, nil
}

func (d *RpmDB) Close() error {
	return d.Db.Close()
}

func (d *RpmDB) Package(name string) (*PackageInfo, error) {
	pkgs, err := d.ListPackages()
	if err != nil {
		return nil, xerrors.Errorf("unable to list packages: %w", err)
	}

	for _, pkg := range pkgs {
		if pkg.Name == name {
			return pkg, nil
		}
	}
	return nil, xerrors.Errorf("%s is not installed", name)
}

func (d *RpmDB) ListPackages() ([]*PackageInfo, error) {
	var pkgList []*PackageInfo

	for entry := range d.Db.Read() {
		if entry.Err != nil {
			return nil, entry.Err
		}

		indexEntries, err := headerImport(entry.Value)
		if err != nil {
			return nil, xerrors.Errorf("error during importing header: %w", err)
		}
		pkg, err := getNEVRA(indexEntries)
		if err != nil {
			return nil, xerrors.Errorf("invalid package info: %w", err)
		}

		pkg.BdbFirstOverflowPgNo = entry.BdbFirstOverflowPgNo
		pkg.RawHeader = entry.Value
		pkg.IndexEntries = lo.Map(indexEntries, func(x IndexEntry, _ int) IndexEntry {
			x.Data = nil
			return x
		})

		pkgList = append(pkgList, pkg)
	}

	return pkgList, nil
}
