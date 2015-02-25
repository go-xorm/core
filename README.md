[中文](https://github.com/go-xorm/xorm/blob/master/README_CN.md)

Core is a light wrapper of sql.DB.

# Open
```Go
db, _ := core.Open(db, connstr)
```

# SetMapper
```Go
db.SetMapper(SameMapper())
```

# More scan usage
```Go

rows, _ := db.Query()
for rows.Next() {
    rows.Scan()
    rows.ScanMap()
    rows.ScanSlice()
    rows.ScanStructByName()
    rows.ScanStructByIndex()
}
```

# More Query usage
```Go
rows, err := db.Query("select * from table where name = ?", name)

rows, err := db.QueryStruct("select * from table where name = ?Name",
            &user)

var user = map[string]interface{}{
    "name": "lunny",
}
rows, err = db.QueryMap("select * from table where name = ?name",
            &user)
```

# More QueryRow usage
```Go
rows, err := db.QueryRow("select * from table where name = ?", name)

rows, err := db.QueryRowStruct("select * from table where name = ?Name",
            &user)
var user = map[string]interface{}{
    "name": "lunny",
}
rows, err = db.QueryRowMap("select * from table where name = ?name",
            &user)
```

# More Exec usage
```Go
db.Exec("insert into user (`name`, title, age, alias, nick_name,created) values (?,?,?,?,?,?)", name, title, age, alias...)

user = User{
    Name:"lunny",
    Title:"test",
    Age: 18,
}
result, err = db.ExecStruct("insert into user (`name`, title, age, alias, nick_name,created) values (?Name,?Title,?Age,?Alias,?NickName,?Created)",
            &user)

var user = map[string]interface{}{
    "Name": "lunny",
    "Title": "test",
    "Age": 18,
}
result, err = db.ExecMap("insert into user (`name`, title, age, alias, nick_name,created) values (?Name,?Title,?Age,?Alias,?NickName,?Created)",
            &user)
```