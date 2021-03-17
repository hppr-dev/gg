package gg

import (
  "github.com/gin-gonic/gin"
)

type Mutator func(interface{}) interface{}

func MutateByID(url_param string, mutator Mutator) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    DefaultOutput(ctx, 200, "Not Implemented")
  }
}
