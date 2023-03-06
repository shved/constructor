# Constructor

Some REST APIs have filter parameters to query a particular list of resources fields included into response. Sometimes those parameters are quiet complicated and thats why GraphQL was introduced. But for good old APIs query language may become less or more complex and sometimes it isn't that trivial to build a query parameter.  

With Constructor is this easy:
```go
type HugeResourceStruct struct {
	Name    string  `json:"username"`
	Address Address `constructor:"omit"`
	Object  Object  `json:"obj"`
	// ...
	Field108 string `json:""`
}

type Object struct {
	Prop1 string
	Prop2 string
}

b := constructor.NewBuilder(constructor.Options{
	ParamKey:       "select",
	Delimiter:      constructor.DefaultDelimiter,
	FieldDelimiter: constructor.DefaultFieldDelimiter,
})

queryParam := b.QueryStringFromStruct(HugeResourceStruct{})
```
This call will produce a string of the following format:  
`select=username,Object*Prop1,Object*Prop2`  

It tries to take json tags of a given structure. If json tag is explicitly set to empty, it is ignored. If there is no json tag, the name will be the same as name of the field itself. You may also explicitly ignore the field setting constructor tag to "omit". The field will also be ignored if it is unexported. Find more usage examples in tests.
