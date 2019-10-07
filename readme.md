# Introduction
This is a small plugin to help insert bulk data into mysql using gorm golang

# Get started
## Define your model
Mark the property that you want to insert, update in bulk query by using the tag `column:"insert,udpate"`
Example:
```
type Car struct {
    Id int  // use generated on insert so that we dont need to update it
    Name string `column:"insert,update"`
    Model string `column:"insert,update"`
    Status string `column:"update"`
}
```
## Build query
- Insert query
- Insert on update duplicate


