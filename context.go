package gg

import (
  "reflect"
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "gorm.io/gorm/schema"
  "hppr.dev/gg/mdl"
)

type DefaultMap = map[string]interface{}

type ModelSchemaMap map[reflect.Type]schema.Schema

func (m ModelSchemaMap) GetSchema(model interface{}) schema.Schema {
  return m[reflect.TypeOf(model)]
}

func GetModelSchema(ctx *gin.Context) schema.Schema {
  return getContextVar("ModelSchema", ctx).(schema.Schema)
}

func GetModel(ctx *gin.Context) interface{} {
  return getContextVar("Model", ctx)
}

func GetDatabase(ctx *gin.Context) *gorm.DB {
  return getContextVar("DB", ctx).(*gorm.DB)
}

func GetSchemaMap(ctx *gin.Context) ModelSchemaMap {
  return getContextVar("SchemaMap", ctx).(ModelSchemaMap)
}

func GetDefaultOutputFunction(ctx *gin.Context) output {
  return getContextVar("DefaultOutput", ctx).(func(*gin.Context) output)(ctx)
}

func getContextVar(key string, ctx *gin.Context) interface{}{
  value, exist := ctx.Get(key)
  if !exist {
    panic("Failed to get " + key + " from gin context")
  }
  return value
}

func setupSchemaMap(cfg Config, db *gorm.DB) ModelSchemaMap {
  var schemaMap = make(ModelSchemaMap)
  for _, model := range cfg.Models {
    schemaMap[reflect.TypeOf(model)] = mdl.GetDBSchema(model, db)
  }
  return schemaMap
}
