package gormPlugin

import (
	"fmt"

	"github.com/SkyAPM/go2sky"
	agent "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"gorm.io/gorm"
)

const componentMySQL = 5

func GormCallback(tracer *go2sky.Tracer, dbDsn string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		sql := fmt.Sprintf("%s, %v", db.Statement.SQL.String(), db.Statement.Vars)
		span, _ := tracer.CreateExitSpan(db.Statement.Context, db.Statement.Table, dbDsn, func(header string) error {
			return nil
		})
		span.SetComponent(5)
		span.SetSpanLayer(agent.SpanLayer_Database)
		span.Tag(go2sky.TagDBStatement, sql)
		defer span.End()
	}
}

func RegisterAll(db *gorm.DB, tracer *go2sky.Tracer, dbDsn string, callback func(*go2sky.Tracer, string) func(db *gorm.DB)) {
	db.Callback().Create().Register("skywalking", callback(tracer, dbDsn))
	db.Callback().Query().Register("skywalking", callback(tracer, dbDsn))
	db.Callback().Update().Register("skywalking", callback(tracer, dbDsn))
	db.Callback().Delete().Register("skywalking", callback(tracer, dbDsn))
	db.Callback().Row().Register("skywalking", callback(tracer, dbDsn))
	db.Callback().Raw().Register("skywalking", callback(tracer, dbDsn))
}
