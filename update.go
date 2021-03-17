package gg

import (
  "github.com/gin-gonic/gin"
  "hppr.dev/gg/mdl"
)

func UpdateByID(url_param string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    dataMap := make(DefaultMap)
    model := GetModel(ctx)
    schema := GetModelSchema(ctx)
    ctx.BindJSON(&dataMap)
    if err := mdl.MatchAnyMapToModel(dataMap, schema); err != nil {
      DefaultOutput(ctx, 500, gin.H{"error": err.Error()})
      return
    }
    db := GetDatabase(ctx)
    result := db.Model(model).Where("id = ?", ctx.Param(url_param)).Updates(dataMap)
    if result.Error != nil {
      DefaultOutput(ctx, 500, gin.H{"error": result.Error.Error()})
      return
    }
    DefaultOutput(ctx, 200, gin.H{"message": dataMap})
  }
}
