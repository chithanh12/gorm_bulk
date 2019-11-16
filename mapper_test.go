package gorm_bulk

import (
	"encoding/json"
	"testing"
	"time"
)

type Property struct {
	Prop1 string `json:"prop1"`
	Prop2 string `json:"prop2"`
	Age   int    `json:"age"`
}

type Car struct {
	Id        int
	Name      string
	Property  Property `json:"property" gorm:"type:json"`
	CreatedAt time.Time
	UpdatedBy string
}

func TestMapperGetColumn(t *testing.T) {
	now := time.Now()

	car := &Car{
		Id:   1,
		Name: "AAA",
		Property: Property{
			Prop1: "p1",
			Prop2: "p2",
			Age:   1,
		},
		CreatedAt: now,
		UpdatedBy: "nct",
	}

	// Test get insert column
	insertColumns, err := Mapper().GetColumns(car)
	if err != nil {
		t.Error(err)
	}
	if len(insertColumns) != 5 {
		t.Errorf("The number of insert columns is not correct: expected %v but get %v", 5, len(insertColumns))
	}

	if !StringHelper().SliceContains(insertColumns, "id") {
		t.Error("Missing column `id`")
	}

	if !StringHelper().SliceContains(insertColumns, "name") {
		t.Error("Missing column `name`")
	}

	if !StringHelper().SliceContains(insertColumns, "property") {
		t.Error("Missing column `property`")
	}

	if !StringHelper().SliceContains(insertColumns, "updated_by") {
		t.Error("Missing column `updated_by`")
	}

	if !StringHelper().SliceContains(insertColumns, "created_at") {
		t.Error("Missing column `created_at`")
	}
}

func TestGetInsertParams(t *testing.T) {
	now := time.Now()

	car := &Car{
		Id:   1,
		Name: "AAA",
		Property: Property{
			Prop1: "p1",
			Prop2: "p2",
			Age:   1,
		},
		CreatedAt: now,
		UpdatedBy: "nct",
	}

	param, err := Mapper().GetValues(car)
	if err != nil {
		t.Error(err)
	}

	if len(param) != 5 {
		t.Errorf("Number of params is not correct. Expected 4 but got %v", len(param))
	}

	if param["id"] != 1 {
		t.Errorf("Name value is not correct %v", param["name"])
	}

	if param["name"] != "AAA" {
		t.Errorf("Name value is not correct %v", param["name"])
	}

	if param["updated_by"] != "nct" {
		t.Errorf("Created by value is not correct %v", param["updated_by"])
	}

	if param["created_at"] != now {
		t.Errorf("time value is not correct")
	}

	property, err := json.Marshal(car.Property)
	if err != nil {
		t.Errorf("Can not marshal property value %v", err)
	}

	if param["property"] != string(property) {
		t.Error("Can not get property `property`")
	}
}

func TestBuildInsertQuery(t *testing.T) {
	now := time.Now()
	now1 := now.Add(1 * time.Minute)
	var rows []interface{}

	rows = append(rows, &Car{
		Id:   1,
		Name: "aa",
		Property: Property{
			Prop1: "p11",
			Prop2: "p12",
			Age:   11,
		},
		CreatedAt: now,
		UpdatedBy: "nct1"})
	rows = append(rows,
		&Car{
			Id:   2,
			Name: "bb",
			Property: Property{
				Prop1: "p21",
				Prop2: "p22",
				Age:   12,
			},
			CreatedAt: now1,
			UpdatedBy: "nct2",
		},
	)

	query, err := QueryBuilder().BuildInsertQuery("cars", rows)
	if err != nil {
		t.Error(err)
	}
	if query == nil {
		t.Error("error query is nill")
	}
	expectedQuery := "insert into `cars` (`id`,`name`,`property`,`created_at`,`updated_by`) values (?,?,?,?,?), (?,?,?,?,?)"

	// check syntax
	if query.Query != expectedQuery {
		t.Error("Generated query is not correct:")
		t.Log("Expected query:", expectedQuery)
		t.Log("Generated query:", query)
	}

	//check params
	prop1, _ := json.Marshal(rows[0].(*Car).Property)
	prop2, _ := json.Marshal(rows[1].(*Car).Property)

	expectedParams := []interface{}{}

	expectedParams = append(expectedParams, 1)
	expectedParams = append(expectedParams, "aa")
	expectedParams = append(expectedParams, string(prop1))
	expectedParams = append(expectedParams, now)
	expectedParams = append(expectedParams, "nct1")

	expectedParams = append(expectedParams, 2)
	expectedParams = append(expectedParams, "bb")
	expectedParams = append(expectedParams, string(prop2))
	expectedParams = append(expectedParams, now1)
	expectedParams = append(expectedParams, "nct2")

	if len(query.Parameters) != len(expectedParams) {
		t.Errorf("Number of params is not correct. Expedted %v rows but got %v", len(expectedParams), len(query.Parameters))
	}

	for idx, _ := range query.Parameters {
		if query.Parameters[idx] != expectedParams[idx] {
			t.Errorf("Expected param at %v is %v but got %v", idx, expectedParams[idx], query.Parameters[idx])
		}
	}
}

func TestBuildInsertOnDuplicateQuery(t *testing.T) {
	now := time.Now()
	now1 := now.Add(1 * time.Minute)
	var rows []interface{}

	rows = append(rows, &Car{Id: 1,
		Name: "aa",
		Property: Property{
			Prop1: "p11",
			Prop2: "p12",
			Age:   11,
		},
		CreatedAt: now,
		UpdatedBy: "nct1"})
	rows = append(rows, &Car{
		Id:   2,
		Name: "bb",
		Property: Property{
			Prop1: "p21",
			Prop2: "p22",
			Age:   12,
		},
		CreatedAt: now1,
		UpdatedBy: "nct2",
	},
	)

	query, err := QueryBuilder().BuildInsertOnDuplicateUpdate("cars", rows)
	if err != nil {
		t.Error(err)
	}

	expectedQuery := "insert into `cars` (`id`,`name`,`property`,`created_at`,`updated_by`) values (?,?,?,?,?), (?,?,?,?,?) " +
		"on duplicate key update `id`=values(`id`),`name`=values(`name`),`property`=values(`property`),`created_at`=values(`created_at`),`updated_by`=values(`updated_by`)"

	// check syntax
	if query.Query != expectedQuery {
		t.Error("Generated query is not correct:")
		t.Log("Expected query:", expectedQuery)
		t.Log("Generated query:", query.Query)
	}

	//check params

	prop1, _ := json.Marshal(rows[0].(*Car).Property)
	prop2, _ := json.Marshal(rows[1].(*Car).Property)

	expectedParams := []interface{}{}

	expectedParams = append(expectedParams, 1)
	expectedParams = append(expectedParams, "aa")
	expectedParams = append(expectedParams, string(prop1))
	expectedParams = append(expectedParams, now)
	expectedParams = append(expectedParams, "nct1")

	expectedParams = append(expectedParams, 2)
	expectedParams = append(expectedParams, "bb")
	expectedParams = append(expectedParams, string(prop2))
	expectedParams = append(expectedParams, now1)
	expectedParams = append(expectedParams, "nct2")

	if len(query.Parameters) != len(expectedParams) {
		t.Errorf("Number of params is not correct. Expedted %v rows but got %v", len(expectedParams), len(query.Parameters))
	}

	for idx, _ := range query.Parameters {
		if query.Parameters[idx] != expectedParams[idx] {
			t.Errorf("Expected param at %v is %v but got %v", idx, expectedParams[idx], query.Parameters[idx])
		}
	}
}
