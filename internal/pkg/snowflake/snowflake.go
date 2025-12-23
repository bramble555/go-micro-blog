package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

func NewNode(machineID int64) (*snowflake.Node, error) {
	// 统一 Epoch，便于论文和维护
	snowflake.Epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6
	return snowflake.NewNode(machineID)
}
