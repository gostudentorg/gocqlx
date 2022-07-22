package gocqlx

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/maraino/go-mock"
	"github.com/stretchr/testify/assert"

	"testing"
)

const (
	LoggedBatch gocql.BatchType = 0
)

type MockQB struct {
	names []string
}

type QB interface {
	//qb.Builder

	Query(session Session) *Queryx
	QueryContext(ctx context.Context, session Session) *Queryx
}

func (mqb *MockQB) Query(session Session) *Queryx {
	return &Queryx{
		Names: []string{"Name", "Age", "First", "Last"},
	}
}

func (mqb *MockQB) QueryContext(ctx context.Context, session Session) *Queryx {
	return &Queryx{}
}

type user struct {
	insertQry MockQB
}

type SessionMock struct {
	mock.Mock
}

func (s *SessionMock) NewBatch(typ gocql.BatchType) *Batch {
	return &Batch{}
}

func (s *SessionMock) ExecuteBatch(batch *gocql.Batch) error {
	return nil
}

func TestBatch_BindStruct(t *testing.T) {
	v := &struct {
		Name  string
		Age   int
		First string
		Last  string
	}{
		Name:  "name",
		Age:   30,
		First: "first",
		Last:  "last",
	}

	t.Run("simple", func(t *testing.T) {
		var sessionMock = &SessionMock{}
		var session Session

		batchWrapper := sessionMock.NewBatch(LoggedBatch)
		usertmp := new(user)
		//user.insertQry.QueryContext(context.Background(), session).BindStruct(v).ExecRelease()

		err := batchWrapper.BindStruct(*usertmp.insertQry.Query(session), v)

		err = sessionMock.ExecuteBatch(batchWrapper.Batch)
		assert.NoError(t, err)
	})
}
