package gg

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"hppr.dev/gg/mdl"
	"reflect"
)

type DefaultMap = map[string]interface{}

type ModelSchemaMap map[reflect.Type]schema.Schema

// GetSchema get the given gorm model schema from the schema map
func (m ModelSchemaMap) GetSchema(model interface{}) schema.Schema {
	return m[reflect.TypeOf(model)]
}

// GetModelSchema is a utility function that retrieves the gorm schema from the gin.Context
func GetModelSchema(ctx *gin.Context) schema.Schema {
	return getContextVar("ModelSchema", ctx).(schema.Schema)
}

// GetModel is a utility funciton that retrieves the current model from the gin.Context
func GetModel(ctx *gin.Context) interface{} {
	return getContextVar("Model", ctx)
}

// GetDatabase is a utility function that retreives the gorm.DB instance from the gin.Context
func GetDatabase(ctx *gin.Context) *gorm.DB {
	return getContextVar("DB", ctx).(*gorm.DB)
}

func getSchemaMap(ctx *gin.Context) ModelSchemaMap {
	return getContextVar("SchemaMap", ctx).(ModelSchemaMap)
}

func getDefaultOutputFunction(ctx *gin.Context) output {
	return getContextVar("DefaultOutput", ctx).(func(*gin.Context) output)(ctx)
}

func getContextVar(key string, ctx *gin.Context) interface{} {
	value, _ := ctx.Get(key)
	return value
}

func setupSchemaMap(cfg Config, db *gorm.DB) ModelSchemaMap {
	var schemaMap = make(ModelSchemaMap)
	for _, model := range cfg.Models {
		schemaMap[reflect.TypeOf(model)] = mdl.GetDBSchema(model, db)
	}
	return schemaMap
}
