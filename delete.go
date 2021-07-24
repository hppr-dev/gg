package gg

import (
	"github.com/gin-gonic/gin"
)

// DeleteByID returns a handler that will delete a given record with the matching primary key given in the url
func DeleteByID(urlParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := GetDatabase(ctx)
		model := GetModel(ctx)
    if cast, ok := model.(BeforeDeleteWithContexter) ; ok {
      if err := cast.BeforeDeleteWithContext(ctx, db) ; err != nil {
        return
      }
    }
		result := db.Model(model).Delete(model, ctx.Param(urlParam))
		if result.RowsAffected == 0 {
			DefaultOutput(ctx, 404, gin.H{"error": "not found"})
			return
		}
    if cast, ok := model.(AfterDeleteWithContexter) ; ok {
      if err := cast.AfterDeleteWithContext(ctx, db) ; err != nil {
        return
      }
    }
		DefaultOutput(ctx, 200, gin.H{"message": "Record deleted"})
	}
}
