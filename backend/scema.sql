-- Drop existing tables and sequences (order matters due to foreign keys)
DROP TABLE IF EXISTS grades CASCADE;
DROP TABLE IF EXISTS enrollments CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
DROP TABLE IF EXISTS faculty CASCADE;
DROP TABLE IF EXISTS students CASCADE;

-- Drop existing sequences if they exist (for clean recreation)
DROP SEQUENCE IF EXISTS students_id_seq CASCADE;
DROP SEQUENCE IF EXISTS faculty_id_seq CASCADE;
DROP SEQUENCE IF EXISTS courses_id_seq CASCADE;
DROP SEQUENCE IF EXISTS enrollments_id_seq CASCADE;
DROP SEQUENCE IF EXISTS grades_id_seq CASCADE;

-- Create tables
CREATE TABLE students (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL DEFAULT '',
  date_of_birth DATE NOT NULL,
  address TEXT NOT NULL DEFAULT '',
  contact VARCHAR(255) NOT NULL DEFAULT '',
  program VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE faculty (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL DEFAULT '',
  date_of_birth DATE NOT NULL,
  info TEXT NOT NULL DEFAULT ''
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

-- Seed data (WARNING: Passwords should be hashed in a real application)
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

CREATE OR REPLACE FUNCTION authenticate_user(
  p_id INT,
  p_password VARCHAR,
  p_role VARCHAR
)
RETURNS TABLE (
  user_id INT,
  user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_stored_password VARCHAR;
  v_date_of_birth DATE;
  v_dob_password VARCHAR;
BEGIN
  IF p_role = 'student' THEN
    SELECT password, date_of_birth INTO v_stored_password, v_date_of_birth FROM students WHERE id = p_id;
  ELSIF p_role = 'faculty' THEN
    SELECT password, date_of_birth INTO v_stored_password, v_date_of_birth FROM faculty WHERE id = p_id;
  ELSE
    RAISE EXCEPTION 'Invalid role specified';
  END IF;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Invalid credentials';
  END IF;

  IF v_stored_password = '' THEN
    IF v_date_of_birth IS NULL THEN
       RAISE EXCEPTION 'Invalid credentials';
    END IF;
    v_dob_password := to_char(v_date_of_birth, 'YYYY-MM-DD');
    IF p_password != v_dob_password THEN
      RAISE EXCEPTION 'Invalid credentials';
    END IF;
  ELSE
    IF p_password != v_stored_password THEN
      RAISE EXCEPTION 'Invalid credentials';
    END IF;
  END IF;

  RETURN QUERY SELECT p_id, p_role;

EXCEPTION
  WHEN NO_DATA_FOUND THEN
    RAISE EXCEPTION 'Invalid credentials';
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Authentication failed: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION create_student(
  p_name VARCHAR,
  p_password VARCHAR,
  p_date_of_birth DATE,
  p_address TEXT,
  p_contact VARCHAR,
  p_program VARCHAR
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_student_id INT;
BEGIN
  IF p_name IS NULL OR p_name = '' THEN
    RAISE EXCEPTION 'Student name is required';
  END IF;
  IF p_date_of_birth IS NULL THEN
    RAISE EXCEPTION 'Student date of birth is required';
  END IF;

  INSERT INTO students (name, password, date_of_birth, address, contact, program)
  VALUES (p_name, p_password, p_date_of_birth, p_address, p_contact, p_program)
  RETURNING id INTO v_student_id;

  RETURN v_student_id;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to create student: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION get_students(
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS SETOF students
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role = 'student' THEN
    RETURN QUERY SELECT * FROM students WHERE id = p_user_id;
  ELSIF p_user_role = 'faculty' THEN
    RETURN QUERY SELECT * FROM students;
  ELSE
    RAISE EXCEPTION 'Access denied. Invalid user role.';
  END IF;
END;
$$;

CREATE OR REPLACE FUNCTION get_student_by_id(
  p_student_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS SETOF students
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role = 'student' AND p_student_id != p_user_id THEN
    RAISE EXCEPTION 'Access denied. Students can only view their own details.';
  END IF;

  RETURN QUERY SELECT * FROM students WHERE id = p_student_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Student not found';
  END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE update_student(
  p_student_id INT,
  p_name VARCHAR,
  p_date_of_birth DATE,
  p_address TEXT,
  p_contact VARCHAR,
  p_program VARCHAR,
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can update student details.';
  END IF;

   IF p_name IS NULL OR p_name = '' THEN
    RAISE EXCEPTION 'Student name is required';
  END IF;
   IF p_date_of_birth IS NULL THEN
    RAISE EXCEPTION 'Student date of birth is required';
  END IF;

  UPDATE students
  SET name = p_name,
    date_of_birth = p_date_of_birth,
    address = p_address,
    contact = p_contact,
    program = p_program
  WHERE id = p_student_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Student not found';
  END IF;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to update student: %', SQLERRM;
END;
$$;

CREATE OR REPLACE PROCEDURE delete_student(
  p_student_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can delete students.';
  END IF;

  DELETE FROM students WHERE id = p_student_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Student not found';
  END IF;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to delete student: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION create_course(
  p_code VARCHAR,
  p_title VARCHAR,
  p_credits DECIMAL(3, 2),
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_course_id INT;
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can create courses.';
  END IF;

  IF p_code IS NULL OR p_code = '' THEN
    RAISE EXCEPTION 'Course code is required';
  END IF;
  IF p_title IS NULL OR p_title = '' THEN
    RAISE EXCEPTION 'Course title is required';
  END IF;
  IF p_credits IS NULL OR p_credits <= 0 THEN
    RAISE EXCEPTION 'Positive credits are required';
  END IF;

  INSERT INTO courses (code, title, credits)
  VALUES (p_code, p_title, p_credits)
  RETURNING id INTO v_course_id;

  RETURN v_course_id;

EXCEPTION
  WHEN unique_violation THEN
    RAISE EXCEPTION 'Course with code % already exists', p_code;
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to create course: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION get_all_courses()
RETURNS SETOF courses
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY SELECT * FROM courses;
END;
$$;

CREATE OR REPLACE FUNCTION get_course_by_id(
  p_course_id INT
)
RETURNS SETOF courses
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY SELECT * FROM courses WHERE id = p_course_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Course not found';
  END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE update_course(
  p_course_id INT,
  p_code VARCHAR,
  p_title VARCHAR,
  p_credits DECIMAL(3, 2),
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can update courses.';
  END IF;

  IF p_code IS NULL OR p_code = '' THEN
    RAISE EXCEPTION 'Course code is required';
  END IF;
  IF p_title IS NULL OR p_title = '' THEN
    RAISE EXCEPTION 'Course title is required';
  END IF;
  IF p_credits IS NULL OR p_credits <= 0 THEN
    RAISE EXCEPTION 'Positive credits are required';
  END IF;

  UPDATE courses
  SET code = p_code,
    title = p_title,
    credits = p_credits
  WHERE id = p_course_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Course not found';
  END IF;

EXCEPTION
  WHEN unique_violation THEN
    RAISE EXCEPTION 'Course with code % already exists', p_code;
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to update course: %', SQLERRM;
END;
$$;

CREATE OR REPLACE PROCEDURE delete_course(
  p_course_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can delete courses.';
  END IF;

  DELETE FROM courses WHERE id = p_course_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Course not found';
  END IF;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to delete course: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION create_enrollment(
  p_student_id INT,
  p_course_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS TABLE (
  v_id INT,
  v_date DATE
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_enrollment_id INT;
  v_enrollment_date DATE;
  v_student_exists BOOLEAN;
  v_course_exists BOOLEAN;
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can create enrollments.';
  END IF;

  IF p_student_id IS NULL OR p_student_id = 0 OR p_course_id IS NULL OR p_course_id = 0 THEN
    RAISE EXCEPTION 'Student ID and Course ID are required';
  END IF;

  SELECT EXISTS(SELECT 1 FROM students WHERE id = p_student_id) INTO v_student_exists;
  SELECT EXISTS(SELECT 1 FROM courses WHERE id = p_course_id) INTO v_course_exists;

  IF NOT v_student_exists OR NOT v_course_exists THEN
    RAISE EXCEPTION 'Invalid student ID or course ID';
  END IF;

  INSERT INTO enrollments (student_id, course_id)
  VALUES (p_student_id, p_course_id)
  RETURNING id, enrollment_date INTO v_enrollment_id, v_enrollment_date;

  RETURN QUERY SELECT v_enrollment_id, v_enrollment_date;

EXCEPTION
  WHEN unique_violation THEN
    RAISE EXCEPTION 'Student is already enrolled in this course';
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to create enrollment: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION get_enrollments(
  p_user_id INT,
  p_user_role VARCHAR,
  p_filter_student_id INT DEFAULT NULL
)
RETURNS SETOF enrollments
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role = 'student' THEN
    RETURN QUERY SELECT * FROM enrollments WHERE student_id = p_user_id;
  ELSIF p_user_role = 'faculty' THEN
    IF p_filter_student_id IS NOT NULL THEN
      RETURN QUERY SELECT * FROM enrollments WHERE student_id = p_filter_student_id;
    ELSE
      RETURN QUERY SELECT * FROM enrollments;
    END IF;
  ELSE
    RAISE EXCEPTION 'Access denied. Invalid user role.';
  END IF;
END;
$$;

CREATE OR REPLACE FUNCTION get_enrollment_by_id(
  p_enrollment_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS SETOF enrollments
LANGUAGE plpgsql
AS $$
DECLARE
  v_student_id INT;
  v_enrollment enrollments;
BEGIN
  SELECT * INTO v_enrollment FROM enrollments WHERE id = p_enrollment_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Enrollment not found';
  END IF;

  IF p_user_role = 'student' AND v_enrollment.student_id != p_user_id THEN
    RAISE EXCEPTION 'Access denied. Students can only view their own enrollments.';
  END IF;

  RETURN NEXT v_enrollment;
END;
$$;

CREATE OR REPLACE PROCEDURE delete_enrollment(
  p_enrollment_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_student_id INT;
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can delete enrollments.';
  END IF;

  DELETE FROM enrollments WHERE id = p_enrollment_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Enrollment not found';
  END IF;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to delete enrollment: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION add_grade(
  p_enrollment_id INT,
  p_grade DECIMAL(3, 2),
  p_semester INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
  v_grade_id INT;
  v_enrollment_student_id INT;
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can add grades.';
  END IF;

  IF p_enrollment_id IS NULL OR p_enrollment_id = 0 OR p_semester IS NULL OR p_semester = 0 THEN
    RAISE EXCEPTION 'Enrollment ID and Semester are required';
  END IF;

  SELECT student_id INTO v_enrollment_student_id FROM enrollments WHERE id = p_enrollment_id;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'Invalid enrollment ID';
  END IF;

  INSERT INTO grades (enrollment_id, grade, semester)
  VALUES (p_enrollment_id, p_grade, p_semester)
  RETURNING id INTO v_grade_id;

  RETURN v_grade_id;

EXCEPTION
  WHEN unique_violation THEN
    RAISE EXCEPTION 'Grade for this enrollment and semester already exists';
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to add grade: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION get_all_grades(
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS SETOF grades
LANGUAGE plpgsql
AS $$
BEGIN
  IF p_user_role = 'student' THEN
    RETURN QUERY
    SELECT g.*
    FROM grades g
    JOIN enrollments e ON g.enrollment_id = e.id
    WHERE e.student_id = p_user_id;
  ELSIF p_user_role = 'faculty' THEN
    RETURN QUERY SELECT * FROM grades;
  ELSE
    RAISE EXCEPTION 'Access denied. Invalid user role.';
  END IF;
END;
$$;

-- Corrected function definition (removed duplicate OR)
CREATE OR REPLACE FUNCTION get_grade_by_id(
  p_grade_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS SETOF grades
LANGUAGE plpgsql
AS $$
DECLARE
  v_grade grades;
  v_enrollment_student_id INT;
BEGIN
  SELECT * INTO v_grade FROM grades WHERE id = p_grade_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Grade not found';
  END IF;

  IF p_user_role = 'student' THEN
    SELECT student_id INTO v_enrollment_student_id FROM enrollments WHERE id = v_grade.enrollment_id;
    IF NOT FOUND THEN
     RAISE EXCEPTION 'Internal error: Enrollment not found for grade.';
    END IF;
    IF v_enrollment_student_id != p_user_id THEN
     RAISE EXCEPTION 'Access denied. Students can only view grades for their own enrollments.';
    END IF;
  END IF;

  RETURN NEXT v_grade;
END;
$$;

CREATE OR REPLACE PROCEDURE update_grade(
  p_grade_id INT,
  p_enrollment_id INT,
  p_grade DECIMAL(3, 2),
  p_semester INT,
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
   v_enrollment_student_id INT;
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can update grades.';
  END IF;

  IF p_enrollment_id IS NULL OR p_enrollment_id = 0 OR p_semester IS NULL OR p_semester = 0 THEN
    RAISE EXCEPTION 'Enrollment ID and Semester are required';
  END IF;

  SELECT student_id INTO v_enrollment_student_id FROM enrollments WHERE id = p_enrollment_id;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'Invalid enrollment ID';
  END IF;

  UPDATE grades
  SET enrollment_id = p_enrollment_id,
    grade = p_grade,
    semester = p_semester
  WHERE id = p_grade_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Grade not found';
  END IF;

EXCEPTION
  WHEN unique_violation THEN
    RAISE EXCEPTION 'Grade for this enrollment and semester already exists';
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to update grade: %', SQLERRM;
END;
$$;

CREATE OR REPLACE PROCEDURE delete_grade(
  p_grade_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_enrollment_student_id INT;
BEGIN
  IF p_user_role != 'faculty' THEN
    RAISE EXCEPTION 'Access denied. Only faculty can delete grades.';
  END IF;

  DELETE FROM grades WHERE id = p_grade_id;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Grade not found';
  END IF;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to delete grade: %', SQLERRM;
END;
$$;

CREATE OR REPLACE FUNCTION get_student_transcript(
  p_student_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS TABLE (
  enrollment_id INT,
  course_code VARCHAR,
  course_title VARCHAR,
  credits DECIMAL(3, 2),
  grade_id INT,
  grade DECIMAL(3, 2),
  semester INT
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_student_exists BOOLEAN;
BEGIN
  IF p_user_role = 'student' AND p_student_id != p_user_id THEN
    RAISE EXCEPTION 'Access denied. Students can only view their own transcript.';
  END IF;

  SELECT EXISTS(SELECT 1 FROM students WHERE id = p_student_id) INTO v_student_exists;
  IF NOT v_student_exists THEN
    RAISE EXCEPTION 'Student not found';
  END IF;

  RETURN QUERY
  SELECT
    e.id AS enrollment_id,
    c.code,
    c.title,
    c.credits,
    g.id AS grade_id,
    g.grade,
    g.semester
  FROM
    enrollments e
  JOIN
    courses c ON e.course_id = c.id
  LEFT JOIN
    grades g ON e.id = g.enrollment_id
  WHERE
    e.student_id = p_student_id
  ORDER BY
    g.semester NULLS LAST, c.code;

END;
$$;

CREATE OR REPLACE FUNCTION calculate_student_gpa(
  p_student_id INT,
  p_user_id INT,
  p_user_role VARCHAR
)
RETURNS DECIMAL(3, 2)
LANGUAGE plpgsql
AS $$
DECLARE
  v_gpa DECIMAL(3, 2);
  v_student_exists BOOLEAN;
BEGIN
  IF p_user_role = 'student' AND p_student_id != p_user_id THEN
    RAISE EXCEPTION 'Access denied. Students can only calculate their own GPA.';
  END IF;

  SELECT EXISTS(SELECT 1 FROM students WHERE id = p_student_id) INTO v_student_exists;
  IF NOT v_student_exists THEN
    RAISE EXCEPTION 'Student not found';
  END IF;


  SELECT
    SUM(g.grade * c.credits) / NULLIF(SUM(c.credits), 0)
  INTO v_gpa
  FROM
    enrollments e
  JOIN
    courses c ON e.course_id = c.id
  JOIN
    grades g ON e.id = g.enrollment_id
  WHERE
    e.student_id = p_student_id AND g.grade IS NOT NULL;

  IF v_gpa IS NULL THEN
    RETURN 0.0;
  END IF;

  RETURN v_gpa;

EXCEPTION
  WHEN OTHERS THEN
    RAISE EXCEPTION 'Failed to calculate GPA: %', SQLERRM;
END;
$$;

