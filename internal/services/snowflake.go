package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/everywan/identifier"
)

/*
*
* 1                                               42           52             64
* +-----------------------------------------------+------------+---------------+
* | timestamp(ms)                                 | workerid   | sequence      |
* +-----------------------------------------------+------------+---------------+
* | 0000000000 0000000000 0000000000 0000000000 0 | 0000000000 | 0000000000 00 |
* +-----------------------------------------------+------------+---------------+
*
* 1. 41位时间截(秒级), 大概可用 ((1<<40)-1555868387)/60/60/24/365 == 34815.9 年(20190422), 后续改为毫秒也可用34年
* 2. 10位数据机器位，可以部署在1024个节点
* 3. 12位序列，毫秒内的计数，同一机器，同一时间截并发4096个序号
 */

const (
	workeridBits = uint(10) //机器id所占的位数
	sequenceBits = uint(12) //序列所占的位数
	workeridMax  = 1<<workeridBits - 1
	sequenceMax  = 1<<sequenceBits - 1
)

// SnowflakeService is
type SnowflakeService struct {
	sync.Mutex
	sequence  int64
	workerID  int64
	timestamp int64
}

// NewSnowflakeService 创建服务
func NewSnowflakeService(workerID int64) (snowflakeSvc *SnowflakeService, err error) {
	snowflakeSvc = new(SnowflakeService)
	if workerID < 0 || workerID > workeridMax {
		return snowflakeSvc, fmt.Errorf("workerid must be between 0 and %d", workeridMax)
	}
	snowflakeSvc.workerID = workerID
	return snowflakeSvc, nil
}

var _ identifier.SnowflakeService = &SnowflakeService{}

// Generate 生成唯一id
func (s *SnowflakeService) Generate(ctx context.Context) (id int64, err error) {
	s.Lock()
	defer s.Unlock()
	now := time.Now().Unix()

	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & sequenceMax
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().Unix()
			}
		}
	} else if now > s.timestamp {
		s.sequence = 0
	} else {
		return 0, fmt.Errorf("Clock moved backwards,  Refusing to generate id for %d milliseconds", s.timestamp-now)
	}
	s.timestamp = now

	r := int64(s.timestamp<<(workeridBits+sequenceBits) | s.workerID<<sequenceBits | s.sequence)
	return r, nil
}
