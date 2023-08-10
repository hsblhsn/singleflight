# Singleflight

Package singleflight provides a duplicate function call suppression mechanism.
It's a generic supported wrapper of [`golang.org/x/sync/singleflight`](https://pkg.go.dev/golang.org/x/sync/singleflight).

### Install

```bash
go get -u github.com/hsblhsn/singleflight
```

### Usage

For detailed usage please refer to the [`golang.org/x/sync/singleflight`](https://pkg.go.dev/golang.org/x/sync/singleflight) package documentation.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hsblhsn/singleflight"
)

func main() {
	g := singleflight.NewGroup[int]()
	inputs := []int{0, 1, 1, 1, 2, 2, 2, 3, 3, 3}
	for _, v := range inputs {
		go func(input int) {
			key := fmt.Sprintf("key-%d", input)
			result, _, _ := g.Do(key, func() (int, error) {
				// simulate a slow function
				time.Sleep(time.Second)
				log.Println("[SLOW_FUNCTION_CALL]", input)
				return input * input, nil
			})
			log.Println("[OUTPUT]", input, result)
		}(v)
	}
	// wait for all goroutines to finish.
	// in real world this should be done by a sync.WaitGroup or similar.
	time.Sleep(2 * time.Second)
}
```
[Playground](https://go.dev/play/p/mNE_wvKHmWc)

```txt
2009/11/10 23:00:01 [SLOW_FUNCTION_CALL] 3
2009/11/10 23:00:01 [OUTPUT] 3 9
2009/11/10 23:00:01 [SLOW_FUNCTION_CALL] 1
2009/11/10 23:00:01 [OUTPUT] 1 1
2009/11/10 23:00:01 [OUTPUT] 1 1
2009/11/10 23:00:01 [OUTPUT] 3 9
2009/11/10 23:00:01 [OUTPUT] 3 9
2009/11/10 23:00:01 [SLOW_FUNCTION_CALL] 0
2009/11/10 23:00:01 [OUTPUT] 0 0
2009/11/10 23:00:01 [SLOW_FUNCTION_CALL] 2
2009/11/10 23:00:01 [OUTPUT] 2 4
2009/11/10 23:00:01 [OUTPUT] 2 4
2009/11/10 23:00:01 [OUTPUT] 2 4
2009/11/10 23:00:01 [OUTPUT] 1 1
```

So you can see that the slow function is called only once for each unique input.
But the output is returned to all callers. That's how singleflight works.


### License

No license. Do whatever you want with this code. I am not responsible for anything.