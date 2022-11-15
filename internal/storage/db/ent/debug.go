// Code generated by ent, DO NOT EDIT.

package ent

import "entgo.io/ent/dialect"

func (c *PageInfoClient) Debug() *PageInfoClient {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	return &PageInfoClient{config: cfg}
}

func (c *PageLinkClient) Debug() *PageLinkClient {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	return &PageLinkClient{config: cfg}
}