package models

import (
	gormaudit "github.com/potibm/kasseapparat/internal/app/store/gorm"
)

type GormOwnedModel struct {
	GormModel
	gormaudit.AuditModel
} // @name models.gormOwnedModel
