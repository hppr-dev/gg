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
      DefaultOutput(ctx, 400, gin.H{"error" : err.Error()})
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
  inputData := mdl.New(marshal)
	model := GetModel(ctx)
	gormData := mdl.New(model)
	err := bindData(ctx, inputData)
  mdl.CopyFields(inputData, gormData)
  return gormData, err
}

func bindData(ctx *gin.Context, output interface{}) error {
  return ctx.ShouldBind(output)
}

func outputGormData(ctx *gin.Context, gormData interface{}) {
  outputGormDataStatus(ctx, gormData, 201)
}

func outputGormDataStatus(ctx *gin.Context, gormData interface{}, status int) {
  sch := GetModelSchema(ctx)
  jsonOutputData := mdl.ModelToMap(gormData, sch)
	DefaultOutput(ctx, status, jsonOutputData)
}
