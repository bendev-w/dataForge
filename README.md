# DataForge - Database Seeding CLI Tool

A powerful Go CLI tool for seeding PostgreSQL and MySQL databases with realistic, lifelike data at specified intervals.

## Features

- **Interval-based Seeding**: Automatically seed your database at regular intervals
- **Complete Schema**: Pre-built tables for a typical e-commerce database
  - `customers` - Customer information
  - `products` - Product catalog
  - `orders` - Customer orders
  - `order_items` - Items in each order
  - `payments` - Payment records
- **Realistic Data**: Uses `gofakeit` to generate realistic fake data
- **Transaction Support**: All seeding operations are atomic
- **Easy CLI**: Simple, intuitive command-line interface
- **Hasura GraphQL Admin Console**: Built-in web UI to inspect and query seeded data in real-time

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher or MySQL 8.0 or higher (or use the included Docker Compose setup)
- Docker & Docker Compose (optional, for PostgreSQL and Hasura)
- A PostgreSQL or MySQL database to seed

## Installation

### Build from Source

```bash
go build -o dataforge.exe
```

## Usage

### Commands

#### Start Seeding

Continuously seed your database at specified intervals.

```bash
dataforge start --db "postgres://user:password@localhost:5432/dbname" --interval 5s
# MySQL
dataforge start --db "mysql://user:password@localhost:3306/dbname" --interval 5s
```

**Parameters:**
- `--db` (required): PostgreSQL or MySQL connection URL
- `--interval` (optional): Seeding interval (default: `5s`, format: `5s`, `1m`, `10m`, `1h`)

**Examples:**
```bash
# Seed every 5 seconds
dataforge start --db postgres://localhost:5432/shop --interval 5s

# Seed every 1 minute
dataforge start --db postgres://localhost:5432/shop --interval 1m

# Seed every 10 minutes
dataforge start --db postgres://localhost:5432/shop --interval 10m

# Seed MySQL every minute
dataforge start --db mysql://dataforge:dataforge@localhost:3306/shop --interval 1m
```

#### Clear Database

Remove all data from database tables (sequences will be reset to 1).

```bash
dataforge clear --db "postgres://user:password@localhost:5432/dbname"
# MySQL
dataforge clear --db "mysql://user:password@localhost:3306/dbname"
```

**Example:**
```bash
dataforge clear --db postgres://localhost:5432/shop
```

#### Refresh Database

Clear all existing data and seed with fresh new data (in one operation).

```bash
dataforge refresh --db "postgres://user:password@localhost:5432/dbname"
# MySQL
dataforge refresh --db "mysql://user:password@localhost:3306/dbname"
```

**Example:**
```bash
dataforge refresh --db postgres://localhost:5432/shop
```

## Hasura GraphQL Admin Console

View and query your seeded data in real-time with Hasura's powerful GraphQL interface.

### Accessing Hasura

1. **Start the services**:
   ```bash
   docker-compose up -d
   ```

2. **Open the console** in your browser:
   ```
   http://localhost:8080/console
   ```

3. **Login with**:
   - Admin Secret: `dataforge_admin_secret`

### Features in Hasura Console

- **Data Browser**: View all customers, products, orders, and payments in a table format
- **GraphQL Explorer**: Write and execute GraphQL queries
- **Relationships**: Automatically discovered relationships between tables (e.g., orders linked to customers)
- **Mutations**: Insert, update, or delete data directly
- **Real-time Subscriptions**: Watch data changes as they happen
- **API Explorer**: Browse auto-generated GraphQL APIs

### Example GraphQL Query

```graphql
query GetOrderDetails {
  orders(limit: 10, order_by: {order_date: desc}) {
    id
    total_amount
    status
    customer {
      first_name
      last_name
      email
    }
    order_items {
      product {
        name
        price
      }
      quantity
      subtotal
    }
    payments {
      amount
      payment_method
      status
    }
  }
}
```

## Database Schema

### customers
- `id` (Serial Primary Key)
- `first_name`, `last_name` - Customer name
- `email` - Unique email address
- `phone` - Contact number
- `address`, `city`, `state`, `zip_code`, `country` - Address details
- `created_at`, `updated_at` - Timestamps

### products
- `id` (Serial Primary Key)
- `name` - Product name
- `description` - Product description
- `price` - Product price
- `stock_quantity` - Available stock
- `category` - Product category
- `sku` - Unique stock-keeping unit
- `created_at`, `updated_at` - Timestamps

### orders
- `id` (Serial Primary Key)
- `customer_id` - References customers
- `order_date` - When the order was placed
- `total_amount` - Order total
- `status` - Order status (pending, processing, shipped, delivered, cancelled)
- `shipping_address` - Delivery address
- `created_at`, `updated_at` - Timestamps

### order_items
- `id` (Serial Primary Key)
- `order_id` - References orders
- `product_id` - References products
- `quantity` - Number of items
- `unit_price` - Price per item
- `subtotal` - Total for this line item
- `created_at` - Creation timestamp

### payments
- `id` (Serial Primary Key)
- `order_id` - References orders
- `amount` - Payment amount
- `payment_method` - Method used (credit_card, debit_card, paypal, etc.)
- `status` - Payment status (pending, completed, failed, refunded)
- `transaction_id` - Unique transaction ID
- `created_at`, `updated_at` - Timestamps

## Data Generation

The tool generates:
- **50 customers** per seeding cycle with realistic names, emails, and addresses
- **100 products** per cycle with various categories and prices
- **5 orders per customer** with dates in the past 90 days
- **2-5 order items per order** with quantities and pricing
- **1 payment per order** with various payment methods and statuses

All data is generated using the `gofakeit` library to ensure variety and realism.

## Development

### Project Structure

```
DataForge/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # Root command setup
│   └── start.go           # Start command implementation
├── db/
│   ├── connection.go      # Database connection
│   └── schema.go          # Schema initialization
├── seeders/
│   ├── seeders.go         # Main seeding orchestrator
│   ├── customer_seeder.go # Customer data seeding
│   ├── product_seeder.go  # Product data seeding
│   ├── order_seeder.go    # Order data seeding
│   ├── order_item_seeder.go # Order items seeding
│   └── payment_seeder.go  # Payment data seeding
├── go.mod                 # Go module definition
└── README.md              # This file
```

### Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/brianvoe/gofakeit/v6` - Fake data generation

## License

MIT License
