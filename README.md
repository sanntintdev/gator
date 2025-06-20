# Gator

Gator is a command-line RSS feed aggregator built with Go. It helps you manage RSS feeds, follow your favorite content sources, and browse posts all from your terminal.

## What it does

This application allows you to:
- Manage user accounts
- Add and follow RSS feeds
- Browse posts from your followed feeds
- Track feed updates automatically

## Requirements

- Go 1.24.2 or higher
- PostgreSQL database
- Internet connection for fetching RSS feeds

## Installation

1. Clone this repository:
```bash
git clone https://github.com/sanntintdev/gator.git
cd gator
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up your PostgreSQL database and update the configuration file with your database URL.

4. Run the database migrations:
```bash
# The SQL schema files are in sql/schema/ directory
```

## Usage

The application uses a command-line interface. Here are the main commands:

### User Management
- Create and manage user accounts
- Set up your profile for feed management

### Feed Management  
- Add RSS feeds to the system
- Follow feeds you want to track
- Unfollow feeds you no longer need

### Browse Posts
- View posts from your followed feeds
- Browse through recent content
- Stay updated with your favorite sources

## Configuration

The application uses a configuration file to store database connection details and other settings. Make sure to set up your database URL correctly.

## Project Structure

```
gator/
   main.go              # Application entry point
   internal/
      commands/        # Command handlers and logic
      config/          # Configuration management
      database/        # Database operations and models
   sql/
      queries/         # SQL query files
      schema/          # Database schema migrations
   sqlc.yaml           # SQLC configuration
```

## Technologies Used

- **Go**: Main programming language
- **PostgreSQL**: Database for storing feeds, users, and posts
- **SQLC**: SQL code generation tool
- **UUID**: For unique identifiers

## Contributing

This project is part of a learning exercise. Feel free to explore the code and suggest improvements.

## License

This project is for educational purposes.