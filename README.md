<div align="center">

# 🐊 GATOR 🐊

</div>

## 🔌 Requirements 🔌

- Postgres
- Go

## ⚙️ Installation ⚙️

### 🐘 Postgres

1. Ensure PostgreSQL is installed and running on your system. You can refer to the [official PostgreSQL documentation](https://www.postgresql.org/docs/) for installation and setup instructions. Once installed, connect to PostgreSQL. We use the default user, 'postgres' to connect to it, but if you are using a different user, replace 'postgres' with the user's name.
```terminal
sudo -u postgres psql
```
2. Run the following SQL commands in the PostgreSQL shell to create the database and its tables:
```postgres
CREATE DATABASE gator;
```
- Connect to the new db

```
\c gator
```

- Migrate through all of the following:

```
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE feeds
(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
);

ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;

CREATE TABLE posts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_url UNIQUE (url)
);
```
 ***You can also find these migrations in my sql/schema directory with goose metadata in the comments if you'd like to control the schema with that.***

### 🐊 Gator

1. Install gator

```terminal
    go install github.com/ChipsAhoyEnjoyer/gator@latest
```

2. Verify the installation:
```terminal
gator version
```
This should display the installed version of Gator.

3. Create a file name ".gatorconfig.json" in your home directory (~)

4. Add the following to the gator config file:
```json
    {
        "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
        // Replace "username:password" with your actual PostgreSQL username and password.
        // For example: "postgres:mysecurepassword"
        "current_user_name": ""
    }
```
- *Make sure to replace "username:password" with your postgress username and password*

<!-- ...existing code... -->

## 🚀 Usage 🚀

Gator CLI provides the following commands:

- **register** - Create a new user
  ```terminal
  gator register '<username>'
  ```

- **login** - Log in as an existing user
  ```terminal
  gator login '<username>'
  ```

- **addfeed** - Add a new RSS feed to the database
  ```terminal
  gator addfeed <name> <url>
  ```

- **follow** - Follow an existing feed
  ```terminal
  gator follow '<link>'
  ```

- **unfollow** - Unfollow a feed you're currently following
  ```terminal
  gator unfollow '<link>'
  ```

- **browse** - Browse posts from feeds you follow (defaults to 2 posts)
  ```terminal
  gator browse [limit]
  ```

- **feeds** - List all available feeds in the database
  ```terminal
  gator feeds
  ```

- **following** - List all feeds you are currently following
  ```terminal
  gator following
  ```

- **users** - List all registered users
  ```terminal
  gator users
  ```

- **agg** - Aggregate/scrape feeds on a schedule. Refresh-rate can be '1s', '1m', '1h', etc.
  ```terminal
  gator agg '<refresh-rate>'
  ```

- **reset** - Delete all users, posts and feeds
  ```terminal
  gator reset
  ```

- **version** - Display the current version of Gator
  ```terminal
  gator version
  ```

- **help** - Show all available commands
  ```terminal
  gator help
  ```

