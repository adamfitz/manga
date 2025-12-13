# Manga Tracker - Fyne GUI Application

A complete cross-platform manga tracking application built with Go and Fyne, featuring PostgreSQL database integration.

## Features

- **Library Management**: Add, view, update, and delete manga entries
- **Bookmarking**: Bookmark manga with notes
- **MangaDex Integration**: Search and import manga from MangaDex
- **PostgreSQL Database**: All data stored in PostgreSQL
- **GUI Database Configuration**: Easy-to-use configuration dialog
- **Cross-platform**: Runs on Windows, macOS, and Linux
- **Modular Architecture**: Easy to extend with new features

## Setup Instructions

### 1. Prerequisites

- Go 1.21 or higher
- PostgreSQL database server
- Git

### 2. Database Setup

Create a PostgreSQL database:
```sql
CREATE DATABASE manga_db;
```

### 3. Installation
```bash
# Clone or create the project directory
mkdir manga
cd manga

# Initialize Go module
go mod init github.com/yourusername/manga

# Download dependencies
go get fyne.io/fyne/v2@latest
go get github.com/lib/pq

# Build the application
go build

# Run the application
./manga  # On Linux/Mac
manga.exe  # On Windows
```

### 4. First Run Configuration

When you run the application for the first time:

1. The application will start and show "No Database Connection"
2. Go to **File > Database Settings** in the menu
3. Enter your database configuration:
   - **Server**: Your PostgreSQL server address (e.g., `localhost`)
   - **Port**: Database port (default: `5432`)
   - **Username**: Your PostgreSQL username
   - **Password**: Your PostgreSQL password
   - **Database**: The database name you created (e.g., `manga_db`)
4. Click **Test Connection** to verify your settings
5. Click **Save & Connect** to save and connect

The configuration is automatically saved to `~/.config/manga.config` and will be loaded on subsequent launches.

## Using the Application

### Menu Options

- **File > Database Settings**: Open database configuration dialog
- **File > Quit**: Exit the application
- **Help > About**: View application information

### Managing Your Library

1. **Adding Manga Manually**:
   - Go to the **Library** tab
   - Click **Add Manga**
   - Fill in the manga details
   - Click **Add**

2. **Searching and Importing from MangaDex**:
   - Go to the **Search** tab
   - Enter a manga title
   - Click **Search**
   - Select a manga from the results
   - Click **Add to Library**

3. **Deleting Manga**:
   - Select a manga from your library
   - Click **Delete Manga**
   - Confirm the deletion

### Managing Bookmarks

1. **Adding a Bookmark**:
   - Select a manga from your library
   - Click **Add Bookmark**
   - Add an optional note
   - Click **Add**

2. **Viewing Bookmarks**:
   - Go to the **Bookmarks** tab
   - Select a bookmark to view details

3. **Removing Bookmarks**:
   - Select a bookmark
   - Click **Remove Bookmark**

## License

MIT License