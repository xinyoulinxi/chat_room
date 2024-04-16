package utils

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

var snowId *snowflake.Node

func init() {
	id, err := snowflake.NewNode(rand.New(rand.NewSource(time.Now().UnixMilli())).Int63n(1024))
	if err != nil {
		panic(err)
	}
	snowId = id
}

func GenerateId() string {
	return snowId.Generate().String()
}
