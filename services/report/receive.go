package main

import (
    "fmt"
    "github.com/streadway/amqp"
    "log"
    "os"
    "encoding/json"
)

type Product struct {
    Name  string  `json:"name"`
    Value string `json:"value"`
}

type CreditCard struct {
	number, securityNumber string
}

type Address struct {
	zipCode, street, number, neighborhood, city, state string
}

type Order struct {
    Name       string     `json:"name"`
    Email      string     `json:"email"`
    CPF        string     `json:"cpf"`
    CreditCard CreditCard `json:"creditCard"`
    Products   []Product  `json:"products"`
    Address    Address    `json:"address"`
}

func updateReports(report map[string]int, products []Product) map[string]int {
	for _, product := range(products) {
		if(len(product.Name) < 0) {
            continue
        }else if _, ok := report[product.Name]; !ok {
            report[product.Name] = 1;
        }else {
            report[product.Name]++;
        }
	}
	return report
}

func printReport(report map[string]int) {
	for key, value := range(report) {
		log.Printf("%s = %d vendas\n", key, value)
	}
}

func processMessage(msg amqp.Delivery) {
	log.Println("Pedido recebido com sucesso!")
	var order Order
    err := json.Unmarshal(msg.Body, &order)
    if err != nil {
        log.Printf("Failed to parse message: %v", err)
        return
    }

	report := make(map[string]int)
	report = updateReports(report, order.Products)
	printReport(report)
}

func main() {
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
        os.Getenv("RABBITMQ_LOGIN"),
        os.Getenv("RABBITMQ_PASSWORD"),
        os.Getenv("RABBITMQ_HOST"),
        os.Getenv("RABBITMQ_PORT"),
        os.Getenv("RABBITMQ_VHOST"),
    )

    rabbitService, err := NewRabbitMQService(rabbitURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer rabbitService.Close()

    queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
    log.Printf("Subscribed successfully to queue: %s", queueName)

    err = rabbitService.Consume(queueName, processMessage)
    if err != nil {
        log.Fatalf("Failed to consume messages: %v", err)
    }
}