# 🐇 RabbitMQ Dead Letter Queue with Per-Message TTL (Golang)

This project demonstrates how to implement **message-level TTL** with a **Dead Letter Queue (DLQ)** using RabbitMQ and Go (`amqp091-go`).

---

## 🧰 Requirements

```bash
Go 1.18+
RabbitMQ (local or Docker)
```

Install Go RabbitMQ client:

```bash
go get github.com/rabbitmq/amqp091-go
```

---

## 🐳 RabbitMQ in Docker

```bash
docker run -d --hostname rabbit --name rabbitmq   -p 5672:5672 -p 15672:15672   rabbitmq:3-management
```

Access UI: [http://localhost:15672](http://localhost:15672)  
Login: `guest` / `guest`

---

## 📁 Project Structure

<pre>
.
├── producer.go         # Publishes a message with TTL
├── consumer.go         # Listens only to the Dead Letter Queue (DLQ)
└── README.md           # You're here
</pre>

---

## 🔧 Key RabbitMQ Configuration

<table>
<thead>
<tr><th>Parameter</th><th>Description</th></tr>
</thead>
<tbody>
<tr><td><code>x-dead-letter-exchange</code></td><td>The DLX where dead messages are sent</td></tr>
<tr><td><code>x-dead-letter-routing-key</code></td><td>The routing key used when forwarding to the DLQ</td></tr>
<tr><td><code>Expiration</code></td><td>Set per-message TTL (e.g., "5000" = 5 seconds)</td></tr>
</tbody>
</table>

---

## 🚀 How to Run

<details>
<summary><strong>Step 1: Start DLQ Consumer</strong></summary>

```bash
go run consumer.go
```

This sets up:

- DLX: `dead_letter_exchange`
- DLQ: `dead_letter_queue`
- Main queue: `main_queue` (with DLX configured)

It only listens to the DLQ (not the main queue).
</details>

<details>
<summary><strong>Step 2: Run Producer</strong></summary>

```bash
go run producer.go
```

This sends a message with a 5-second TTL to `main_queue`.  
Since no one consumes `main_queue`, the message will expire and go to the DLQ.
</details>

<details>
<summary><strong>Step 3: Observe Expired Message</strong></summary>

After ~5 seconds, the consumer prints:

```
[DLQ] Received expired message: TTL test message
```
</details>

---

## 💡 Why It Works

- Messages published with an **Expiration** are held in `main_queue`.
- No consumer reads `main_queue`, so they expire.
- RabbitMQ forwards them to the `dead_letter_queue` using the DLX config.

---

## 🧠 Behavior Matrix

| Action                        | TTL Triggered | Sent to DLQ |
|------------------------------|---------------|-------------|
| Message not consumed         | ✅ Yes        | ✅ Yes      |
| Message immediately consumed | ❌ No         | ❌ No       |
| Message rejected (no requeue)| ❌ N/A        | ✅ Yes      |
| Message acked                | ❌ No         | ❌ No       |

---

## 📚 References

- <a href="https://www.rabbitmq.com/dlx.html" target="_blank">RabbitMQ Dead Letter Exchange</a>
- <a href="https://pkg.go.dev/github.com/rabbitmq/amqp091-go" target="_blank">amqp091-go library</a>

---

## 📜 License

MIT – Use freely for personal or commercial projects.

---