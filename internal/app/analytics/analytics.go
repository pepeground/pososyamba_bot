package analytics

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"os"
	"strconv"
	"time"
)

func SendToInflux(username string, userID int, chatID int64, chatTitle, messageType, command string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("INFLUX_URL"),
		Username: os.Getenv("INFLUX_USERNAME"),
		Password: os.Getenv("INFLUX_PASSWORD"),
	})
	if err != nil {
		log.Printf("%+v\n", "Error creating InfluxDB Client: "+err.Error())
	}
	defer c.Close()

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "web_services",
	})

	// Create a point and add to batch
	tags := map[string]string{
		"command":  command,
		"chat_id":  strconv.Itoa(int(chatID)),
		"user_id":  strconv.Itoa(int(userID)),
		"username": username,
	}

	fields := map[string]interface{}{
		"username":    username,
		"user_id":     userID,
		"chat_title":  chatTitle,
		"chat_id":     chatID,
		"messageType": messageType,
	}

	pt, err := client.NewPoint("pososyamba_usage", tags, fields, time.Now())
	if err != nil {
		log.Println("Error: ", err.Error())
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		bp.AddPoint(pt)

		// Write the batch
		c.Write(bp)
	}
}
