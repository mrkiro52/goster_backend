package gorm

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Context struct {
	context.Context

	originalRequest bool
}

func (c *Context) Commit() error {
	if !c.originalRequest {
		return nil
	}

	if tx, err := GetTransaction(c); err == nil {
		return tx.Commit().Error
	} else {
		return err
	}
}

func (c *Context) Rollback() error {
	if !c.originalRequest {
		return nil
	}

	if tx, err := GetTransaction(c); err == nil {
		return tx.Rollback().Error
	} else {
		return err
	}
}

func withTransaction(ctx context.Context, db *gorm.DB) *Context {
	if _, err := GetTransaction(ctx); err == nil {
		if v, ok := ctx.(*Context); ok {
			return &Context{Context: v.Context, originalRequest: false}
		} else {
			// possible when someone's else code (e.g. third-party library) overwrites context.
			return &Context{Context: ctx, originalRequest: false}
		}
	}

	return &Context{
		Context:         context.WithValue(ctx, "gorm_tx", db.Begin().WithContext(ctx)),
		originalRequest: true,
	}
}

func GetTransaction(ctx context.Context) (*gorm.DB, error) {
	if v := ctx.Value("gorm_tx"); v != nil {
		return v.(*gorm.DB), nil
	} else {
		return nil, fmt.Errorf("can't extract transaction from context")
	}
}

type TransactionContextFactory func(context.Context) *Context

func provideTransactionContext(db *gorm.DB) TransactionContextFactory {
	return func(c context.Context) *Context {
		return withTransaction(c, db)
	}
}
