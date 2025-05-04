DROP TABLE IF EXISTS grades CASCADE;
DROP TABLE IF EXISTS enrollments CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
DROP TABLE IF EXISTS faculty CASCADE;
DROP TABLE IF EXISTS students CASCADE;

CREATE TABLE students (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL DEFAULT '',
  date_of_birth DATE NOT NULL,
  address TEXT NOT NULL,
  contact VARCHAR(255) NOT NULL,
  program VARCHAR(255) NOT NULL
);

CREATE TABLE faculty (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL DEFAULT '',
  date_of_birth DATE NOT NULL,
  info TEXT NOT NULL
);

CREATE TABLE courses (
  id SERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE NOT NULL,
  title VARCHAR(255) NOT NULL,
  credits DECIMAL(3, 2) NOT NULL
);

CREATE TABLE enrollments (
  id SERIAL PRIMARY KEY,
  student_id INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
  course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  enrollment_date DATE DEFAULT CURRENT_DATE NOT NULL,
  UNIQUE (student_id, course_id)
);

CREATE TABLE grades (
  id SERIAL PRIMARY KEY,
  enrollment_id INT NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
  grade DECIMAL(3, 2),
  semester INT NOT NULL,
  UNIQUE (enrollment_id, semester)
);

-- Seed data - NOTE: Passwords should be hashed in a real application
INSERT INTO students (name, password, date_of_birth, address, contact, program) VALUES
('Alice Smith', '', '2002-05-15', '123 Main St, Anytown', '555-1234', 'Computer Science'),
('Bob Johnson', '', '2003-11-20', '456 Oak Ave, Somewhere', '555-5678', 'Electrical Engineering'),
('Charlie Brown', '', '2001-07-01', '789 Pine Ln, Nowhere', '555-9012', 'Physics'),
('Diana Prince', '', '2004-03-10', '101 Hero Way, Themyscira', '555-3456', 'History'),
('Ethan Hunt', '', '2003-09-25', '246 Spy Blvd, IMF HQ', '555-7890', 'International Relations');

INSERT INTO faculty (name, password, date_of_birth, info) VALUES
('prof_davis', '', '1975-08-22', 'Dr. Emily Davis, Head of Computer Science'),
('dr_wilson', '', '1968-04-11', 'Dr. John Wilson, Professor of Electrical Engineering'),
('prof_jones', '', '1980-12-03', 'Dr. Sarah Jones, Professor of History');

INSERT INTO courses (code, title, credits) VALUES
('CS101', 'Introduction to Programming', 3.00),
('EE201', 'Circuit Analysis', 4.00),
('PHY101', 'General Physics I', 4.00),
('HIS201', 'World History II', 3.00),
('IR301', 'Global Politics', 3.00);

INSERT INTO enrollments (student_id, course_id, enrollment_date) VALUES
(1, 1, '2023-09-01'),
(1, 3, '2023-09-01'),
(2, 2, '2023-09-01'),
(3, 3, '2023-09-01'),
(4, 4, '2023-09-01'),
(5, 5, '2023-09-01'),
(1, 4, '2024-01-15'),
(2, 1, '2024-01-15');

INSERT INTO grades (enrollment_id, grade, semester) VALUES
(1, 3.80, 20231),
(2, 3.50, 20231),
(3, 4.00, 20231),
(4, 3.20, 20231),
(5, 3.90, 20231),
(6, 3.70, 20231),
(7, 3.00, 20242);

