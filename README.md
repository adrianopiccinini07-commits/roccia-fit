# 🏋️‍♂️ RocciaFit - Full-Stack Fitness Coaching Platform

RocciaFit is a comprehensive, full-stack fitness coaching platform designed to bridge the gap between personal trainers and their athletes. Built with a focus on cross-platform mobile experience and a robust, scalable backend, it streamlines workout programming, nutritional tracking, and daily compliance gamification.

## 🚀 Tech Stack

**Frontend (Mobile)**
* **Framework:** React Native (Expo)
* **Language:** TypeScript
* **Routing:** Expo Router (File-based routing)
* **State Management & UI:** React Hooks, Custom Theme Engine (Light/Dark mode)
* **Data Visualization:** `react-native-chart-kit`

**Backend (API)**
* **Language:** Go (Golang)
* **Router:** `go-chi/chi` (Lightweight, idiomatic REST routing)
* **Database:** PostgreSQL
* **ORM:** GORM (Auto-migrations, relational data modeling)
* **Security:** JWT (JSON Web Tokens) Authentication, Bcrypt password hashing

## ✨ Core Features

### 🦅 Coach Dashboard (Management & Programming)
* **Athlete Roster:** Real-time overview of assigned athletes, tracking their XP (Experience Points), body weight trends, and posing gallery.
* **Workout Builder:** Create multi-week training blocks. Searchable exercise catalog, dynamic sets/reps/RIR assignment, and specific execution notes.
* **Diet Planner:** Macro-based diet creation with meal splitting. Includes food alternatives logic to give athletes flexibility while hitting target macros.
* **Cardio & NEAT:** Assign daily step targets and structured cardio sessions (machinery, duration, intensity).
* **Direct Feedback:** Ticket system and direct coaching notes to adjust protocols on the fly.

### ⚡ Athlete Dashboard (Compliance & Gamification)
* **Daily Tracker:** Interactive daily checklist for Workout, Diet, and Cardio compliance.
* **Gamification Engine:** Earn XP for completing daily tasks. Level up and rank progression (Copper ➔ Champion) to boost retention and motivation.
* **Live Workout Mode:** Start a session, log actual weights/reps, utilize the built-in rest timer (with screen wake-lock), and auto-calculate 1RM (Brzycki formula).
* **Physiological Check-ins:** Upload morning weight, sleep quality (1-5), hunger levels (1-5), and progress pictures directly to the Coach.
* **Progress Analytics:** Visual charts mapping body weight trends and training volume over time.

## 🏗️ Architecture & Engineering Highlights
* **Role-Based Access Control (RBAC):** Distinct UI namespaces and API endpoints for `COACH` and `ATLETA` roles.
* **Relational Database Design:** Complex nested structures (User -> Scheda -> Giorno -> Esercizio) handled efficiently via PostgreSQL and GORM.
* **Cross-Platform UI:** Fully responsive design utilizing Safe Area insets and conditional rendering to look native on both iOS and Android.

## 👨‍💻 Author
**Adriano Piccinini**
*Computer Science and Automation Engineering Student @ Sapienza University of Rome*