//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type TokensCache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
}

type record struct {
	value     interface{}
	timestamp time.Time
	ttl       time.Duration
}

type MemTokensCache struct {
	mutex  sync.RWMutex
	tokens map[string]*record
}

func NewMemTokensCache() *MemTokensCache {
	c := &MemTokensCache{
		tokens: make(map[string]*record),
	}

	go func() {
		for now := range time.Tick(time.Second) {
			c.mutex.Lock()
			for k, v := range c.tokens {
				if int64(now.Sub(v.timestamp).Seconds()) > int64(v.ttl.Seconds()) {
					delete(c.tokens, k)
				}
			}
			c.mutex.Unlock()
		}
	}()
	return c
}

func (c *MemTokensCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.tokens[key] != nil {
		return errors.New("already existing")
	}

	c.tokens[key] = &record{
		value:     value,
		timestamp: time.Now().UTC(),
		ttl:       ttl,
	}
	return nil
}

func (c *MemTokensCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	record := c.tokens[key]
	if record == nil {
		return nil, nil
	}

	return record.value, nil
}
