package config

type Queue struct {
    RabbitMQ RabbitMQ `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
}

type RabbitMQ struct {
    Uri string `mapstructure:"uri" json:"uri" yaml:"uri"`
    ExchangeName string `mapstructure:"exchange_name" json:"exchange_name" yaml:"exchange_name"`
    DelayExchangeName string `mapstructure:"delay_exchange_name" json:"delay_exchange_name" yaml:"delay_exchange_name"`
}
