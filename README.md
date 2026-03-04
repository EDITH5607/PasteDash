# PasteDash

**PasteDash** is a high-performance, security-focused snippet sharing platform built with **Go**. It emphasizes production-grade architecture, leveraging the Go standard library to provide a fast and reliable environment for managing and sharing code snippets or text.

The project serves as a deep dive into scalable backend development, focusing on robust middleware, session management, and modern security practices.

---

### 🚀 Features

* **Secure Snippet Sharing:** Create and share snippets with built-in CSRF protection and secure HTTP headers.
* **Dynamic Session Management:** Efficient handling of user sessions to ensure data persistence and security.
* **Custom Middleware:** A powerful middleware stack for logging, recovery, and security enforcement.
* **Local HTTPS Development:** Fully implemented TLS support using self-signed certificates for an encrypted development workflow.
* **Optimized Routing:** Designed for high-throughput and low-latency performance.
* **Template Caching:** Intelligent server-side rendering with optimized template management for faster response times.

### 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Environment:** Ubuntu Linux
* **Editor:** LazyVim / Terminal-based workflow
* **Core Concepts:** * `net/http` standard library
* Secure Session Management
* TLS/HTTPS Implementation
* RESTful System Design



---

### 🏗️ Project Structure

The architecture follows a clean, "production-grade" pattern inspired by industry-standard Go workflows:

* **`cmd/web/`**: Entry point of the application, containing the main server logic and configuration.
* **`internal/`**: Core business logic, data models, and database interactions to prevent external package leakage.
* **`ui/`**: HTML templates and static assets (CSS/JS).

---

### 🛡️ Security Implementations

* **TLS Handshake:** Configured with specific cipher suites to ensure modern encryption standards.
* **Security Headers:** Implemented headers like `X-Content-Type-Options`, `X-Frame-Options`, and `HSTS`.
* **Input Validation:** Strict sanitization of user-submitted snippets.

---

### 🏁 Getting Started

#### Prerequisites

* Go 1.22+
* OpenSSL (for certificate generation)

#### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/EDITH5607/PasteDash
cd Pastedash

```


2. **Generate TLS certificates:**
```bash
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost

```


3. **Run the application:**
```bash
go run ./cmd/web

```



### 🚀 Live Demo

You can find the live version of **PasteDash** hosted and running at the link below. This deployment showcases the full production-grade workflow, including secure session handling and optimized template rendering.

**Live Link:** coming soon

---



