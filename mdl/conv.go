package mdl

import (
  "fmt"
  "reflect"
  "gorm.io/gorm/schema"
)

// Creates a runtime type that has form, json, xml tags for marshalling to/from json
// By default uses the db name from the schema.
func CreateMarshalType(model interface{}, sch schema.Schema) interface{} {
  t := reflect.TypeOf(model).Elem()
  fields := make([]reflect.StructField, t.NumField())
  for i := 0; i < len(fields) ; i++ {
    field := t.Field(i)
    if field.Type.Kind() != reflect.Struct {
      fields[i] = reflect.StructField{
        Name : field.Name,
        Type : field.Type,
        Tag  : addDefaultTags(sch.FieldsByName[field.Name], field.Tag),
      }
    } else {
      fields[i] = reflect.StructField{
        Name : field.Name,
        Type : field.Type,
        Tag  : "form:\"-\" json:\"-\" xml:\"-\"",
      }
    }
  }
  return reflect.New(reflect.StructOf(fields)).Elem().Addr().Interface()
}

// Copy equivalent fields from one interface{} to another
// Fields must have the same field names and the interfaces should be pointers to model-like structs
func CopyFields(src, dest interface{}) {
  typeDef := reflect.TypeOf(src).Elem()
  srcStruct := followPtr(src)
  destStruct := followPtr(dest)
  for i := 0 ; i < typeDef.NumField() ; i++ {
    destStruct.FieldByName(typeDef.Field(i).Name).Set(srcStruct.Field(i))
  }
}

//Convert gorm model to map
func ModelToMap(model interface{}, sch schema.Schema) map[string]interface{} {
  outMap := make(map[string]interface{})
  modelStruct := followPtr(model)
  for _, field := range sch.Fields {
    outMap[field.DBName] = modelStruct.FieldByName(field.Name).Interface()
  }
  return outMap
}

//Convert map to gorm Model
func MapToModel(inMap map[string]interface{}, model interface{}, sch schema.Schema) {
  modelStruct := followPtr(model)
  for key, value := range inMap {
    field := sch.LookUpField(key)
    modelStruct.FieldByName(field.Name).Set(reflect.ValueOf(value))
  }
}

func addDefaultTags(schemaField *schema.Field, tag reflect.StructTag) reflect.StructTag {
  if schemaField != nil {
    computedTags := string(tag)
    defaultName := schemaField.DBName
    computedTags += defaultTag(tag, "form", defaultName)
    computedTags += defaultTag(tag, "json", defaultName)
    computedTags += defaultTag(tag, "xml", defaultName)
    computedTags += bindingTag(tag, schemaField)
    return reflect.StructTag(computedTags)
  }
  return "form:\"-\" json:\"-\" xml:\"-\""
}

func bindingTag(tag reflect.StructTag, field *schema.Field) string {
  if _, found := tag.Lookup("binding") ; !found && !isAutoCreatable(field) {
    return "binding:\"required\""
  }
  return ""
}

func defaultTag(tag reflect.StructTag, tagName, defaultValue string) string {
  if _, found := tag.Lookup(tagName) ; !found {
    return fmt.Sprintf(" %s:\"%s\"", tagName, defaultValue)
  }
  return ""
}
