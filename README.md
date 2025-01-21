## How to build
After cloning the repository run 
```sh
go build -tags=jsoniter .
```
This will build the backend server with [jsoniter](https://github.com/json-iterator/go) which is a [faster alternative to encoding/json](https://github.com/json-iterator/go?tab=readme-ov-file#benchmark)

You can also build the backend server with debug data stripped:
```sh
-ldflags "-s -w"
```
