package gg

import (
  "github.com/gin-gonic/gin"
  "hppr.dev/gg/database"
)

func SetModel(model interface{}) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    schemaMap := GetSchemaMap(ctx)
    ctx.Set("Schema", schemaMap.GetSchema(model))
    ctx.Set("Model", model)
    ctx.Next()
  }
}

func Middleware(cfg Config) gin.HandlerFunc {
  db, err := database.Open(cfg.Database.Configure(), cfg.OnDBOpen, cfg.Gorm)
  if err != nil {
    panic("Could not connect to database")
  }
  schemaMap := setupSchemaMap(cfg, db)
  defaultOutput := GetOutputFunction(cfg.DefaultOutputFormat)
  return func(ctx *gin.Context) {
    ctx.Set("DB", db)
    ctx.Set("SchemaMap", schemaMap)
    ctx.Set("DefaultOutput", defaultOutput)
    ctx.Next()
  }
}

