// rabbitmq-service.go
package main

import (
    "sync"

    "github.com/streadway/amqp"
)

type RabbitMQService struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    url     string
    mu      sync.Mutex
}

// Initialize RabbitMQ connection and channel
func NewRabbitMQService(url string) (*RabbitMQService, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    ch.Qos(1, 0, false)

    return &RabbitMQService{
        conn:    conn,
        channel: ch,
        url:     url,
    }, nil
}

// Close the connection
func (r *RabbitMQService) Close() {
    r.channel.Close()
    r.conn.Close()
}

// Publish a message to a queue
func (r *RabbitMQService) Publish(queue string, message []byte) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    _, err := r.channel.QueueDeclare(
        queue,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    err = r.channel.Publish(
        "",
        queue,
        false,
        false,
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "application/json",
            Body:         message,
        },
    )
    return err
}

// Consume messages from a queue
func (r *RabbitMQService) Consume(queue string, callback func(amqp.Delivery)) error {
    msgs, err := r.channel.Consume(
        queue,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    for msg := range msgs {
        callback(msg)
    }
    return nil
}
