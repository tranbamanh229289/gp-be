package main

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())
	
	_, err = c.AddFunc("", func(){
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
})
}

