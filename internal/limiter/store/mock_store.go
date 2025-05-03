package store

import (
	"fmt"
	"sync"
	"time"
)

// MockStore implements the Store interface for testing purposes
type MockStore struct {
	mu            sync.Mutex
	requestCounts map[string]int64
	blockedKeys   map[string]time.Time
	failOnCall    bool
}

// NewMockStore creates a new MockStore
func NewMockStore(failOnCall bool) *MockStore {
	return &MockStore{
		requestCounts: make(map[string]int64),
		blockedKeys:   make(map[string]time.Time),
		failOnCall:    failOnCall,
	}
}

// Increment increments the counter for a key and returns whether the request is allowed
// Increment in MockStore
func (m *MockStore) Increment(key string, limit int, blockSeconds int) (bool, time.Duration, error) {
	if m.failOnCall {
		return false, 0, ErrStorageFailure
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Verificar se a chave está bloqueada e remover bloqueios expirados
	for key, blockedUntil := range m.blockedKeys {
		if time.Now().After(blockedUntil) {
			delete(m.blockedKeys, key) // Remove a chave bloqueada se expirou
		}
	}

	// Verificar se chave está bloqueada
	if blockedUntil, ok := m.blockedKeys[key]; ok {
		if time.Now().Before(blockedUntil) {
			ttl := time.Until(blockedUntil)
			return false, ttl, nil
		}
		// Se bloqueio expirou, remover da lista de bloqueados
		delete(m.blockedKeys, key)
	}

	// Incrementar o contador
	m.requestCounts[key]++
	count := m.requestCounts[key]

	// Se o limite for excedido, bloquear a chave
	if count > int64(limit) {
		blockedUntil := time.Now().Add(time.Duration(blockSeconds) * time.Second)
		m.blockedKeys[key] = blockedUntil
		ttl := time.Until(blockedUntil)
		return false, ttl, nil
	}

	return true, 0, nil
}

// IsBlocked checks if a key is blocked
func (m *MockStore) IsBlocked(key string) (bool, time.Duration, error) {
	if m.failOnCall {
		return false, 0, ErrStorageFailure
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Verificar e remover bloqueios expirados
	if blockedUntil, ok := m.blockedKeys[key]; ok {
		if time.Now().After(blockedUntil) {
			// Bloqueio expirado, removendo a chave
			delete(m.blockedKeys, key)
			return false, 0, nil
		}
		// Caso contrário, ainda está bloqueado
		ttl := time.Until(blockedUntil)
		return true, ttl, nil
	}

	// Caso não esteja bloqueado
	return false, 0, nil
}

// GetRequestCount returns the current request count for a key
func (m *MockStore) GetRequestCount(key string) (int64, error) {
	if m.failOnCall {
		return 0, ErrStorageFailure
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	return m.requestCounts[key], nil
}

// ResetCounter reseta o contador de requisições para uma chave
func (m *MockStore) ResetCounter(key string) error {
	if m.failOnCall {
		return ErrStorageFailure
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.requestCounts, key)
	return nil
}

// ErrStorageFailure is returned when the storage operation fails
var ErrStorageFailure = fmt.Errorf("storage operation failed")
