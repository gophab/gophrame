package server

import (
	"context"
	"errors"
	"sync"

	"github.com/wjshen/gophrame/core/database"
	"github.com/wjshen/gophrame/core/inject"

	"github.com/go-oauth2/oauth2/v4"
)

var (
	theClientStore oauth2.ClientStore
)

func ClientStore() oauth2.ClientStore {
	if theClientStore == nil {
		theClientStore = InitClientStore()
	}
	return theClientStore
}

func InitClientStore() oauth2.ClientStore {
	if theClientStore == nil {
		if result, err := NewDatabaseClientStore(); err == nil {
			inject.InjectValue("clientStore", result)
			theClientStore = result
		}
	}

	return theClientStore
}

// NewClientStore create client store
/**
 * Database Client Store
 */
func NewDatabaseClientStore() (oauth2.ClientStore, error) {
	return &DatabaseClientStore{
		data: make(map[string]oauth2.ClientInfo),
	}, nil
}

// ClientStore client information store
type DatabaseClientStore struct {
	sync.RWMutex
	data map[string]oauth2.ClientInfo
}

// GetByID according to the ID for the client information
func (cs *DatabaseClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	cs.RLock()
	defer cs.RUnlock()

	if c, ok := cs.data[id]; ok {
		return c, nil
	}

	var client OAuthClient

	result := database.DB().Where("client_id = ? AND del_flag = false ", id).First(&client)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected > 0 {
		cs.data[id] = &client
		return &client, nil
	}

	return nil, errors.New("not found")
}
