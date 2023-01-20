# ss-cache
Zero-dependency Distributed statically-sharded in-memory key-value cache based on a thread-safe generic LRU cache.

## Usage
- clone the repo and run `go build` to get the binary.
- running the binary spins up a new server with a cache instance on a specific port.
- check `main.go` flags section to know what flags are possible to provide when you run the binary.
- you need to have a `config.json` file somewhere to specify shards configuration and use it when spinning up new servers. (if not provided it defaults to `./config.json`).
- issue requests to `/get?key=<something>` or `/set?key=<something>&value=<something>`
- the values are stored as an array of bytes for better performance and the client is responsible for decoding it back which shouldn't be hard.

## CLI
Although you can use it, the project is still in it's early stages and a CLI isn't yet provided, so you have to use the binary directly.
