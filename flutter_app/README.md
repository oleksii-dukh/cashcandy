# CashCandy Flutter App

A complete money-saving app built with Flutter and Go backend.

## Getting Started

This Flutter app connects to your Go backend to provide a full-featured savings goal tracking application.

### 🚀 Quick Setup

1. **Install dependencies:**

   ```bash
   flutter pub get
   ```

2. **Generate JSON serialization:**

   ```bash
   dart run build_runner build
   ```

3. **Start the Go backend:**

   ```bash
   cd ../go-backend
   go run main.go
   ```

4. **Run the Flutter app:**
   ```bash
   flutter run
   ```

### 📱 Features

- User authentication (login/register)
- Dashboard with savings statistics
- Create and manage savings goals
- Add/remove money from goals
- Progress tracking with visual indicators
- Real-time updates

### 🏗️ Architecture

The app uses a clean architecture with:

- **Models:** Data structures for API communication
- **Services:** API calls and local storage
- **Providers:** State management (similar to Svelte stores)
- **Screens:** UI pages/routes
- **Widgets:** Reusable UI components

For detailed documentation, see the full README in the project root.
