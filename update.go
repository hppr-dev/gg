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
      DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
      return
    }
    db := GetDatabase(ctx)
    result := db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).Updates(dataMap)
    if result.RowsAffected == 0 {
      DefaultOutput(ctx, 404, gin.H{"error": "not found"})
      return
    }
    db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).First(updated)
    DefaultOutput(ctx, 200, updated)
  }
}
