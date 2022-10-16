package gg

import (
	"github.com/gin-gonic/gin"
  "hppr.dev/gg/mdl"
  "fmt"
)

// UpdateByID returns a handler that uses post data to update the model record with the id from the url parameter
func UpdateByID(urlParam string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schema := GetModelSchema(ctx)
    var bodyData = make(DefaultMap)
    bindData(ctx, &bodyData)
    fmt.Printf("gormData:%+v\n", bodyData)
    if err := mdl.MatchAnyMapToModel(bodyData, schema) ; err != nil {
			DefaultOutput(ctx, 400, gin.H{"error": err.Error()})
			return
    }
    model := GetModel(ctx)
    gormData := mdl.New(model)
    mdl.MapToModel(bodyData, gormData, schema)
		pKeyColumn := schema.PrioritizedPrimaryField.DBName
		db := GetDatabase(ctx)
    if cast, ok := gormData.(BeforeUpdateWithContexter) ; ok {
      if err := cast.BeforeUpdateWithContext(ctx, db); err != nil {
        return
      }
    }
		result := db.Model(gormData).Where(pKeyColumn+" = ?", ctx.Param(urlParam)).Updates(gormData)
		if result.RowsAffected == 0 {
			DefaultOutput(ctx, 404, gin.H{"error": "not found"})
			return
		}
    if cast, ok := gormData.(AfterUpdateWithContexter) ; ok {
      if err := cast.AfterUpdateWithContext(ctx, db); err != nil {
        return
      }
    }
		db.Model(gormData).Where(pKeyColumn+" = ?", ctx.Param(urlParam)).First(gormData)
    outputGormDataStatus(ctx, gormData, 200)
	}
}
