package db

// This file contains a mock of the KV interface that can fail on
// demand in order to make the tests around the db implementation
// simpler.

import (
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	types "github.com/netauth/protocol"
)

var (
	goodEntityBytes1, _ = proto.Marshal(&types.Entity{ID: proto.String("entity1"), Number: proto.Int32(1)})
	goodEntityBytes2, _ = proto.Marshal(&types.Entity{ID: proto.String("entity2"), Number: proto.Int32(7)})
	goodGroupBytes1, _  = proto.Marshal(&types.Group{Name: proto.String("group1"), Number: proto.Int32(1)})
	goodGroupBytes2, _  = proto.Marshal(&types.Group{Name: proto.String("group1"), Number: proto.Int32(7)})
)

type mockKV struct {
	mock.Mock
}

func newMockKV(hclog.Logger) (KVStore, error) {
	x := &mockKV{}
	x.On("SetEventFunc", mock.Anything).Return()
	return x, nil
}

func newMockKVError(hclog.Logger) (KVStore, error) {
	return nil, errors.New("Initialization error")
}

func (mkv *mockKV) Put(k string, v []byte) error {
	args := mkv.Called(k, v)
	return args.Error(0)
}

func (mkv *mockKV) Get(k string) ([]byte, error) {
	args := mkv.Called(k)
	return args.Get(0).([]byte), args.Error(1)
}

func (mkv *mockKV) Del(k string) error {
	return mkv.Called(k).Error(0)
}

func (mkv *mockKV) Keys(f string) ([]string, error) {
	args := mkv.Called(f)
	return args.Get(0).([]string), args.Error(1)
}

func (mkv *mockKV) Close() error {
	args := mkv.Called()
	return args.Error(0)
}

func (mkv *mockKV) Capabilities() []KVCapability {
	return mkv.Called().Get(0).([]KVCapability)
}

func (mkv *mockKV) SetEventFunc(f func(Event)) {
	mkv.Called(f)
}
