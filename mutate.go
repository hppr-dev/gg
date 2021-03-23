package gg

import (
  "encoding/json"
  "gorm.io/gorm/schema"
  "github.com/gin-gonic/gin"
  "hppr.dev/gg/mdl"
)

type Mutator func(interface{}) interface{}

// MutateByID returns a handler that mutates, i.e. changes, the model using the given function then pushes it back to the database
func MutateByID(urlParam string, mutator Mutator) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    model := GetModel(ctx)
    sch := GetModelSchema(ctx)
    instance := mdl.NewModel(model)
    db := GetDatabase(ctx)
    pKeyColumn := sch.PrioritizedPrimaryField.DBName
    result := db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).Find(instance)
    if result.RowsAffected == 0 {
      DefaultOutput(ctx, 404, gin.H{"error": "not found"})
      return
    }
    updated := mutator(instance)
    db.Model(model).Where(pKeyColumn + " = ?", ctx.Param(urlParam)).Updates(updated)
    DefaultOutput(ctx, 200, convertStructToOutMap(updated, sch))
  }
}

func convertStructToOutMap(model interface{}, sch schema.Schema) map[string]interface{}{
  var jsonOut map[string]interface{}
  jsonBytes, _ := json.Marshal(model)
  json.Unmarshal(jsonBytes, &jsonOut)
  conv := make(map[string]interface{})
  for _, field := range sch.Fields {
    conv[field.DBName] = jsonOut[field.Name]
  }
  return conv
}
