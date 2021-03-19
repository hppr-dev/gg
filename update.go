package gg

import (
  "github.com/gin-gonic/gin"
  "hppr.dev/gg/mdl"
)

func UpdateByID(urlParam string) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    dataMap := make(DefaultMap)
    updated := make(DefaultMap)
    model := GetModel(ctx)
    schema := GetModelSchema(ctx)
    ctx.BindJSON(&dataMap)
    pKeyColumn := schema.PrioritizedPrimaryField.DBName
    if err := mdl.MatchAnyMapToModel(dataMap, schema); err != nil {
      DefaultOutput(ctx, 500, gin.H{"error": err.Error()})
      return
    }
    db := GetDatabase(ctx)
    result := db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).Updates(dataMap)
    if result.Error != nil {
      DefaultOutput(ctx, 500, gin.H{"error": result.Error.Error()})
      return
    }
    findResult := db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).Take(&updated)
    if findResult.Error != nil {
      DefaultOutput(ctx, 500, gin.H{"error": findResult.Error.Error()})
      return
    }
    DefaultOutput(ctx, 200, updated)
  }
}
