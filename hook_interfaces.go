package gg

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
)

type BeforeCreateWithContexter interface {
  BeforeCreateWithContext(ctx *gin.Context, db *gorm.DB) error
}

type AfterCreateWithContexter interface {
  AfterCreateWithContext(ctx *gin.Context, db *gorm.DB) error
}

type BeforeUpdateWithContexter interface {
  BeforeUpdateWithContext(ctx *gin.Context, db *gorm.DB) error
}

type AfterUpdateWithContexter interface {
  AfterUpdateWithContext(ctx *gin.Context, db *gorm.DB) error
}

type BeforeDeleteWithContexter interface {
  BeforeDeleteWithContext(ctx *gin.Context, db *gorm.DB) error
}

type AfterDeleteWithContexter interface {
  AfterDeleteWithContext(ctx *gin.Context, db *gorm.DB) error
}

type AfterFindWithContexter interface {
  AfterFindWithContext(ctx *gin.Context, db *gorm.DB) error
}
