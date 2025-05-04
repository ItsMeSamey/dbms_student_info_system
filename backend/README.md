# Student Information System

A web application designed to manage academic data for students and faculty. It provides different interfaces and capabilities based on user roles, with a significant portion of the business logic implemented within the database using PL/SQL.

## Features

* **Faculty Portal:** Faculty can view lists of students, courses, and enrollments, and access detailed student information including transcripts for grade management.
* **User Authentication:** Secure login for both students and faculty members.
* **Student Management (Faculty Only):** Faculty can perform CRUD (Create, Read, Update, Delete) operations on student records.
* **Course Management (Faculty Only):** Faculty can manage course information (Create, Read, Update, Delete).
* **Enrollment Management (Faculty Only):** Faculty can enroll students in courses and remove enrollments.
* **Grade Management (Faculty Only):** Faculty can add, edit, and delete grades for student enrollments.

* **Student Portal:** Students can view their personal information, academic transcript, and calculated GPA.

## Architecture Overview

This project follows a layered architecture. The frontend is a React application that communicates with a Go backend via RESTful APIs. A key aspect of the backend is that it delegates much of the data interaction, validation, and core business logic to the PostgreSQL database using **PL/SQL functions and procedures**.

The Go backend acts as a thin layer, primarily handling:
* Receiving HTTP requests.
* Parsing request data.
* Authenticating users (by calling a PL/SQL function).
* Calling the appropriate PL/SQL function or procedure in the database.
* Handling database responses and errors (including interpreting PL/SQL exceptions).
* Formatting responses for the frontend.

The PL/SQL layer handles:
* Direct database queries (SELECT, INSERT, UPDATE, DELETE).
* Data validation (checking for required fields, valid IDs, unique constraints).
* Authorization checks based on the `user_id` and `user_role` passed from the backend.
* Complex operations like calculating GPA and generating transcripts.
* Raising exceptions to signal errors or access violations back to the backend.

## Technologies Used

**Backend:**

* **Go:** The primary language for the API server.
* **Fiber:** A web framework for building the API.
* **pgx:** A high-performance PostgreSQL driver for Go.
* **golang-jwt/jwt/v5:** For handling JWT authentication.

**Database:**

* **PostgreSQL:** The relational database storing all application data.
* **PL/SQL:** Used to write stored procedures and functions that encapsulate database logic.

## Key PL/SQL Functions and Procedures

The following PL/SQL functions and procedures are central to the application's logic:

* `authenticate_user(p_id INT, p_password VARCHAR, p_role VARCHAR)`:
    * **Purpose:** Verifies user credentials against the `students` or `faculty` tables.
    * **Logic:** Checks if a user with the given ID and role exists and if the provided password (or DOB for default) matches.
    * **Returns:** `user_id` and `user_role` if successful.
    * **Raises Exception:** 'Invalid credentials' or 'Invalid role specified'.

* `create_student(p_name VARCHAR, p_password VARCHAR, p_date_of_birth DATE, p_address TEXT, p_contact VARCHAR, p_program VARCHAR)`:
    * **Purpose:** Inserts a new student record.
    * **Logic:** Performs basic validation (name, DOB) and inserts into the `students` table.
    * **Returns:** The ID of the newly created student.
    * **Raises Exception:** 'Student name is required', 'Student date of birth is required', or database errors.

* `get_students(p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Retrieves student records based on the requesting user's role.
    * **Logic:** If the user is a 'student', returns only their own record. If 'faculty', returns all student records.
    * **Returns:** A set of `students` records.
    * **Raises Exception:** 'Access denied. Invalid user role.'

* `get_student_by_id(p_student_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Retrieves a single student record by ID with authorization.
    * **Logic:** Checks if the requesting user is authorized (student can only get their own). Selects the student record.
    * **Returns:** A single `students` record.
    * **Raises Exception:** 'Access denied...', 'Student not found'.

* `update_student(p_student_id INT, p_name VARCHAR, ..., p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Updates an existing student record.
    * **Logic:** Performs authorization (only faculty) and basic validation. Updates the `students` table.
    * **Raises Exception:** 'Access denied...', 'Student not found', validation errors, or database errors.

* `delete_student(p_student_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Deletes a student record.
    * **Logic:** Performs authorization (only faculty) and deletes from the `students` table.
    * **Raises Exception:** 'Access denied...', 'Student not found', or database errors.

* `create_course(p_code VARCHAR, p_title VARCHAR, p_credits DECIMAL, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Inserts a new course record.
    * **Logic:** Performs authorization (only faculty) and validation. Inserts into `courses`. Handles unique code constraint.
    * **Returns:** The ID of the new course.
    * **Raises Exception:** 'Access denied...', validation errors, 'Course with code already exists', or database errors.

* `get_all_courses()`:
    * **Purpose:** Retrieves all course records.
    * **Logic:** Selects all from `courses`. Authorization is expected to be handled by the calling Go handler's middleware.
    * **Returns:** A set of `courses` records.

* `get_course_by_id(p_course_id INT)`:
    * **Purpose:** Retrieves a single course by ID.
    * **Logic:** Selects from `courses`.
    * **Returns:** A single `courses` record.
    * **Raises Exception:** 'Course not found'.

* `update_course(p_course_id INT, p_code VARCHAR, ..., p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Updates a course record.
    * **Logic:** Performs authorization (only faculty) and validation. Updates `courses`. Handles unique code constraint.
    * **Raises Exception:** 'Access denied...', 'Course not found', validation errors, 'Course with code already exists', or database errors.

* `delete_course(p_course_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Deletes a course record.
    * **Logic:** Performs authorization (only faculty) and deletes from `courses`.
    * **Raises Exception:** 'Access denied...', 'Course not found', or database errors.

* `create_enrollment(p_student_id INT, p_course_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Creates a new enrollment record.
    * **Logic:** Performs authorization (only faculty) and validation (required IDs, existence of student/course). Inserts into `enrollments`. Handles unique student+course constraint. Uses a CTE to return the new ID and date.
    * **Returns:** The ID and enrollment date of the new enrollment.
    * **Raises Exception:** 'Access denied...', validation errors, 'Invalid student ID or course ID', 'Student is already enrolled...', or database errors.

* `get_enrollments(p_user_id INT, p_user_role VARCHAR, p_filter_student_id INT DEFAULT NULL)`:
    * **Purpose:** Retrieves enrollment records based on user role and optional student filter.
    * **Logic:** If 'student', returns own enrollments. If 'faculty', returns all or filtered by student ID.
    * **Returns:** A set of `enrollments` records.
    * **Raises Exception:** 'Access denied...'.

* `get_enrollment_by_id(p_enrollment_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Retrieves a single enrollment by ID with authorization.
    * **Logic:** Checks if the requesting user is authorized (student can only get their own). Selects the enrollment record.
    * **Returns:** A single `enrollments` record.
    * **Raises Exception:** 'Access denied...', 'Enrollment not found'.

* `delete_enrollment(p_enrollment_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Deletes an enrollment record.
    * **Logic:** Performs authorization (only faculty) and deletes from `enrollments`.
    * **Raises Exception:** 'Access denied...', 'Enrollment not found', or database errors.

* `add_grade(p_enrollment_id INT, p_grade DECIMAL, p_semester INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Adds a grade to an enrollment.
    * **Logic:** Performs authorization (only faculty) and validation (required enrollment ID, semester, existence of enrollment). Inserts into `grades`. Handles unique enrollment+semester constraint.
    * **Returns:** The ID of the new grade.
    * **Raises Exception:** 'Access denied...', validation errors, 'Invalid enrollment ID', 'Grade for this enrollment and semester already exists', or database errors.

* `get_all_grades(p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Retrieves grade records based on user role.
    * **Logic:** If 'student', returns grades for their enrollments. If 'faculty', returns all grades.
    * **Returns:** A set of `grades` records.
    * **Raises Exception:** 'Access denied...'.

* `get_grade_by_id(p_grade_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Retrieves a single grade by ID with authorization.
    * **Logic:** Checks if the requesting user is authorized (student can only get grades for their own enrollments). Selects the grade record.
    * **Returns:** A single `grades` record.
    * **Raises Exception:** 'Access denied...', 'Grade not found'.

* `update_grade(p_grade_id INT, p_enrollment_id INT, p_grade DECIMAL, p_semester INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Updates a grade record.
    * **Logic:** Performs authorization (only faculty) and validation. Updates `grades`. Handles unique enrollment+semester constraint.
    * **Raises Exception:** 'Access denied...', 'Grade not found', validation errors, 'Grade for this enrollment and semester already exists', or database errors.

* `delete_grade(p_grade_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Deletes a grade record.
    * **Logic:** Performs authorization (only faculty) and deletes from `grades`.
    * **Raises Exception:** 'Access denied...', 'Grade not found', or database errors.

* `get_student_transcript(p_student_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Generates a student's academic transcript.
    * **Logic:** Performs authorization. Joins `enrollments`, `courses`, and `grades` (using a LEFT JOIN to include courses without grades). Orders by semester and course code.
    * **Returns:** A set of transcript rows including enrollment details, course info, and grade info (if available).
    * **Raises Exception:** 'Access denied...', 'Student not found'.

* `calculate_student_gpa(p_student_id INT, p_user_id INT, p_user_role VARCHAR)`:
    * **Purpose:** Calculates a student's GPA.
    * **Logic:** Performs authorization. Calculates the weighted average of grades for courses with assigned grades. Handles cases with no graded courses (returns 0.0).
    * **Returns:** The calculated GPA as a DECIMAL.
    * **Raises Exception:** 'Access denied...', 'Student not found', or database errors.

## Usage (Backend only)

1.  Install Go.
3.  Set up backend dependencies (`go mod tidy`).
4.  Configure backend `.env` (DATABASE\_URL, JWT\_SECRET). **Note: Securely manage secrets in production.**
5.  **Run the `scema.sql` script on your PostgreSQL database.** This creates tables and defines all PL/SQL functions/procedures. **Warning: This script drops existing tables and data.**
6.  Start the backend (`go run main.go`).

