# gocache

Go library that provide abstract cache support.

The library currently supports two type of cache:

- Redis.
- In Memory.

## Example usage

### Using sugar Getxxx methods

```go
package main

import (
	"fmt"
	"github.com/creekorful/gocache"
	"log"
	"time"
)

func main() {
	c := gocache.NewMemoryCache("test")

	totalUsers, err := c.GetInt64("totalUsers", func() (int64, time.Duration) {
		// Expensive operation to count number of users...
		totalUsers := int64(1000)
		return totalUsers, 5 * time.Minute
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(totalUsers)
}
```

### Using classic Get / Set

```go
package main

import (
	"fmt"
	"github.com/creekorful/gocache"
	"log"
	"time"
)

func main() {
	c := gocache.NewMemoryCache("test")

	totalUsers, exists, err := c.Int64("totalUsers")
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		// Expensive operation to count number of users...
		totalUsers = int64(1000)

		if err := c.SetInt64("totalUsers", totalUsers, 5*time.Minute); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(totalUsers)
}
```

## Who's using the library?

- [darkspot-org/bathyscaphe](https://github.com/darkspot-org/bathyscaphe) : Fast, highly configurable, cloud native dark
  web crawler.