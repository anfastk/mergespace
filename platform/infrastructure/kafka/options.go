package kafka

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func ProducerOpts(brokers []string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(10),
		kgo.ProducerLinger(10 * time.Millisecond),
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),
	}
}

func ConsumerOpts(brokers []string, group string, topics []string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
		kgo.BlockRebalanceOnPoll(),
	}
}
