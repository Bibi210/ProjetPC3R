```go
/* type PrimaryKey struct {
	autoIncrement bool
}

type ForeignKey struct {
	ReferedTable Table
	ReferedField string
}

type Key interface {
	toSQL() string
}

type BasicField struct {
	t SQLType
}

type KeyField struct {
	field BasicField
	key   Key
}

type Field interface {
	toSQL() string
}

type SQLType interface {
	typeToSQL()
}

type Table struct {
	Name   string
	Fields map[string]Field
} */

```