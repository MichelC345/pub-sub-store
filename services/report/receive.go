package main

import (
    //"bufio"
    "fmt"
    //"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
    //"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
    "log"
    "os"
)

type CreditCard struct {
	number, securityNumber string
}

type Product struct {
	name string
	value float64
}

type Address struct {
	zipCode, street, number, neighborhood, city, state string
}

type Order struct {
	name, email, cpf string
	creditCard CreditCard
	products []Product
	address Address
}

func updateReports(report map[string]int, products []Product) map[string]int {
	for _, product := range(products) {
		if(len(product.name) < 0) {
            continue
        }else if _, ok := report[product.name]; !ok {
            report[product.name] = 1;
        }else {
            report[product.name]++;
        }
	}
	return report
}

func printReport(report map[string]int) {
	for key, value := range(report) {
		log.Printf("%s = %d vendas\n", key, value)
	}
}

func processMessage(msg Order) {
	log.Println("Pedido recebido com sucesso!")
	/*msgJson, err := json.MarshalIndent(mailData, "", "  ")
    if err != nil {
        log.Printf("Failed to convert to JSON: %v", err)
        return
    }
	log.Println(string(msgJson))*/

	report := make(map[string]int)
	report = updateReports(report, msg.products)
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