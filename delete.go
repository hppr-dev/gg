package gg

import (
  "github.com/gin-gonic/gin"
)

func DeleteByID(url_param string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    db := GetDatabase(ctx)
    model := GetModel(ctx)
    if err := db.Model(model).Delete(model, ctx.Param(url_param)).Error; err != nil {
      DefaultOutput(ctx, 404, gin.H{"error": err.Error()})
      return
    }
    DefaultOutput(ctx, 200, gin.H{"message": "Record deleted"})
  }
}
