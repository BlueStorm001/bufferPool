# bufferPool


``` golang
import (
    "fmt"
    "github.com/BlueStorm001/bufferPool"
)
```

``` golang
func main() {
    buff :=  bufferPool.NewDefault()
    b := buff.Get()
    //.......
    b.Write([]byte("AAA")).Write([]byte("BBB"))
    //.......
    fmt.Println(b.Bytes())
     //.......
    buff.Put(b)
}
```