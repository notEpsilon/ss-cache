# ss-cache
Zero-dependency Distributed statically-sharded in-memory cache based on a thread-safe generic LRU cache.

## Usage
- clone the repo and run `go build` to get the binary.
- running the binary spins up a new server with a cache instance on a specific port.
- check `main.go` flags section to know what flags are possible to provide when you run the binary.
- you need to have a `config.json` file somewhere to specify shards configuration and use it when spinning up new servers. (if not provided it defaults to `./config.json`).

## CLI
Although you can use it, the project is still in it's early stages and a CLI isn't yet provided, so you have to use the binary directly.
