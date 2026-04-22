# Gator (Blog Aggregator)

Gator is a command-line RSS feed aggregator written in Go. It allows users to register, follow various RSS feeds, and browse collected posts from a PostgreSQL database.

## Prerequisites

Before installing, ensure you have the following installed on your system:

- **Go**: version 1.23.1 or higher (as specified in `go.mod`)

- **PostgreSQL**: A running instance to store users, feeds, and posts.

## Installation

To install the `gator` binary to your `$GOPATH/bin` folder, run:

Bash

```
go install github.com/Bravnar/gator@latest

```

_Note: Ensure your Go bin directory is in your system's PATH to run the command directly._

## Configuration

Gator expects a `.gatorconfig.json` file in your home directory to manage the database connection and current user.

Example `~/.gatorconfig.json`:

JSON

```
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "your_username"
}

```

## Commands

Usage: `gator <command> [arguments]`

| **Command** | **Arguments**         | **Description**                                           |
| ----------- | --------------------- | --------------------------------------------------------- |
| `register`  | `<name>`              | Create a new user account.                                |
| `login`     | `<name>`              | Switch the current active user.                           |
| `users`     |                       | List all registered users.                                |
| `addfeed`   | `<name> <url>`        | Add a new RSS feed to follow (requires login).            |
| `feeds`     |                       | Show all feeds and the users who added them.              |
| `follow`    | `<url>`               | Follow an existing feed (requires login).                 |
| `following` |                       | List feeds the current user follows.                      |
| `unfollow`  | `<url>`               | Stop following a feed.                                    |
| `agg`       | `<time_between_reqs>` | Start the aggregator (e.g., `1m`, `1h`).                  |
| `browse`    | `[limit]`             | View posts from followed feeds.                           |
| `reset`     |                       | **Danger:** Deletes all users and data from the database. |

## Dependencies

This project uses the following Go packages:

- `github.com/google/uuid`: For unique identifiers.

- `github.com/lib/pq`: PostgreSQL driver for Go.

---

**Repository:** <https://github.com/Bravnar/boot_dev_blog_aggregator>
