package main

type Config struct {
	DBUser           string `yaml:"db_user" validate:"required"`
	DBPassword       string `yaml:"db_password" validate:"required"`
	DBHost           string `yaml:"db_host" validate:"required"`
	DBPort           string `yaml:"db_port" validate:"required,numeric"`
	DBName           string `yaml:"db_name" validate:"required"`
	ElasticsearchURL string `yaml:"elasticsearch_url" validate:"required"`
}
