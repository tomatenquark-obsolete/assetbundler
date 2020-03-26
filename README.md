assetbundler
------------

The `assetbundler` allows for downloading map contents of a Tomatenquark server.

It works in following simple steps:

- if server has the `servercontent` variable set it will send a `N_SERVERCONTENT` packet
- the client receives this packet on `N_MAPCHANGE`. If the client has `downloadmaps` enabled, it will start downloading
- the `libassetbundler` library will download all the configs referenced from `map.cfg` and collect all the resources defined in these files
- after this is done `libassetbundler` creates a temporary ZIP archive and uses `addzip` to load this into the game

Why go? Why not write it in C?

This is a quite obvious choice when you think about it:

- we want to support multiple storage backends in the future, not just HTTP
- adding those as C/C++ libraries is a major refactoring task of the build system and a lot of avoidable work
- since both `CFG` files (plain text) and `HTTP` servers are a well known technology it makes testing in isolation (without a tomato client) a nice experience
- it allows independent improvement and development as a standalone tool without need to iterate the client much

## How to build

```
git clone https://github.com/tomatenquark/assetbundler.git`
cd assetbundler/
go build -o libassetbundler.[so|dylib|dll] -buildmode=c-shared pkg/assetbundler/assetbundler.go # given you have a recent version of go installed
```
