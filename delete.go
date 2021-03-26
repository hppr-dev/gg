package gg

import (
	"github.com/gin-gonic/gin"
)

// DeleteByID returns a handler that will delete a given record with the matching primary key given in the url
func DeleteByID(urlParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := GetDatabase(ctx)
		model := GetModel(ctx)
		result := db.Model(model).Delete(model, ctx.Param(urlParam))
		if result.RowsAffected == 0 {
			DefaultOutput(ctx, 404, gin.H{"error": "not found"})
			return
		}
		DefaultOutput(ctx, 200, gin.H{"message": "Record deleted"})
	}
}
