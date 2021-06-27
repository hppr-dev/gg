package mdl

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

type ModelColumnInfo map[string]ColumnInfo

type ColumnInfo struct{
  DBName string
  StructName string
  Required bool
  HasDefault bool
  ColumnType reflect.Type
  SchemaField *schema.Field
}

func NewModel(model interface{}) interface{} {
	return reflect.New(followPtr(model).Type()).Interface()
}

func followPtr(model interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(model))
}

func GetDBSchema(model interface{}, db *gorm.DB) schema.Schema {
	var count int64
	return *db.Model(model).Count(&count).Statement.Schema
}

func GetModelColumnInfo(model interface{}, db *gorm.DB) ModelColumnInfo {
  modelSchema := GetDBSchema(model, db)
	names := make(ModelColumnInfo)
	for _, field := range modelSchema.Fields {
		names[field.DBName] = ColumnInfo{
      DBName     : field.DBName,
      StructName : field.Name,
      Required   : !isAutoCreatable(field),
      HasDefault : field.HasDefaultValue,
      ColumnType : field.FieldType,
      SchemaField: field,
    }
	}
	return names
}


func ExtractModelIDs(schema schema.Schema, anonModel interface{}) map[string]interface{} {
	primaryFields := schema.PrimaryFields
	values := make(map[string]interface{})
	for _, field := range primaryFields {
		values[field.DBName] = getModelField(anonModel, field.Name)
	}
	return values
}

func getModelField(model interface{}, fieldName string) interface{} {
	return followPtr(model).FieldByName(fieldName).Interface()
}
