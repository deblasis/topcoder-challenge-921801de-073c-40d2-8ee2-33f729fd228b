// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE. 
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
