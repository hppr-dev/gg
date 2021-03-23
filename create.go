package gg

import (
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "hppr.dev/gg/mdl"
)

// BodyCreate returns a handler that creates a model record based on post data
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
    ctx.ShouldBindBodyWith(&modelRef, binding.JSON)
    db := GetDatabase(ctx)
    db.Model(model).Create(modelRef)
    DefaultOutput(ctx, 201, convertStructToOutMap(modelRef, schema))
  }
}
