# varchar
```regexp
([a-z]*)\s*varchar\((\d+)\)

field.String("$1").MaxLen($2).StructTag(`json:"$1"`)
```


# bigint

```regexp
([a-z_0-9]*)\s*bigint
field.Int8("$1") .StructTag(`json:"$1"`)
```

# num
```regexp
([a-z_0-9]*)\s*(numeric[^)]*\)*)
field.Int8("$1"). GoType(new(BigInt)) .StructTag(`json:"$1"`).SchemaType(map[string]string{ dialect.Postgres: "$2", })
```