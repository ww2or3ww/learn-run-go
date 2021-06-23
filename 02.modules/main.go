package main

import (
	"fmt"
	"time"

	"github.com/Songmu/go-httpdate"
)

func main() {
	strTime := httpdate.Time2Str(time.Now())
	fmt.Println(strTime)
}
