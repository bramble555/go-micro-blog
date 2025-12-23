package global

import (
	"go-micro-blog/internal/pkg/snowflake"
	"log"

	sf "github.com/bwmarrin/snowflake"
)

var SnowflakeNode *sf.Node

func InitSnowflake(machineID int64) {
	node, err := snowflake.NewNode(machineID)
	if err != nil {
		log.Fatalf("init snowflake failed: %v", err)
	}
	SnowflakeNode = node
}

func GenID() int64 {
	return SnowflakeNode.Generate().Int64()
}
