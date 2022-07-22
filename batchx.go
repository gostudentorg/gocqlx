package gocqlx

import (
	"github.com/gocql/gocql"
)

type Batch struct {
	*gocql.Batch
}

func (s *Session) NewBatch(typ gocql.BatchType) *Batch {
	return &Batch{
		Batch: s.Session.NewBatch(typ),
	}
}

func (b *Batch) BindStruct(qry Queryx, arg interface{}) error {
	arglist, err := qry.bindStructArgs(arg, nil)
	if err != nil {
		return err
	}
	b.Query(qry.Statement(), arglist...)
	return nil
}
