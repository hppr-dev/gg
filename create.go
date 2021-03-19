package gg

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "hppr.dev/gg/mdl"
)

func BodyCreate() gin.HandlerFunc {
  return func(ctx *gin.Context) {
    schema := GetModelSchema(ctx)
    model := GetModel(ctx)
    dataMap := make(DefaultMap)
    ctx.ShouldBindBodyWith(&dataMap, binding.JSON)
    if err := mdl.MatchAllMapToModel(dataMap, schema); err != nil {
      DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
      return
    }
    modelRef := mdl.NewModel(model)
    if err := ctx.ShouldBindBodyWith(&modelRef, binding.JSON); err != nil {
      DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
      return
    }
    db := GetDatabase(ctx)
    if err := db.Model(model).Create(modelRef).Error; err != nil {
      DefaultOutput(ctx, 500, gin.H{"error": err.Error()})
    }
    DefaultOutput(ctx, 201, mdl.ExtractModelIDs(schema, modelRef))
  }
}
