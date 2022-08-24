package gormPlugin

import (
	"fmt"
	"strings"

	"github.com/SkyAPM/go2sky"
	"gorm.io/gorm"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentMySQL = 5

func GormCallback(tracer *go2sky.Tracer, dbName string) func(db *gorm.DB) {
	if tracer == nil {
		return func(db *gorm.DB) {}
	}
	return func(db *gorm.DB) {
		sql := fmt.Sprintf("%s, %v", db.Statement.SQL.String(), db.Statement.Vars)
		tableName := db.Statement.Table
		if tableName == "" {
			tableName = "undefined"
		}
		span, _ := tracer.CreateExitSpan(db.Statement.Context, tableName, dbName, func(key, value string) error {
			return nil
		})
		span.SetComponent(componentMySQL)
		span.SetSpanLayer(agentv3.SpanLayer_Database)
		span.Tag(go2sky.TagDBStatement, sql)
		defer span.End()
	}
}

func RegisterAll(db *gorm.DB, tracer *go2sky.Tracer, dbDsn string, callback func(*go2sky.Tracer, string) func(db *gorm.DB)) {
	// get db name if it's dsn
	dbName := dbDsn
	prefix := strings.Split(dbDsn, "?")[0]
	splited := strings.Split(prefix, "/")
	if len(splited) > 1 {
		dbName = splited[1]
	}

	db.Callback().Create().Register("skywalking", callback(tracer, dbName))
	db.Callback().Query().Register("skywalking", callback(tracer, dbName))
	db.Callback().Update().Register("skywalking", callback(tracer, dbName))
	db.Callback().Delete().Register("skywalking", callback(tracer, dbName))
	db.Callback().Row().Register("skywalking", callback(tracer, dbName))
	db.Callback().Raw().Register("skywalking", callback(tracer, dbName))
}
