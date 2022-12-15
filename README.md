### ent-bank.io

This is a supporting repository for a take-home assignment.

It contains a very dumbed down version of an API server that can create users and update their balance.

#### Running

To run the server:

* This is a template repository. [Create a repo from this template](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template)
* [Install Go](https://go.dev/doc/install) 
* `go run -mod=mod main.go`

#### Testing

Two unit tests exist in [main_test.go](main_test.go) and are there to help you quickly iterate on the modifications to
the server. To run them:

```text
go test ./...
```

#### Schema

The schema for the User entity is:

```go
package schema

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("email").
			Unique(),
		field.Float("balance"),
	}
}
```

It can be found in [ent/schema/user.go](ent/schema/user.go).

#### Endpoints

The server has two endpoints:

1. `/v1/user` - used to create a new User.

```shell
curl -X POST --location "http://localhost:8080/v1/user" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer pupu" \
    -d "{
        \"name\": \"rotem\", \"email\": \"r@t.com\"
        }"
```

2. `/v1/user/:id/balance` - used to update the balance of an existing user.

3. `/v1/user/:id/balance/:date` - TODO: Implement this. It should return the balance of the user at the given date.

```shell
curl -X PATCH --location "http://localhost:8080/v1/user/1/balance" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 1234" \
    -d "{
          \"delta\": -100
        }"
```

#### Auth

To simplify the implementation the server will accept any bearer token in the
`Authorization` header. If no token is supplied a 401 error is returned.
