# Introduction
This is a small plugin to help insert bulk data into mysql using gorm golang

# Get started
## Define your model
Mark the property that you want to insert, update in bulk query by using the tag `column:"insert,udpate"`
Example:
```
type Car struct {
	Id        int       `column:"insert"`
	Name      string    `column:"insert,update"`
	Property  Property  `json:"property" gorm:"type:json" column:"insert,update"`
	CreatedAt time.Time `column:"insert"`
	UpdatedBy string    `column:"insert"`
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

query, params := QueryBuilder().BuildInsertQuery(tableName, rows)
```

The generated `query` result as: 
```
insert into `cars` (`id`,`name`,`property`,`created_at`,`updated_by`) values (?,?,?,?,?)
```

And you can execute the query with the output param as follow:
```
err = db.Exec(query, params...).Error
```
- Insert on update duplicate
# Note:
If you have a large rows you must divide it into smaller job (such as: 200 rows/query).


