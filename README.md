
# MangaDB — Desktop frontend (Fyne)

![Version](https://img.shields.io/github/v/tag/adamfitz/mangadb?label=version)
![Go Reference](https://img.shields.io/badge/Go-Module-blue)
![Go Version](https://img.shields.io/github/go-mod/go-version/adamfitz/mangadb)




A compact desktop frontend written in Go using the Fyne GUI toolkit. MangaDB is a client application that reads and writes a PostgreSQL database you must provide and maintain. It is a hobby/experimental project — do not use this with production or critical databases.

Summary

- Frontend-only desktop app (no DB server included).
- You must create, secure, and back up the PostgreSQL database yourself.
- The app can run built-in migrations on startup, or you can apply SQL manually.

Quick start

- Prerequisites: Go 1.21+, PostgreSQL
- Latest release: download the prebuilt binaries from the releases page: https://github.com/adamfitz/mangadb/releases

Database (minimal requirement)

You must provide your own PostgreSQL database. The database name does not matter and no DB server setup instructions are provided here — this project does not include or manage a database server.

The application requires the following tables and columns to exist. Create these tables in your database before connecting (the app's code contains the canonical migrations in `database/migrations.go`).

Required tables and columns

- `manga`
  - `id` (serial primary key)
  - `title` (varchar(255), not null)
  - `alt_title` (varchar(255))
  - `author` (varchar(255))
  - `description` (text)
  - `cover_url` (text)
  - `url` (text, not null)
  - `mangadex_id` (varchar(255))
  - `status` (varchar(50), expected values: 'ongoing','completed','hiatus','cancelled')
  - `created_at` (timestamp)
  - `updated_at` (timestamp)

- `anime`
  - `id` (serial primary key)
  - `name` (varchar(255), not null)
  - `alt_name` (varchar(255))
  - `url` (text, not null)
  - `status` (varchar(50), expected values: 'ongoing','completed','hiatus','cancelled')

- `lightnovel` (same as `anime`, plus `volumes`)
  - `id` (serial primary key)
  - `name` (varchar(255), not null)
  - `alt_name` (varchar(255))
  - `url` (text, not null)
  - `status` (varchar(50), expected values: 'ongoing','completed','hiatus','cancelled')
  - `volumes` (integer)

- `webnovel` (same as `anime`)
  - `id` (serial primary key)
  - `name` (varchar(255), not null)
  - `alt_name` (varchar(255))
  - `url` (text, not null)
  - `status` (varchar(50), expected values: 'ongoing','completed','hiatus','cancelled')

- `webtoons` (same as `anime`)
  - `id` (serial primary key)
  - `name` (varchar(255), not null)
  - `alt_name` (varchar(255))
  - `url` (text, not null)
  - `status` (varchar(50), expected values: 'ongoing','completed','hiatus','cancelled')

Configuration

- Open File → Database Settings in the app and enter host, port (default 5432), username, password and database name.
- Configuration is saved to `~/.config/mangadb.config` (JSON).

Fonts

- The app uses static CJK OTF fonts for proper rendering (e.g. `NotoSansCJKjp-Regular.otf`). Avoid variable-font (VF) files.


Support and disclaimers

- This is an experimental/hobby project. Do not use this on production or critical databases.
- Review migrations before running them. Back up your data regularly.

What the app does — usage (short)

- Library: view your library in the Library tab. Add a new entry with "Add Manga" (title, alt title, author, description, cover URL, source URL, status). Edit or delete entries from the library view.
- Search & import: use the Search/MangaDex interface to find titles and import them into your local DB (imported rows populate the `manga` table fields shown in the schema).
- Bookmarks: add, view and remove bookmarks for items in your library (notes are stored alongside bookmark entries).
- Window/menu: use File → Database Settings to configure DB connection, File → Quit to exit, Help → About for version info.

How to use (step-by-step)

1. Build and run the app (see Quick start).
2. On first run, open File → Database Settings. Enter host, port, username/password, and database name. Test connection, then Save & Connect.
3. If you prefer the app to create tables automatically, let it run migrations on connect. If you prefer manual control, apply the SQL in the Schema section first.
4. Use the Library and Search tabs to add content. Use Bookmarks to add notes/positions.


Notes about migrations

- Migrations are kept in code only (`database/migrations.go`). The app will attempt to run them on startup if connected.
- If you want an SQL file to run manually, extract the statements from `database/migrations.go` and apply them in your DB tool.

License: MIT
