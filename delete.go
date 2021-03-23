package gg

import (
  "github.com/gin-gonic/gin"
)

func DeleteByID(url_param string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    db := GetDatabase(ctx)
    model := GetModel(ctx)
    result := db.Model(model).Delete(model, ctx.Param(url_param))
    if result.RowsAffected == 0{
      DefaultOutput(ctx, 404, gin.H{"error": "not found"})
      return
    }
    DefaultOutput(ctx, 200, gin.H{"message": "Record deleted"})
  }
}
