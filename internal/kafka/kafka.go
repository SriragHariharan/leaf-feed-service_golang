package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"

	kgo "github.com/segmentio/kafka-go"
)

const envKafkaBrokers = "KAFKA_BROKERS"

var (
	Brokers []string
)

// Connect parses KAFKA_BROKERS (comma-separated host:port list), verifies connectivity,
// and stores the broker list for consumers.
func Connect() error {
	raw := strings.TrimSpace(os.Getenv(envKafkaBrokers))
	if raw == "" {
		raw = "localhost:9092"
	}

	Brokers = splitBrokers(raw)
	if len(Brokers) == 0 {
		return fmt.Errorf("%s must contain at least one host:port", envKafkaBrokers)
	}

	conn, err := kgo.Dial("tcp", Brokers[0])
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	_ = conn.Close()

	log.Printf("kafka: connected (brokers=%v)", Brokers)
	return nil
}

func splitBrokers(s string) []string {
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
