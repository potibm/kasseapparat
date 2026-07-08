package gorm

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuditModel struct {
	CreatedBy  string  `json:"created_by"           gorm:"column:created_by;type:varchar(255);<-:create"`
	ModifiedBy string  `json:"modified_by"          gorm:"column:modified_by;type:varchar(255)"`
	DeletedBy  *string `json:"deleted_by,omitempty" gorm:"column:deleted_by;type:varchar(255);index"`
}

func RegisterAuditCallbacks(db *gorm.DB) error {
	if err := db.Callback().
		Create().
		Before("gorm:create").
		Register("audit:before_create", beforeCreateCallback); err != nil {
		return err
	}

	if err := db.Callback().
		Update().
		Before("gorm:update").
		Register("audit:before_update", beforeUpdateCallback); err != nil {
		return err
	}

	if err := db.Callback().
		Delete().
		Before("gorm:delete").
		Register("audit:before_delete", beforeDeleteCallback); err != nil {
		return err
	}

	return nil
}

func getUserIDFromContext(tx *gorm.DB) string {
	if tx.Statement.Schema == nil {
		return ""
	}

	userID, ok := tx.Statement.Context.Value(UserIDKey).(string)
	if !ok {
		return ""
	}

	return userID
}

func setAuditColumn(tx *gorm.DB, columnName, userID string) {
	if field := tx.Statement.Schema.LookUpField(columnName); field != nil {
		tx.Statement.SetColumn(columnName, userID)
	}
}

func beforeCreateCallback(tx *gorm.DB) {
	if userID := getUserIDFromContext(tx); userID != "" {
		setAuditColumn(tx, "CreatedBy", userID)
		setAuditColumn(tx, "ModifiedBy", userID)
	}
}

func beforeUpdateCallback(tx *gorm.DB) {
	if userID := getUserIDFromContext(tx); userID != "" {
		setAuditColumn(tx, "ModifiedBy", userID)
	}
}

func beforeDeleteCallback(tx *gorm.DB) {
	userID := getUserIDFromContext(tx)
	if userID == "" {
		return
	}

	if tx.Statement.Schema == nil {
		return
	}

	field := tx.Statement.Schema.LookUpField("DeletedBy")
	if field == nil {
		return
	}

	updateQuery := tx.Session(&gorm.Session{NewDB: true})

	if tx.Statement.Dest != nil {
		updateQuery = updateQuery.Model(tx.Statement.Dest)
	} else {
		updateQuery = updateQuery.Table(tx.Statement.Table)
	}

	if whereClause, ok := tx.Statement.Clauses["WHERE"]; ok {
		if whereExpr, ok := whereClause.Expression.(clause.Where); ok {
			updateQuery.Statement.AddClause(whereExpr)
		}
	}

	err := updateQuery.Updates(map[string]interface{}{
		field.DBName: userID,
	}).Error
	if err != nil {
		_ = tx.AddError(err)
	}

	setAuditColumn(tx, "DeletedBy", userID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
