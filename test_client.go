// test_client.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Notification struct {
	Type    string    `json:"type"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Clients []string  `json:"clients,omitempty"`
}

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	for {
		fmt.Println("\n========================================")
		fmt.Println("🔔 NATS Notification Test Client")
		fmt.Println("========================================")
		fmt.Println("Connected to NATS server ✅")
		fmt.Println("\nChoose test:")
		fmt.Println("1. Broadcast notification")
		fmt.Println("2. Targeted notification")
		fmt.Println("3. Multiple notifications (priority test)")
		fmt.Println("4. List notifications")
		fmt.Println("5. Exit")

		var choice int
		fmt.Print("\nEnter choice (1-5): ")

		// More robust input handling
		n, err := fmt.Scanf("%d", &choice)
		if err != nil || n != 1 {
			fmt.Println("❌ Invalid input. Please enter a number between 1-5.")
			// Clear the input buffer
			var discard string
			fmt.Scanln(&discard)
			continue
		}

		switch choice {
		case 1:
			testBroadcast(nc)
		case 2:
			testTargeted(nc)
		case 3:
			testPriority(nc)
		case 4:
			testList(nc)
		case 5:
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please enter a number between 1-5.")
		}

		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}

func testBroadcast(nc *nats.Conn) {
	fmt.Println("\n📡 Testing Broadcast Notification...")

	notifications := []Notification{
		{
			Type:    "info",
			Message: "Hello World from Go client!",
			Time:    time.Now(),
		},
	}

	data, err := json.Marshal(notifications)
	if err != nil {
		fmt.Printf("❌ Failed to marshal notifications: %v\n", err)
		return
	}

	fmt.Println("📤 Sending broadcast notification...")
	resp, err := nc.Request("NOTIFICATION.send-to-all", data, time.Second*5)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}

	fmt.Printf("✅ Success! Response: %s\n", string(resp.Data))
	fmt.Println("💡 Check your WebSocket client to see the message!")
}

func testTargeted(nc *nats.Conn) {
	fmt.Println("\n🎯 Testing Targeted Notification...")

	notifications := []Notification{
		{
			Type:    "warning",
			Message: "Targeted message from Go client!",
			Time:    time.Now(),
			Clients: []string{"1", "2"},
		},
	}

	data, err := json.Marshal(notifications)
	if err != nil {
		fmt.Printf("❌ Failed to marshal notifications: %v\n", err)
		return
	}

	fmt.Println("📤 Sending targeted notification to clients 1 and 2...")
	resp, err := nc.Request("NOTIFICATION.send-to-clients", data, time.Second*5)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}

	fmt.Printf("✅ Success! Response: %s\n", string(resp.Data))
	fmt.Println("💡 Only WebSocket clients with ID 1 or 2 should receive this message!")
}

func testPriority(nc *nats.Conn) {
	fmt.Println("\n⚡ Testing Priority Ordering...")

	notifications := []Notification{
		{
			Type:    "info",
			Message: "Info message (priority 3)",
			Time:    time.Now(),
		},
		{
			Type:    "error",
			Message: "Error message (priority 1)",
			Time:    time.Now(),
		},
		{
			Type:    "warning",
			Message: "Warning message (priority 2)",
			Time:    time.Now(),
		},
	}

	data, err := json.Marshal(notifications)
	if err != nil {
		fmt.Printf("❌ Failed to marshal notifications: %v\n", err)
		return
	}

	fmt.Println("📤 Sending multiple notifications...")
	resp, err := nc.Request("NOTIFICATION.send-to-all", data, time.Second*5)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}

	fmt.Printf("✅ Success! Response: %s\n", string(resp.Data))
	fmt.Println("💡 WebSocket clients should receive messages in priority order:")
	fmt.Println("   1st: Error message (highest priority)")
	fmt.Println("   2nd: Warning message (medium priority)")
	fmt.Println("   3rd: Info message (lowest priority)")
}

func testList(nc *nats.Conn) {
	fmt.Println("\n📋 Testing Notification List...")

	fmt.Println("📤 Requesting stored notifications...")
	resp, err := nc.Request("NOTIFICATION.list", []byte("{}"), time.Second*5)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}

	fmt.Printf("✅ Success! Response: %s\n", string(resp.Data))
	fmt.Println("💡 This shows all notifications stored in the database!")
}
