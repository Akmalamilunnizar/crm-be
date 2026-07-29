# 📊 CRM Project Architecture & Functional Summary

## 1. Executive Overview
This project is a comprehensive and highly modular Customer Relationship Management (CRM) backend system tailored for an Internet Service Provider (ISP) or telecommunications business. It is designed to manage the entire customer lifecycle, from initial installation and asset allocation to recurring billing and ongoing technical support. The system is built for high performance and scalability, featuring automated notification systems and a structured approach to managing field technician reports.

## 2. Architecture & Tech Stack
The backend is built with modern, performant technologies utilizing a clean, layered architectural pattern:

*   **Core Language:** Go (Golang 1.22+) ensuring high concurrency and fast execution.
*   **Web Framework:** Go Fiber v2, chosen for its extreme speed and low memory footprint in handling HTTP APIs.
*   **Database & ORM:** MySQL handled via GORM, allowing for complex relational data mapping, soft deletes, and seamless migrations.
*   **Microservices:** A separate Node.js service running `whatsapp-web.js` to manage session-based WhatsApp web automation securely outside the core Go application.
*   **Design Pattern:** The project follows standard Go project layouts (`cmd/`, `internal/`) and implements a layered architecture. The `internal/` directory strictly separates concerns into:
    *   `api/` & `routes/`: Endpoint definitions.
    *   `handlers/`: Request/response formatting.
    *   `services/`: Core business logic.
    *   `models/`: Data representation, further divided into `entities` (database schemas), `dto` (data transfer objects), and `repository` (database interaction layer).
    *   `middleware/`: Authentication (JWT) and request validation.

## 3. Key Modules & Capabilities
The system is divided into several powerful domain-driven modules:

*   **Asset & Inventory Management System:**
    *   Tracks physical inventory, network devices, and infrastructure items.
    *   Manages "Asset-Company" relationships to easily monitor which equipment belongs to which client or branch.
*   **Customer & Installation Management:**
    *   Maintains detailed customer profiles.
    *   Supports mapping multiple products or service plans to a single customer.
    *   Tracks the lifecycle of customer network installations.
*   **Installation Reports & Technician Tracking:**
    *   Allows field technicians to submit comprehensive installation reports.
    *   Handles multipart file uploads for site documentation and photo evidence.
    *   Assigns and tracks specific technician teams to installation tasks.
*   **Billing & Invoicing Engine:**
    *   Automates the generation of financial invoices for services rendered.
    *   Features a robust recurring invoice system for monthly subscription billing.
    *   Includes capabilities for generating PDF versions of invoices for customers.
*   **Ticketing & Support System:**
    *   Features a streamlined ticket classification system to handle customer complaints, network outages, and maintenance requests effectively.

## 4. External Integrations
To proactively engage with customers and automate internal alerts, the CRM integrates seamlessly with external messaging platforms:

*   **WhatsApp Broadcasting:** Utilizes the custom Node.js microservice (`whatsapp-service`) to send automated billing reminders, installation updates, and mass broadcasts directly to customers' WhatsApp numbers.
*   **Telegram Bot Integration:** Includes an `internal/telegram` module, likely used to push immediate alerts, error logs, or dispatch notifications directly to internal technician and admin groups.

## 5. Data Management & Machine Learning
*   **Database Migrations:** The database schema is strictly version-controlled using SQL migrations (located in `migrations/`) and Go-based seeder scripts, ensuring the database state is easily reproducible across environments.
*   **Advanced Views:** The database utilizes complex SQL Views to easily aggregate data across installations, assets, and customers for the frontend to consume rapidly.
*   **Machine Learning (ML):** The presence of an `internal/ml` package indicates forward-looking capabilities, potentially for predicting customer churn, bandwidth utilization forecasting, or automated ticket classification.
