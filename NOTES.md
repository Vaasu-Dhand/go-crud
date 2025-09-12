## Go Struct Return Values: nil vs Error Patterns

In Go, you **can't return nil for a struct value** (like your Movie), because Movie is not a pointer type. You can only return nil if you return a pointer (*Movie).

So you have two main options:

## Return (Movie, error)

This is the most common Go pattern when the absence of a value is an error-like condition.

```go
func GetById(id string) (Movie, error) {
	for _, movie := range Movies {
		if movie.Id == id {
			return movie, nil
		}
	}
	return Movie{}, fmt.Errorf("movie with id %s not found", id)
}
```

Usage:

```go
m, err := GetById("123")
if err != nil {
	fmt.Println(err)
} else {
	fmt.Println(m.Title)
}
```

✅ **Pros:** Explicit error handling, idiomatic.  
❌ **Con:** You always need to check the error.
