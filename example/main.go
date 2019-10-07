package main

import (
	"fmt"
	"time"

	"github.com/chithanh12/gorm_bulk"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	JsonProp struct {
		Prop1 string `json:"prop1"`
		Prop2 string `json:"prop2"`
	}

	SampleModel struct {
		Id        int64
		Name      string   `column:"insert,update"`
		Items     JsonProp `column:"insert,update" gorm:"type:json"`
		CreatedAt time.Time
	}
)

func main() {
	//Replace the value for your connection string
	db, err := gorm.Open("mysql", connectionString("localhost", "root", "root", "sample"))
	if err != nil {
		fmt.Println("Error when connect to db....")
		return
	}

	defer db.Close()
	tableName := db.NewScope(&SampleModel{}).TableName()

	createTable := fmt.Sprintf("create table if not exists %v("+
		"		`id` int auto_increment primary key,"+
		"		`name` varchar(36),"+
		"		`items` json null,"+
		"		`created_at` timestamp default current_timestamp"+
		"	)engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;", tableName)

	err = db.Exec(createTable).Error
	if err != nil {
		fmt.Println("Can not create sample model tables")
	}

	rows := []interface{}{
		&SampleModel{
			Id:   0,
			Name: "Sample Model 1",
			Items: JsonProp{
				Prop1: "prop 11",
				Prop2: "prop 12",
			},
			CreatedAt: time.Now(),
		},
		&SampleModel{
			Id:   0,
			Name: "Sample Model 2",
			Items: JsonProp{
				Prop1: "prop 21",
				Prop2: "prop 22",
			},
			CreatedAt: time.Now(),
		},
	}
	// Scenario 1: Insert bulk
	query, params := gorm_bulk.QueryBuilder().BuildInsertQuery(tableName, rows)
	err = db.Exec(query+";", params...).Error

	if err != nil {
		fmt.Errorf("Error when process insert bulk %v", err)
	}
}
func connectionString(host, user, password, db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&timeout=90s&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local&multiStatements=true",
		user, password, host, db)
}
