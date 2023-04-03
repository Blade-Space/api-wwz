# wwf - working with files 📃

Package for working with files

## import to progect

```golang
package main

import (
  wwf "api/wwf/routes"

  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  api := r.Group("/api/wwz")
  wwf.RegisterRoutes(api)

  r.Run(":3000")
}
```

## `Zip` methods

- `api/zip` -> ZipFilesHendler
- `api/unzip` -> UnZipHendler
