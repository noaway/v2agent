# godao
postgres orm module

# Usage

call `InitOrm()` first for initialize global postgres connection!

```
import "github.com/noaway/godao"

// init with config
godao.InitOrm(config)

// then use godao.Engine
godao.Engine.Model()...
```

see `example/main.go`