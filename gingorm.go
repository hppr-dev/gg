package gg

import (
  "github.com/gin-gonic/gin"
)

func SetModel(model interface{}) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    schemaMap := GetSchemaMap(ctx)
    ctx.Set("ModelSchema", schemaMap.GetSchema(model))
    ctx.Set("Model", model)
    ctx.Next()
  }
}

func Middleware(cfg Config) gin.HandlerFunc {
  db, _ := cfg.OpenDB()
  schemaMap := setupSchemaMap(cfg, db)
  defaultOutput := GetOutputFunction(cfg.DefaultOutputFormat)
  return func(ctx *gin.Context) {
    ctx.Set("DB", db)
    ctx.Set("SchemaMap", schemaMap)
    ctx.Set("DefaultOutput", defaultOutput)
    ctx.Next()
  }
}

