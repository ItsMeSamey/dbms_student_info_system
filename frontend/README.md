# Student Information System (Frontend)

This is the frontend component of the Student Information System, providing the user interface for students and faculty to interact with the system's data and functionality. It is a single-page application built with React and TypeScript.

## Features

* **User Interface:** Provides distinct views and navigation for both Student and Faculty roles.
* **Authentication Flow:** Allows users to log in and manages the authenticated state.
* **Student Portal Views:**
    * View personal details.
    * View academic transcript.
    * View calculated GPA.
* **Faculty Portal Views:**
    * View lists of students, courses, and enrollments.
    * View detailed information for individual students, including their transcript.
    * Manage (Add, Edit, Delete) grades for student enrollments.
    * Manage (Add, Edit, Delete) student records.
    * Manage (Add, Edit, Delete) course records.
    * Manage (Add, Delete) enrollment records.
* **Data Display:** Presents student, course, enrollment, and grade information in a user-friendly format (tables, forms).
* **Forms and Modals:** Provides forms for data entry (e.g., adding a grade, updating student info) and uses modals for focused interactions.
* **Navigation:** Allows users to navigate between different sections of the application based on their role.
* **Responsive Design:** The UI is built using Tailwind CSS to adapt to various screen sizes.

## Architecture Overview

The frontend is a React single-page application. It communicates with the backend API (built with Go and PostgreSQL/PL/SQL) to fetch and send data.

Key aspects:
* **Component-Based:** The UI is composed of reusable React components.
* **State Management:** Uses React's built-in state management (`useState`, `useEffect`, Context API if needed).
* **API Interaction:** Uses Axios to make asynchronous HTTP requests to the backend endpoints.
* **Routing:** Manages different views/pages within the single-page application structure (implicitly handled by conditional rendering or simple state-based view switching in the provided code structure).
* **Authentication Handling:** Stores the received JWT token (currently in `localStorage`) and includes it in the `Authorization` header for authenticated requests.

The frontend is responsible for:
* Rendering the user interface.
* Handling user input and interactions.
* Making API calls to the backend.
* Displaying data received from the backend.
* Implementing client-side validation for user experience (backend must re-validate).
* Managing the display of loading states and errors received from the backend.

The frontend is *not* responsible for:
* Storing sensitive data securely (passphrases, etc.).
* Implementing core business logic or complex data validation (this is delegated to the backend/PL/SQL).
* Directly accessing the database.

## Technologies Used

* **React:** JavaScript library for building user interfaces.
* **TypeScript:** Adds static typing to improve code quality and maintainability.
* **Vite:** A fast frontend build tool.
* **Tailwind CSS:** A utility-first CSS framework for rapid styling.
* **Axios:** Promise-based HTTP client for making API requests.

## Setup and Installation

### Prerequisites

* Node.js and npm or yarn installed.

### Installation

1.  Navigate to the frontend directory of the project.
2.  Install the dependencies:
    ```bash
    npm install
    # or
    yarn install
    ```
3.  Create a `.env` file in the frontend directory to specify the backend API URL. Replace the URL with the actual address where your backend is running:
    ```env
    VITE_API_URL="[http://127.0.0.1:3000](http://127.0.0.1:3000)"
    ```

## Running the Application

1.  Ensure your backend application is running and accessible at the URL specified in the `.env` file.
2.  Navigate to the frontend directory in your terminal.
3.  Start the development server:
    ```bash
    npm run dev
    # or
    yarn dev
    ```
4.  Open your web browser and navigate to the address provided by the development server (usually `http://localhost:5173`).

## API Interaction

The frontend communicates with the backend API using the following endpoints (based on the backend structure):

* `POST /login`: Authenticates the user.
* `GET /students`: Retrieves student data (all for faculty, own for student).
* `GET /students/:id`: Retrieves details for a specific student.
* `PUT /students/:id`: Updates a student's details (Faculty only).
* `DELETE /students/:id`: Deletes a student (Faculty only).
* `GET /students/:id/transcript`: Retrieves a student's transcript.
* `GET /students/:id/gpa`: Retrieves a student's calculated GPA.
* `GET /courses`: Retrieves all course data.
* `POST /courses`: Creates a new course (Faculty only).
* `PUT /courses/:id`: Updates a course (Faculty only).
* `DELETE /courses/:id`: Deletes a course (Faculty only).
* `GET /enrollments`: Retrieves enrollment data (all/filtered for faculty, own for student).
* `POST /enrollments`: Creates a new enrollment (Faculty only).
* `DELETE /enrollments/:id`: Deletes an enrollment (Faculty only).
* `POST /grades`: Adds a new grade (Faculty only).
* `PUT /grades/:id`: Updates an existing grade (Faculty only).
* `DELETE /grades/:id`: Deletes a grade (Faculty only).

Authenticated requests include the JWT token in the `Authorization: Bearer <token>` header.

## Future Improvements (Frontend)

* Implement more robust form validation and error handling.
* Improve UI/UX and accessibility.
* Add loading indicators and better feedback for asynchronous operations.
* Implement pagination or infinite scrolling for long lists.
* Enhance navigation and routing.
* Add unit and integration tests for React components and API interactions.

