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

type MockQB struct{}

type QB interface {
	Query(session Session) *Queryx
	QueryContext(ctx context.Context, session Session) *Queryx
}

func (mqb *MockQB) Query(session SessionMock) *Queryx {
	return &Queryx{
		Names:  []string{"Name", "Age", "First", "Last"},
		Mapper: DefaultMapper,
	}
}

func (mqb *MockQB) QueryContext(_ context.Context, _ SessionMock) *Queryx {
	return &Queryx{}
}

type User struct {
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

	v2 := &struct {
		User  string
		Age   int
		Email string
		Phone string
	}{
		User:  "Daniel",
		Age:   30,
		Email: "daniel@yahoo.com",
		Phone: "12345678",
	}

	t.Run("simple insert", func(t *testing.T) {
		var sessionMock = &SessionMock{}

		batchWrapper := sessionMock.NewBatch(LoggedBatch)
		usertmp := new(User)

		err := batchWrapper.BindStruct(*usertmp.insertQry.Query(*sessionMock), v)

		err = sessionMock.ExecuteBatch(batchWrapper.Batch)
		assert.NoError(t, err)
	})

	t.Run("multiple inserts", func(t *testing.T) {
		var sessionMock = &SessionMock{}

		batchWrapper := sessionMock.NewBatch(LoggedBatch)
		usertmp := new(User)

		err := batchWrapper.BindStruct(*usertmp.insertQry.Query(*sessionMock), v)

		err = batchWrapper.BindStruct(*usertmp.insertQry.Query(*sessionMock), v2)

		err = sessionMock.ExecuteBatch(batchWrapper.Batch)
		assert.NoError(t, err)
	})
}
