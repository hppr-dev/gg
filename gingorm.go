package gg

import (
  "hppr.dev/gg/mdl"
	"github.com/gin-gonic/gin"
)

// SetModel configures the gin context to have references to the model and model schema
func SetModel(model interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schemaMap := getSchemaMap(ctx)
    schema := schemaMap.GetSchema(model)
		ctx.Set("ModelSchema", schema)
		ctx.Set("MarshalType", mdl.CreateMarshalType(model, schema))
		ctx.Set("Model", model)
		ctx.Next()
	}
}

// Middleware configures gin to use the gingorm middleware
func Middleware(cfg Config) gin.HandlerFunc {
	db, _ := cfg.OpenDB()
	schemaMap := setupSchemaMap(cfg, db)
	defaultOutput := getOutputFunction(cfg.DefaultOutputFormat)
	return func(ctx *gin.Context) {
		ctx.Set("DB", db)
		ctx.Set("SchemaMap", schemaMap)
		ctx.Set("DefaultOutput", defaultOutput)
		ctx.Next()
	}
}
