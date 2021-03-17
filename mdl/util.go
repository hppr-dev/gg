package mdl

import (
  "reflect"
  "gorm.io/gorm"
  "gorm.io/gorm/schema"
)


func NewModel(model interface{}) interface{} {
  return reflect.New(followPtr(model).Type()).Interface()
}

func followPtr(model interface{}) reflect.Value{
  return reflect.Indirect(reflect.ValueOf(model))
}

func GetDBSchema(model interface{}, db *gorm.DB) schema.Schema {
  var count int64
  return *db.Model(model).Count(&count).Statement.Schema
}

func ExtractModelIDs(schema schema.Schema, anonModel interface{}) map[string]interface{} {
  primaryFields := schema.PrimaryFields
  values := make(map[string]interface{})
  for _, field := range primaryFields {
    values[field.DBName] = getModelField(anonModel, field.Name)
  }
  return values
}

func getModelField(model interface{}, fieldName string) interface{}{
  return followPtr(model).FieldByName(fieldName).Interface()
}
