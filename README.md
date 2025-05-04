# Student Information System

A simple web application for managing student, course, enrollment, and grade information. The system has separate views and functionalities for students and faculty members.

## Features

* **User Authentication:** Login for students and faculty.
* **Student Management (Faculty Only):** View, add, update, and delete student records.
* **Course Management (Faculty Only):** View, add, update, and delete course records.
* **Enrollment Management (Faculty Only):** View, create, and delete student enrollments in courses.
* **Grade Management (Faculty Only):** View, add, edit, and delete grades for student enrollments.
* **Student View:** Students can view their personal details, transcript, and calculated GPA.
* **Faculty View:** Faculty can view lists of students, courses, and enrollments, and manage student details, grades, etc.
* **Responsive Frontend:** Designed to be usable on different screen sizes.

## Technologies Used

**Backend:**

* **Go:** Programming language
* **Fiber:** Fast and minimalist web framework
* **pgx:** PostgreSQL driver
* **golang-jwt/jwt/v5:** JWT implementation
* **Viper:** Environment variable management (Implicitly used for JWT Secret)
* **PostgreSQL:** Database

**Frontend:**

* **React:** JavaScript library for building user interfaces
* **TypeScript:** Typed superset of JavaScript
* **Vite:** Frontend build tool
* **Tailwind CSS:** Utility-first CSS framework for styling
* **Axios:** Promise-based HTTP client

## Setup and Installation

### Prerequisites

* Go (version 1.18 or higher recommended)
* Node.js and npm or yarn
* PostgreSQL database server

### Backend Setup

1.  Clone the repository (assuming your code is in a repository).
2.  Navigate to the backend directory.
3.  Install Go dependencies:
```bash
go mod tidy
```
4.  Create a `.env` file in the backend directory with your database connection string and JWT secret:
5.  Run the database schema and seed data script (`scema.sql`) on your PostgreSQL database. You can use a client like `psql`:
```bash
psql -U your_user -d your_dbname -h your_host -p your_port -f scema.sql
```
6.  Build and run the backend application:
```bash
go run main.go
```
The backend should start on `http://127.0.0.1:3000` by default.

### Frontend Setup

1.  Navigate to the frontend directory.
2.  Install Node.js dependencies:
```bash
npm install
# or
yarn install
```
3.  Create a `.env` file in the frontend directory to specify the backend API URL:
```env
VITE_API_URL="http://127.0.0.1:3000"
```
4.  Run the frontend application:
```bash
npm run dev
# or
yarn dev
```
The frontend should start on `http://localhost:5173` by default.

## Running the Application

1.  Ensure your PostgreSQL database is running.
2.  Start the backend application (see Backend Setup).
3.  Start the frontend application (see Frontend Setup).
4.  Open your web browser and navigate to the frontend URL (e.g., `http://localhost:5173`).

## API Endpoints (Overview)

The backend provides RESTful API endpoints for managing the system data. Authentication is required for most endpoints.

* `/login` \[POST]: Authenticate users (student or faculty).
* `/students` \[GET, POST]: Get all students (Faculty) or own student (Student), Create student (Faculty).
* `/students/:id` \[GET, PUT, DELETE]: Get, Update, Delete a specific student.
* `/students/:id/transcript` \[GET]: Get transcript for a student.
* `/students/:id/gpa` \[GET]: Calculate GPA for a student.
* `/courses` \[GET, POST]: Get all courses, Create course (Faculty).
* `/courses/:id` \[GET, PUT, DELETE]: Get, Update, Delete a specific course.
* `/enrollments` \[GET, POST]: Get all enrollments (Faculty) or own enrollments (Student), Create enrollment (Faculty). Supports `?student_id=<id>` query for faculty.
* `/enrollments/:id` \[GET, DELETE]: Get, Delete a specific enrollment.
* `/grades` \[GET, POST]: Get all grades (Faculty) or own grades (Student), Add grade (Faculty).
* `/grades/:id` \[GET, PUT, DELETE]: Get, Update, Delete a specific grade.

## Authentication

The application uses JWT (JSON Web Tokens) for authentication. Upon successful login, the backend returns a JWT, which the frontend stores (currently in `localStorage`) and sends in the `Authorization: Bearer <token>` header for subsequent authenticated requests.

## Future Improvements

* Implement secure password hashing.
* Implement backend rate limiting for login attempts.
* Implement pagination for large lists (students, courses, enrollments, grades).
* Enhance UI/UX and accessibility.
* Write unit and integration tests for both backend and frontend.
* Set up HTTPS for secure communication.

