package main

type Config struct {
	PublicURL        string `yaml:"public_url" validate:"required"`
	ServerPort       int    `yaml:"server_port" validate:"required"`
	ElasticsearchURL string `yaml:"elasticsearch_url" validate:"required"`
}
