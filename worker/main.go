package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	gomail "gopkg.in/mail.v2"
)

type Car struct {
	Model string `json:"model"`
	Year  int    `json:"year"`
	Color string `json:"color"`
	Email string `json:"email"`
}

func main() {
	var ctx = context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println("Redis client connected successfully...")

	subscriber := redisClient.Subscribe(ctx, "send-car-data")
	fmt.Println("Created subscriber")

	car := Car{}

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	ticker := time.NewTicker(time.Duration(20) * time.Second)
	stop := make(chan bool)

	go func() {
		defer func() { stop <- true }()
		for {
			select {
			case <-ticker.C:
				message, err := subscriber.ReceiveMessage(ctx)
				if err != nil {
					fmt.Println("Error receive message")
					return
				}

				err = json.Unmarshal([]byte(message.Payload), &car)
				if err != nil {
					return
				}

				SendEmail(car)
			case <-stop:
				fmt.Println("closing goroutine")
				return
			}
		}
	}()

	// block until receive signal
	<-c
	ticker.Stop()

	stop <- true

	<-stop
	fmt.Println("Stop application")
}

func SendEmail(c Car) {

	m := gomail.NewMessage()
	m.SetHeader("From", "przmld033@gmail.com")
	m.SetHeader("To", c.Email)
	m.SetHeader("Subject", "Gomail test subject")

	body := "<h2 style=\"color: red\">Your Order :</h2><h3>Car model:" + c.Model + "</h3>" + "<h3>Car year:" + strconv.Itoa(c.Year) + "</h3>" + "<h3>Car color:" + c.Color + "</h3>"
	m.SetBody("text/html", body)

	fmt.Println("Message created")

	d := gomail.NewDialer("smtp.gmail.com", 587, "przmld033@gmail.com", "SECRET")
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	fmt.Println("Connecting to SMTP")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email successfully sent")
}
