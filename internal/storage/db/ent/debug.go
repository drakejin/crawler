// Code generated by ent, DO NOT EDIT.

package ent

import "entgo.io/ent/dialect"

func (c *PageClient) Debug() *PageClient {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	return &PageClient{config: cfg}
}

func (c *PageReferredClient) Debug() *PageReferredClient {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	return &PageReferredClient{config: cfg}
}

func (c *PageSourceClient) Debug() *PageSourceClient {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	return &PageSourceClient{config: cfg}
}
