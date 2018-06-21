# GO ES

### Golang event sourcing made easy 
[![GoDoc](https://godoc.org/github.com/z0mbie42/goes?status.svg)](https://godoc.org/github.com/z0mbie42/goes)
[![GitHub release](https://img.shields.io/github/release/z0mbie42/goes.svg)](https://github.com/z0mbie42/goes/releases)

## Usage

**See `_examples/user/main.go` for an usage example**

```bash
$ docker run -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
$ export DATABASE="postgres://postgres:mysecretpassword@localhost/?sslmode=disable"
$ psql $DATABASE -c "CREATE DATABASE goes"
$ export DATABASE="postgres://postgres:mysecretpassword@localhost/goes?sslmode=disable"
$ cd _examples/user
$ go get -u
$ go run main.go
```

## Notes

`Apply` methods should return a pointer
`Validate` methods take a pointer as input

## Todo

- [x] First draft that avoid multiple switch
- [ ] Stop using pointers as an arguments and return values for the `Call` function: use purely immutables aggregates. (for the moment they are actually immutables, but you need to pass a pointer to persis in `gorm`, you can save an interface (which is not a concrete type).

## Glossary

* **Commands** Commands are responsible for: Validating attributes, Validating that the action can be performed given the current state of the application and Building the event. A `Command` can only return 1 `Event`, but it can be return mutiple `Event` types.

* **Events** are the source of truth. They are applied to `Aggregates`

* **Aggregates** represent the current state of the application. They are like models.

* **Calculators** to update the state of the application. This is the `Apply` method of the `Aggregate` interface. There is `Sync Reactors` which are called synchronously in the `Call` funciton, and can fail the transaction if an error occur, and `Async Reactor` which are called asynchronously, and are not checked for error (fire and forget).

* **Reactors** to trigger side effects as events happen. They are registered with the `On` Function.

* **Event Store** PostgresSQL


## Resources

This implementation is sort of the Go implementation of the following event sourcing framework

* https://kickstarter.engineering/event-sourcing-made-simple-4a2625113224
Because of the Go type system, i wasn't able (you can help ?) to use purely immutable aggregates:
You need to pass a pointer to the `Call` function. The underlying data is not modified, but is kind of dirty.

* https://github.com/mishudark/eventhus

