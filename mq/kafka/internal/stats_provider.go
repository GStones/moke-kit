package internal

import (
	k "github.com/segmentio/kafka-go"
)

type StatsProvider interface {
	Stats() (string, interface{})
}

type ReaderStatsProvider struct {
	reader *k.Reader
}

func (s *ReaderStatsProvider) Stats() (string, interface{}) {
	stats := s.reader.Stats()
	msg := stats.ClientID

	return msg, stats
}

type WriterStatsProvider struct {
	writer *k.Writer
}

func (s *WriterStatsProvider) Stats() (string, interface{}) {
	stats := s.writer.Stats()
	msg := stats.ClientID

	return msg, stats
}
