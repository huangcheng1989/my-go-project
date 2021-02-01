package main

import (
	"awesomeProject/dbredis"
	"context"
	"fmt"
)

func main() {
	dbredis.InitRedisCluster("prod")
	count, err := dbredis.GetRedisCluster().ZCard(context.Background(), "popular_talent_rank").Result()
	if err != nil {
		fmt.Println(0)
	}
	fmt.Println(count)
}
