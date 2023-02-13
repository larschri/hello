# atomicmap

Use atomic.Pointer to store a map for read-only access. No fancy libraries
required for thread safety.

`main.go` start go-routines to update and read the map - `-race` should not
detect race condition.
