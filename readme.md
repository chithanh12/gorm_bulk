# Introduction
This is a small plugin to help insert bulk data into mysql using gorm golang

# Get started
## Define your model
Mark the property that you want to insert, update in bulk query by using the tag `column:"insert,udpate"`
Example:
```
type Car struct {
	Id        int       
	Name      string   
	Property  Property  `json:"property" gorm:"type:json"`
	CreatedAt time.Time 
	UpdatedBy string   
}
```
## Build query
- Insert query
```
tableName := db.NewScope(model).TableName()
car := &Car{
    Id:   1,    
    Name: "AAA", 
    Property: Property{
        Prop1: "p1",
        Prop2: "p2",
        Age:   1,
    },
    CreatedAt: time.Now(),
    UpdatedBy: "admin",
}

statement, err := QueryBuilder().BuildInsertQuery(tableName, rows)
```

The generated `statement.Query` result as: 
```
insert into `cars` (`id`,`name`,`property`,`created_at`,`updated_by`) values (?,?,?,?,?)
```

And you can execute the query with the output param as follow:
```
err = db.Exec(statement.Query, statement.Parameters...).Error
```
- Insert on update duplicate
# Note:
If you have a large rows you must divide it into smaller job (such as: 200 rows/query).


