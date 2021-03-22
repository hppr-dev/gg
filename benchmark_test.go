package gg

import (
  "fmt"
  "testing"
  "gorm.io/gorm"
  "github.com/gin-gonic/gin"
  "gorm.io/driver/sqlite"
)


type BenchModel struct {
  gorm.Model
  Name string
}

func BenchmarkCreateGetDeleteGinGorm(b *testing.B) {
  config := Config{
    Gorm: &gorm.Config{},
    Database: InMemConfig{},
    Models: []interface{}{&BenchModel{}},
  }
  db, _ := config.OpenDB()
  db.Migrator().DropTable(&BenchModel{})
  db.AutoMigrate(&BenchModel{})

  engine := gin.New()
  gin.SetMode(gin.TestMode)

  engine.Use(Middleware(config))
  engine.GET("/:id", SetModel(&BenchModel{}), GetByID("id"))
  engine.POST("/", SetModel(&BenchModel{}), BodyCreate())
  engine.DELETE("/:id", SetModel(&BenchModel{}), DeleteByID("id"))

  body := buildJsonBytes("name", "testing")

  var idUrl string
  for n := 0 ; n < b.N ; n++ {
    resp := runRequest(engine, "POST", "/", body)
    mapResp := extractMapBody(resp)
    idUrl = fmt.Sprintf("/%.0f", mapResp["id"])
    runRequest(engine, "GET", idUrl, nil)
    runRequest(engine, "DELETE", idUrl, nil)
  }
}

func BenchmarkCreateGetDeleteBare(b *testing.B) {
  db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
  db.Migrator().DropTable(&BenchModel{})
  db.AutoMigrate(&BenchModel{})

  engine := gin.New()
  gin.SetMode(gin.TestMode)
  engine.GET("/:id", func(ctx *gin.Context) {
    id := ctx.Param("id")
    var model BenchModel
    db.First(&model, id)
    ctx.JSON(200, model)
  })
  engine.POST("/", func(ctx *gin.Context) {
    model := BenchModel{}
    ctx.BindJSON(&model)
    db.Create(&model)
    ctx.JSON(200, model)
  })
  engine.DELETE("/:id", func(ctx *gin.Context) {
    var model BenchModel
    db.Delete(&model, ctx.Param("id"))
    ctx.JSON(200, gin.H{})
  })

  body := buildJsonBytes("name", "testing")

  var idUrl string
  for n := 0 ; n < b.N ; n++ {
    resp := runRequest(engine, "POST", "/", body)
    mapResp := extractMapBody(resp)
    idUrl = fmt.Sprintf("/%.0f", mapResp["ID"])
    runRequest(engine, "GET", idUrl, nil)
    runRequest(engine, "DELETE", idUrl, nil)
  }
}
