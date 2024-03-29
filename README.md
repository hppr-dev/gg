# Gin Gorm Middleware

Gin middleware to create api endpoints for gorm models.

NOTE: This package is still in development and is unreleased as of yet.

# Using The middleware

Activate the middleware by calling Use on your instance of gin.Engine (from gin.Default(), gin.New(), etc).
Then you may call the provided handler functions to set up basic endpoints.

Note: migration is outside of the scope of this project. Make sure to apply migrations before running the server.

Simple example server:

``` go

package main

import (
  "gorm.io/gorm"
  "github.com/gin-gonic/gin"
  "hppr.dev/gg"
)

type User struct {
  gorm.Model
  Name string
  Password string
}
type Post struct {
  gorm.Model
  Text string
  UserID uint64
  User User
}


func main() {
  engine := gin.Default()

  models := []interface{}{
    &User{},
    &Post{},
  }

  cfg := gg.Config{
    Gorm: &gorm.Config{},
    Database: &gg.SQLiteConfig{"user.db"},
    Models: models,
    DefaultOutputFormat: gg.IndentedJSON,
  }

  engine.Use(gg.Middleware(cfg))

  db, _ := cfg.OpenDB()
  db.AutoMigrate(models...)

  engine.GET("/user", gg.SetModel(&User{}), gg.QuerySearch())
  engine.POST("/user", gg.SetModel(&User{}), gg.BodyCreate())
  engine.DELETE("/user/:id", gg.SetModel(&User{}), gg.DeleteByID("id"))
  engine.GET("/post", gg.SetModel(&Post{}), gg.QuerySearch())
  engine.POST("/post", gg.SetModel(&Post{}), gg.BodyCreate())
  engine.DELETE("/post/:id", gg.SetModel(&Post{}), gg.DeleteByID("id"))

  engine.Run()
}

```

# Generic handler functions

## QuerySearch, BodySearch,

QuerySearch and BodySearch serve the same purpose: search tables.
QuerySearch uses the query string for search parameters.
BodySearch uses post data for search parameters.

``` go
  // curl "http://localhost/user?id_gte=5"
  engine.GET('/user', gg.SetModel(&User{}), gg.QuerySearch())

  // curl "http://localhost/user" -XPOST -d '{"id_gte" : 5}'
  engine.POST('/user', gg.SetModel(&User{}), gg.BodySearch())
```

## QuerySearchByColumn, BodySearchByColumn

Queries the database for entries with a column that matches the url parameter.

``` go
  engine.GET('/user/:q', gg.SetModel(&User{}), gg.QuerySearchByColumn("q", "group"))
```

Also supports search parameters in the column field:

``` go
  engine.GET('/user/:q', gg.SetModel(&User{}), gg.QuerySearchByColumn("q", "id_gte"))
```

### Search Parameters

Search parameters are made by adding a suffix to the db column name.
Available suffixs are :
* `_gte, _gt` - greater than or equal, greater than
* `_lte`,`_lt` - less than or equal, less than
* `_ne` - not equal
* `_contains` - string contains

Multiple search parameters are allowed on the same column name:

``` go
  curl "http://localhost/user?id_gte=5&id_lt=10"
```

Only AND conditions are supported for the time being.

## GetByID

Queries the database for an entry that has a primary key that matches the given url parameter.

``` go
  engine.GET('/user/:id', gg.SetModel(&User{}), gg.GetByID("id"))
```

## BodyCreate

The BodyCreate handler creates a model of a specific model type using post data.
All fields are required, except for those that are autogenerated or that have default values.
Note: If a type mismatch occurs, the zero value for the field will be used.

``` go
  engine.POST('/user', gg.SetModel(&User{}), gg.BodyCreate())
```

## DeleteByID

Deletes the entry with aprimary key matching an url parameter

``` go
  engine.DELETE('/user/:id', gg.SetModel(&User{}), gg.DeleteByID("id"))
```

## UpdateByID

Updates the entry with a primary key matching an url parameter.
Data is passed through request post data.

``` go
  engine.POST('/user/:id', gg.SetModel(&User{}), gg.UpdateByID("id"))
```

## MutateByID

Mutates a model using a mutator function.


``` go
  engine.GET('/user/:id/assign', gg.SetModel(&User{}), gg.MutateByID("id", func(m interface{}) interface{} { model := m.(User) })) 
```

# Middleware

The middleware sets three keys in the context:

* "DB" - the gorm database connection
* "ModelSchemaMap" - the mapping between model types and gorm Schemas
* "DefaultOutput" - the default output function to use

# SetModel

The SetModel handler sets two keys in the context:

* "Model" - The type of model the handler deals with
* "Schema" - The gorm schema of the given model

# Example Custom Configuration

The following would create two endpoints that provide the count of Users and Posts in the database

``` go
  func main() {
    engine := gin.Default()
    cfg := gg.Config{
      ...
    }
    engine.Use(gg.Middleware(cfg))
    engine.GET("/user/count", gg.SetModel(&User{}), GetCount)
    engine.GET("/post/count", gg.SetModel(&Post{}), GetCount)
    engine.Run()
  }

  func GetCount(ctx *gin.Context) {
    db := ctx.MustGet("DB").(*gorm.DB)
    model := ctx.MustGet("Model")
    count := db.Model(model).Count()
    ctx.JSON(200, gin.H{ 'count' : count })
  }
```

The GetCount function could be rewritten with provided convenience functions:

``` go
  func GetCount(ctx *gin.Context) {
    db := gg.GetDB(ctx)
    model := gg.GetModel(ctx)
    count := db.Model(model).Count()
    gg.DefaultOutput(ctx, 200, gin.H{ 'count': count})
  }
```

The above also has the added benefit of allowing the switching of output types based on the Config.DefaultOutputFormat value.

# Hooks

GinGorm provides similar model hooks to gorm.
GinGorm adds the ability to access the current *gin.Context, as well as the DB.

Note: These hooks happen at they're alloted time synchronously in the request process (i.e. either before or after the action).
This means:
- If any before hook panics, the models in question will not be committed.
- If any after hook panics, the requested action has already taken place on the database, and will not be rolled back.
For this reason it is advised to handle errors in after hooks and not let them bubble out to gin.

Note: Before hooks happen before the database is touched and after hooks happen after.
A side effect to this is:
- The models provided to the Before hooks will not be populated with anything
- The models provided to the After hooks will have values from the database and will only affect the returned representation of the model.

For all hooks except AfterFindWithContext if it returns a non-nil error, it is assumed that the context has been set with the appropriate response and request handling terminates.
If AfterFindWithContext returns a non-nil error, the server will return a status code 400 with the returned error message.

Available hooks:
- Model Creation:
  - BeforeCreateWithContext(ctx *gin.Context, db *gorm.DB) error 
  - AfterCreateWithContext(ctx *gin.Context, db *gorm.DB) error 

- Model Update:
  - BeforeUpdateWithContext(ctx *gin.Context, db *gorm.DB) error 
  - AfterUpdateWithContext(ctx *gin.Context, db *gorm.DB) error 

- Model Delete
  - BeforeDeleteWithContext(ctx *gin.Context, db *gorm.DB) error 
  - AfterDeleteWithContext(ctx *gin.Context, db *gorm.DB) error 

- Model Find
  - AfterFindWithContext(ctx *gin.Context, db *gorm.DB) error 
    - Note that this is called even if no results are found

If you do not need the context for model hooks see [gorm hooks](https://gorm.io/docs/hooks.html).
