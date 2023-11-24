package main

import (
	"context"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/controller"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource/cache"
	database2 "gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource/database"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/pkg/database"
	"os"
	"os/signal"
	"time"
)

func main() {
	db := database.NewDatabase("db/db.json")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cashCliet := cache.NewClient(4*time.Second, 6*time.Second, database2.NewClient(db))

	client := controller.NewClient(cashCliet)

	var bestUser = "best user, expired after 30 seconds"

	// Создаём запись
	err := client.Set(ctx, "user:12345:profile", bestUser, 4*time.Second)
	if err != nil {
		panic(err)
	}

	// Получаем запись из кэша
	got, err := client.Get(ctx, "user:12345:profile")
	if err != nil {
		panic(err)
	}

	if got != bestUser {
		panic("invalid value")
	}

	select {
	case <-time.After(8 * time.Second):
	case <-ctx.Done():
	}

	// Получаем запись из базы данных и обновляем кэщ
	gotAgain, err := client.Get(ctx, "user:12345:profile")
	if err != nil {
		panic(err)
	}

	if gotAgain != bestUser {
		panic("invalid value")
	}
}
