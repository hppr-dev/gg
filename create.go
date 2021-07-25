package gg

import (
	"github.com/gin-gonic/gin"
	"hppr.dev/gg/mdl"
)

// BodyCreate returns a handler that creates a model record based on post data
func BodyCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model := GetModel(ctx)
		marshal := GetMarshalType(ctx)
    jsonInputData := mdl.New(marshal)
		gormData := mdl.New(model)
		if err := ctx.ShouldBind(jsonInputData) ; err != nil {
      DefaultOutput(ctx, 400, gin.H{"message" : "Failed to parse data."})
      return
    }
    mdl.CopyFields(jsonInputData, gormData)
		db := GetDatabase(ctx)
    if cast, ok := gormData.(BeforeCreateWithContexter) ; ok {
      if err := cast.BeforeCreateWithContext(ctx, db) ; err != nil {
        return
      }
    }
		db.Model(model).Create(gormData)
    if cast, ok := gormData.(AfterCreateWithContexter) ; ok {
      if err := cast.AfterCreateWithContext(ctx, db) ; err != nil {
        return
      }
    }
    jsonOutputData := mdl.New(marshal)
    mdl.CopyFields(gormData, jsonOutputData)
		DefaultOutput(ctx, 201, jsonOutputData)
	}
}
