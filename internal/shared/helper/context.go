package helper

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type txKey struct{}

func TxMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		if tx.Error != nil {
			c.AbortWithError(500, tx.Error)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		ctx := context.WithValue(c.Request.Context(), txKey{}, tx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			tx.Rollback()
		} else {
			tx.Commit()
		}

	}
}

func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func WithTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db.WithContext(ctx)
}
