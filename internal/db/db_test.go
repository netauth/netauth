package db

import (
	"errors"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	types "github.com/netauth/protocol"
)

func TestNew(t *testing.T) {
	RegisterKV("mock", newMockKV)
	_, err := New("mock")
	assert.Nil(t, err)

	RegisterKV("errorKV", newMockKVError)
	res, err := New("errorKV")
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestPrimeIndexes(t *testing.T) {
	RegisterKV("mock", newMockKV)

	// No error, no keys returned
	m, err := New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, nil)
	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{}, nil)
	err = m.PrimeIndexes()
	assert.Nil(t, err)

	// Error, bad entity
	m, err = New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{"/entities/bad-proto"}, nil)
	m.kv.(*mockKV).On("Get", "/entities/bad-proto").Return([]byte{42}, nil)
	err = m.PrimeIndexes()
	assert.NotNil(t, err)

	// Error, bad group
	m, err = New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, nil)
	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{"/groups/bad-proto"}, nil)
	m.kv.(*mockKV).On("Get", "/groups/bad-proto").Return([]byte{42}, nil)
	err = m.PrimeIndexes()
	assert.NotNil(t, err)

}

func TestLoadEntityIndex(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, ErrInternalError)
	err = m.loadEntityIndex()
	assert.NotNil(t, err)

	// Works fine
	m, err = New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{"/entities/good"}, nil)
	m.kv.(*mockKV).On("Get", "/entities/good").Return(goodEntityBytes1, nil)
	err = m.loadEntityIndex()
	assert.Nil(t, err)

	// Fails load
	m, err = New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{"/entities/bad-proto"}, nil)
	m.kv.(*mockKV).On("Get", "/entities/bad-proto").Return([]byte{42}, nil)
	err = m.loadEntityIndex()
	assert.NotNil(t, err)
}

func TestLoadGroupIndex(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{}, ErrInternalError)
	err = m.loadGroupIndex()
	assert.NotNil(t, err)

	// Works fine
	m, err = New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{"/groups/good"}, nil)
	m.kv.(*mockKV).On("Get", "/groups/good").Return(goodGroupBytes1, nil)
	err = m.loadGroupIndex()
	assert.Nil(t, err)

	// Fails load
	m, err = New("mock")
	assert.Nil(t, err)
	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{"/groups/bad-proto"}, nil)
	m.kv.(*mockKV).On("Get", "/groups/bad-proto").Return([]byte{42}, nil)
	err = m.loadGroupIndex()
	assert.NotNil(t, err)
}

func TestDiscoverEntityIDs(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	entList := []string{"/entities/foo", "/entities/bar"}
	m.kv.(*mockKV).On("Keys", "/entities/*").Return(entList, nil)

	res, err := m.DiscoverEntityIDs()
	assert.Nil(t, err)
	assert.Equal(t, res, entList)
}

func TestLoadEntity(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Get", "/entities/missing").Return([]byte{}, ErrNoValue)
	m.kv.(*mockKV).On("Get", "/entities/bad-error").Return([]byte{}, errors.New("something internal"))
	m.kv.(*mockKV).On("Get", "/entities/bad-proto").Return([]byte{42}, nil)
	m.kv.(*mockKV).On("Get", "/entities/good").Return(goodEntityBytes1, nil)

	cases := []struct {
		id      string
		wantErr error
	}{
		{"missing", ErrUnknownEntity},
		{"bad-error", ErrInternalError},
		{"bad-proto", ErrInternalError},
		{"good", nil},
	}

	for _, c := range cases {
		_, err := m.LoadEntity(c.id)
		assert.Equal(t, c.wantErr, err)
	}
}

func TestSaveEntity(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Put", "/entities/good", mock.Anything).Return(nil)
	m.kv.(*mockKV).On("Put", "/entities/bad", mock.Anything).Return(errors.New("something internal"))

	err = m.SaveEntity(&types.Entity{ID: proto.String("good")})
	assert.Nil(t, err)

	err = m.SaveEntity(nil)
	assert.NotNil(t, err)

	err = m.SaveEntity(&types.Entity{ID: proto.String("bad")})
	assert.NotNil(t, err)
}

func TestDeleteEntity(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Del", "/entities/good").Return(nil)
	m.kv.(*mockKV).On("Del", "/entities/missing").Return(ErrNoValue)

	assert.Nil(t, m.DeleteEntity("good"))
	assert.Equal(t, m.DeleteEntity("missing"), ErrUnknownEntity)
}

func TestDiscoverGroupNames(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	grpList := []string{"/groups/foo", "/groups/bar"}
	m.kv.(*mockKV).On("Keys", "/groups/*").Return(grpList, nil)

	res, err := m.DiscoverGroupNames()
	assert.Nil(t, err)
	assert.Equal(t, res, grpList)
}

func TestLoadGroup(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Get", "/groups/missing").Return([]byte{}, ErrNoValue)
	m.kv.(*mockKV).On("Get", "/groups/bad-error").Return([]byte{}, errors.New("something internal"))
	m.kv.(*mockKV).On("Get", "/groups/bad-proto").Return([]byte{42}, nil)
	m.kv.(*mockKV).On("Get", "/groups/good").Return(goodGroupBytes1, nil)

	cases := []struct {
		id      string
		wantErr error
	}{
		{"missing", ErrUnknownGroup},
		{"bad-error", errors.New("something internal")},
		{"bad-proto", errors.New("unexpected EOF")},
		{"good", nil},
	}

	for _, c := range cases {
		_, err := m.LoadGroup(c.id)
		assert.Equal(t, err, c.wantErr)
	}
}

func TestSaveGroup(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Put", "/groups/good", mock.Anything).Return(nil)
	m.kv.(*mockKV).On("Put", "/groups/bad", mock.Anything).Return(errors.New("something internal"))

	err = m.SaveGroup(&types.Group{Name: proto.String("good")})
	assert.Nil(t, err)

	err = m.SaveGroup(nil)
	assert.NotNil(t, err)

	err = m.SaveGroup(&types.Group{Name: proto.String("bad")})
	assert.NotNil(t, err)
}

func TestDeleteGroup(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Del", "/groups/good").Return(nil)
	m.kv.(*mockKV).On("Del", "/groups/missing").Return(ErrNoValue)

	assert.Nil(t, m.DeleteGroup("good"))
	assert.Equal(t, m.DeleteGroup("missing"), ErrUnknownGroup)
}

func TestShutdown(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Close").Return(errors.New("Error syncing KV"))
	m.Shutdown()
}

func TestNextEntityNumber(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Get", "/entities/load-error").Return([]byte{}, errors.New("KV Load error"))
	m.kv.(*mockKV).On("Get", "/entities/entity1").Return(goodEntityBytes1, nil)
	m.kv.(*mockKV).On("Get", "/entities/entity2").Return(goodEntityBytes2, nil)

	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, nil).Once()
	res, err := m.NextEntityNumber()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), res)

	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, errors.New("retrieval error")).Once()
	_, err = m.NextEntityNumber()
	assert.NotNil(t, err)

	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{}, errors.New("retrieval error")).Once()
	_, err = m.NextEntityNumber()
	assert.NotNil(t, err)

	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{"/entities/entity1", "/entities/load-error"}, nil).Once()
	_, err = m.NextEntityNumber()
	assert.NotNil(t, err)

	m.kv.(*mockKV).On("Keys", "/entities/*").Return([]string{"/entities/entity1", "/entities/entity2"}, nil).Once()
	res, err = m.NextEntityNumber()
	assert.Nil(t, err)
	assert.Equal(t, int32(8), res)
}

func TestNextGroupNumber(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Get", "/groups/load-error").Return([]byte{}, errors.New("KV Load error"))
	m.kv.(*mockKV).On("Get", "/groups/group1").Return(goodGroupBytes1, nil)
	m.kv.(*mockKV).On("Get", "/groups/group2").Return(goodGroupBytes2, nil)

	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{}, nil).Once()
	res, err := m.NextGroupNumber()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), res)

	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{}, errors.New("retrieval error")).Once()
	_, err = m.NextGroupNumber()
	assert.NotNil(t, err)

	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{}, errors.New("retrieval error")).Once()
	_, err = m.NextGroupNumber()
	assert.NotNil(t, err)

	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{"/groups/group1", "/groups/load-error"}, nil).Once()
	_, err = m.NextGroupNumber()
	assert.NotNil(t, err)

	m.kv.(*mockKV).On("Keys", "/groups/*").Return([]string{"/groups/group1", "/groups/group2"}, nil).Once()
	res, err = m.NextGroupNumber()
	assert.Nil(t, err)
	assert.Equal(t, int32(8), res)
}

func TestCapabilities(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Capabilities").Return([]KVCapability{})

	assert.Equal(t, []KVCapability{}, m.Capabilities())
}

func TestDBSearchEntities(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	res, err := m.SearchEntities(SearchRequest{})
	assert.Equal(t, ErrBadSearch, err)
	assert.Equal(t, []*types.Entity(nil), res)

	res, err = m.SearchEntities(SearchRequest{Expression: "*"})
	assert.Nil(t, err)
	assert.Equal(t, []*types.Entity{}, res)
}

func TestDBSearchGroups(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	res, err := m.SearchGroups(SearchRequest{})
	assert.Equal(t, ErrBadSearch, err)
	assert.Equal(t, []*types.Group(nil), res)

	res, err = m.SearchGroups(SearchRequest{Expression: "*"})
	assert.Nil(t, err)
	assert.Equal(t, []*types.Group{}, res)
}

func TestLoadEntityBatch(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Get", "/entities/load-error").Return([]byte{}, errors.New("KV Load error"))
	m.kv.(*mockKV).On("Get", "/entities/entity1").Return(goodEntityBytes1, nil)
	m.kv.(*mockKV).On("Get", "/entities/entity2").Return(goodEntityBytes2, nil)

	res, err := m.loadEntityBatch([]string{"entity1", "load-error"})
	assert.NotNil(t, err)
	assert.Nil(t, res)

	_, err = m.loadEntityBatch([]string{"entity1", "entity2"})
	assert.Nil(t, err)
}

func TestLoadGroupBatch(t *testing.T) {
	RegisterKV("mock", newMockKV)
	m, err := New("mock")
	assert.Nil(t, err)

	m.kv.(*mockKV).On("Get", "/groups/load-error").Return([]byte{}, errors.New("KV Load error"))
	m.kv.(*mockKV).On("Get", "/groups/group1").Return(goodGroupBytes1, nil)
	m.kv.(*mockKV).On("Get", "/groups/group2").Return(goodGroupBytes2, nil)

	res, err := m.loadGroupBatch([]string{"group1", "load-error"})
	assert.NotNil(t, err)
	assert.Nil(t, res)

	_, err = m.loadGroupBatch([]string{"group1", "group2"})
	assert.Nil(t, err)
}

func TestSetParentLogger(t *testing.T) {
	lb = nil
	assert.Nil(t, lb)
	SetParentLogger(hclog.NewNullLogger())
	assert.NotNil(t, lb)
}
