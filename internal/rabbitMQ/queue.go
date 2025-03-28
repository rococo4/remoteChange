package rabbitMQ

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

// Публикация конфигурации в очередь
func publishConfig(config string) error {
	// Подключение к RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("не удалось открыть канал: %w", err)
	}
	defer ch.Close()

	// Объявление очереди
	q, err := ch.QueueDeclare(
		"config_deploy_queue", // имя очереди
		true,                  // сохранять в памяти RabbitMQ
		false,                 // не удалять при отсутствии потребителей
		false,                 // доступна другим коннектам
		false,                 // без дополнительных флагов
		nil,
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	// Публикация сообщения
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (имя очереди)
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(config),
		},
	)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	log.Printf("Конфиг отправлен: %s", config)
	return nil
}

func consumeConfigs() error {
	// Подключение к RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("не удалось открыть канал: %w", err)
	}
	defer ch.Close()

	// Объявление очереди (должно совпадать с продюсером)
	q, err := ch.QueueDeclare(
		"config_deploy_queue",
		true,  // сохранять в памяти RabbitMQ
		false, // не удалять при отсутствии потребителей
		false, // доступна другим коннектам
		false, // без дополнительных флагов
		nil,
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	// Получение сообщений
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack (можно выставить false для явного подтверждения)
		false, // не эксклюзивный
		false, // без локального режима
		false, // без ожидания
		nil,
	)
	if err != nil {
		return fmt.Errorf("не удалось начать чтение из очереди: %w", err)
	}

	log.Println("✅ Ожидание новых конфигов...")

	// Читаем сообщения из канала
	for msg := range msgs {
		log.Printf("Получен конфиг: %s", msg.Body)
		// Здесь можно добавить логику деплоя в Kubernetes
	}

	return nil
}
