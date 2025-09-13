## Go Struct Return Values: nil vs Error Patterns

In Go, you **can't return nil for a struct value** (like your Movie), because Movie is not a pointer type. You can only return nil if you return a pointer (*Movie).

So you have two main options:

### Return (Movie, error)

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

## json.Marshal vs json.Unmarshal
In Go's encoding/json package, Marshal and Unmarshal are fundamental functions for converting between Go data structures and JSON data.

### json.Marshal
json.Marshal: This function takes a Go value (like a struct, map, slice, or basic type) and converts it into a JSON-encoded byte slice. This process is known as marshalling or serialization. The output is a []byte representing the JSON data.

```go
    package main

    import (
    	"encoding/json"
    	"fmt"
    )

    type Person struct {
    	Name string `json:"name"`
    	Age  int    `json:"age"`
    }

    func main() {
    	p := Person{Name: "Alice", Age: 30}
    	jsonData, err := json.Marshal(p)
    	if err != nil {
    		fmt.Println(err)
    		return
    	}
    	fmt.Println(string(jsonData)) // Output: {"name":"Alice","age":30}
    }
```

### json.Unmarshal
Unmarshalling is the reverse process, converting a JSON string into a data structure using the json.Unmarshal function. The target v must be a pointer to a Go data structure where the parsed JSON data will be stored.

```go
    package main

    import (
    	"encoding/json"
    	"fmt"
    )

    type Person struct {
    	Name string `json:"name"`
    	Age  int    `json:"age"`
    }

    func main() {
    	jsonData := []byte(`{"name":"Bob","age":25}`)
    	var p Person
    	err := json.Unmarshal(jsonData, &p)
    	if err != nil {
    		fmt.Println(err)
    		return
    	}
    	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age) // Output: Name: Bob, Age: 25
    }
```