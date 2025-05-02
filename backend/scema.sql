-- Create students table
CREATE TABLE students (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  date_of_birth DATE,
  address TEXT,
  contact VARCHAR(255),
  program VARCHAR(255)
);

-- Create courses table
CREATE TABLE courses (
  id SERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE NOT NULL,
  title VARCHAR(255) NOT NULL,
  credits INT NOT NULL
);

-- Create enrollments table
CREATE TABLE enrollments (
  id SERIAL PRIMARY KEY,
  student_id INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
  course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  enrollment_date DATE DEFAULT CURRENT_DATE,
  UNIQUE (student_id, course_id) -- Ensure a student can only enroll in a course once
);

-- Create grades table
CREATE TABLE grades (
  id SERIAL PRIMARY KEY,
  enrollment_id INT NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
  grade DECIMAL(3, 2), -- Example: 4.00 for GPA
  semester VARCHAR(50),
  UNIQUE (enrollment_id, semester) -- Ensure only one grade per enrollment per semester
);

