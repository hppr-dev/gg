package gg

import (
	"github.com/gin-gonic/gin"
	"hppr.dev/gg/mdl"
)

// BodyCreate returns a handler that creates a model record based on post data
func BodyCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
    gormData, err := marshalInputToGormData(ctx)
    if err != nil {
      DefaultOutput(ctx, 400, gin.H{"message" : "Error validating input"})
      return
    }
		db := GetDatabase(ctx)
    if cast, ok := gormData.(BeforeCreateWithContexter) ; ok {
      if err := cast.BeforeCreateWithContext(ctx, db) ; err != nil {
        return
      }
    }
		db.Create(gormData)
    if cast, ok := gormData.(AfterCreateWithContexter) ; ok {
      if err := cast.AfterCreateWithContext(ctx, db) ; err != nil {
        return
      }
    }
    outputGormData(ctx, gormData)
  }
}

func marshalInputToGormData(ctx *gin.Context) (interface{}, error) {
  marshal := GetMarshalType(ctx)
  jsonInputData := mdl.New(marshal)
	model := GetModel(ctx)
	gormData := mdl.New(model)
	err := ctx.ShouldBind(jsonInputData)
  mdl.CopyFields(jsonInputData, gormData)
  return gormData, err
}

func outputGormData(ctx *gin.Context, gormData interface{}) {
  marshal := GetMarshalType(ctx)
  jsonOutputData := mdl.New(marshal)
  mdl.CopyFields(gormData, jsonOutputData)
	DefaultOutput(ctx, 201, jsonOutputData)
}
