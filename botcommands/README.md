# README
Purpose of this package is to implement functions that work standalone. Working off the
parameters `context.Context` and `io.ReadWriter` the struct that implements the interface
`discordbot.Command` should contain the dependency injections to fulfill its logic.

Ideally this should make it easier to test each command in isolation by only having to assert
on the `io.ReadWriter`.
