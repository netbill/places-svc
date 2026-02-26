package messenger

import (
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/pkg/log"
)

type ProducerConfig struct {
	Producer string   `json:"producer"`
	Brokers  []string `json:"brokers"`

	PlacesV1       ProduceKafkaConfig `json:"places_v1"`
	PlaceClassesV1 ProduceKafkaConfig `json:"place_classes_v1"`
}

type ProduceKafkaConfig struct {
	RequiredAcks string        `json:"required_acks"`
	Compression  string        `json:"compression"`
	Balancer     string        `json:"balancer"`
	BatchSize    int           `json:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout"`
}

func NewProducer(log *log.Logger, cfg ProducerConfig) *eventbox.Producer {
	producer := eventbox.NewProducer(log, cfg.Brokers...)

	err := producer.AddWriter(evtypes.PlacesTopicV1, eventbox.WriterTopicConfig{
		RequiredAcks: cfg.PlacesV1.RequiredAcks,
		Compression:  cfg.PlacesV1.Compression,
		Balancer:     cfg.PlacesV1.Balancer,
		BatchSize:    cfg.PlacesV1.BatchSize,
		BatchTimeout: cfg.PlacesV1.BatchTimeout,
	})
	if err != nil {
		panic(err)
	}

	err = producer.AddWriter(evtypes.PlaceClassesTopicV1, eventbox.WriterTopicConfig{
		RequiredAcks: cfg.PlaceClassesV1.RequiredAcks,
		Compression:  cfg.PlaceClassesV1.Compression,
		Balancer:     cfg.PlaceClassesV1.Balancer,
		BatchSize:    cfg.PlaceClassesV1.BatchSize,
		BatchTimeout: cfg.PlaceClassesV1.BatchTimeout,
	})
	if err != nil {
		panic(err)
	}

	return producer
}
