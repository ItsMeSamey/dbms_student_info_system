Here is a list of SQL queries found in the handlers.go and auth_handlers.go files:-- Used in handlers.go -> CreateStudent
-- Inserts a new student record into the students table.
-- RETURNING id fetches the generated primary key ID of the new student.
INSERT INTO students (name, password, date_of_birth, address, contact, program) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- Used in handlers.go -> GetStudents (Student Role)
-- Selects the details of a single student based on their ID.
-- Used when a student user requests their own details.
SELECT id, name, date_of_birth, address, contact, program FROM students WHERE id = $1;

-- Used in handlers.go -> GetStudents (Faculty Role)
-- Selects the details of all students from the students table.
-- Used when a faculty user requests the list of all students.
SELECT id, name, date_of_birth, address, contact, program FROM students;

-- Used in handlers.go -> GetStudent
-- Selects the details of a single student based on their ID.
-- Used when fetching details for a specific student (by ID).
SELECT id, name, date_of_birth, address, contact, program FROM students WHERE id = $1;

-- Used in handlers.go -> UpdateStudent
-- Updates the details of an existing student record based on their ID.
UPDATE students SET name = $1, date_of_birth = $2, address = $3, contact = $4, program = $5 WHERE id = $6;

-- Used in handlers.go -> DeleteStudent
-- Deletes a student record from the students table based on their ID.
DELETE FROM students WHERE id = $1;

-- Used in handlers.go -> CreateCourse
-- Inserts a new course record into the courses table.
-- RETURNING id fetches the generated primary key ID of the new course.
INSERT INTO courses (code, title, credits) VALUES ($1, $2, $3) RETURNING id;

-- Used in handlers.go -> GetCourses
-- Selects the details of all courses from the courses table.
-- Used when any authenticated user requests the list of all courses.
SELECT id, code, title, credits FROM courses;

-- Used in handlers.go -> GetCourse
-- Selects the details of a single course based on its ID.
-- Used when fetching details for a specific course (by ID).
SELECT id, code, title, credits FROM courses WHERE id = $1;

-- Used in handlers.go -> UpdateCourse
-- Updates the details of an existing course record based on its ID.
UPDATE courses SET code = $1, title = $2, credits = $3 WHERE id = $4;

-- Used in handlers.go -> DeleteCourse
-- Deletes a course record from the courses table based on its ID.
DELETE FROM courses WHERE id = $1;

-- Used in handlers.go -> EnrollStudent
-- Checks if a student with the given ID exists.
SELECT EXISTS(SELECT 1 FROM students WHERE id = $1);

-- Used in handlers.go -> EnrollStudent
-- Checks if a course with the given ID exists.
SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1);

-- Used in handlers.go -> EnrollStudent
-- Inserts a new enrollment record into the enrollments table.
-- RETURNING id, enrollment_date fetches the generated ID and default enrollment date.
INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2) RETURNING id, enrollment_date;

-- Used in handlers.go -> GetEnrollments (Student Role)
-- Selects enrollments for a specific student based on their ID.
-- Used when a student user requests their own enrollments.
SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE student_id = $1;

-- Used in handlers.go -> GetEnrollments (Faculty Role, Filtered)
-- Selects enrollments for a specific student based on the provided student_id query parameter.
-- Used when a faculty user requests enrollments for a particular student.
SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE student_id = $1;

-- Used in handlers.go -> GetEnrollments (Faculty Role, All)
-- Selects all enrollment records from the enrollments table.
-- Used when a faculty user requests the list of all enrollments without filtering by student.
SELECT id, student_id, course_id, enrollment_date FROM enrollments;

-- Used in handlers.go -> GetEnrollment
-- Selects the details of a single enrollment based on its ID.
-- Used when fetching details for a specific enrollment (by ID).
SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE id = $1;

-- Used in handlers.go -> GetEnrollment (Authorization Check for Student)
-- Selects the student_id associated with a given enrollment ID.
-- Used to verify if a student user is authorized to view this enrollment.
SELECT student_id FROM enrollments WHERE id = $1;

-- Used in handlers.go -> DeleteEnrollment
-- Deletes an enrollment record from the enrollments table based on its ID.
DELETE FROM enrollments WHERE id = $1;

-- Used in handlers.go -> AddGrade
-- Checks if an enrollment with the given ID exists before adding a grade.
SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1);

-- Used in handlers.go -> AddGrade
-- Inserts a new grade record into the grades table.
-- RETURNING id fetches the generated primary key ID of the new grade.
INSERT INTO grades (enrollment_id, grade, semester) VALUES ($1, $2, $3) RETURNING id;

-- Used in handlers.go -> GetGrades (Student Role)
-- Selects grades for a student by joining grades with enrollments and filtering by student ID.
-- Used when a student user requests their grades.
SELECT
	g.id, g.enrollment_id, g.grade, g.semester
FROM
	grades g
JOIN
	enrollments e ON g.enrollment_id = e.id
WHERE
	e.student_id = $1;

-- Used in handlers.go -> GetGrades (Faculty Role)
-- Selects all grade records from the grades table.
-- Used when a faculty user requests the list of all grades.
SELECT id, enrollment_id, grade, semester FROM grades;

-- Used in handlers.go -> GetGrade
-- Selects the details of a single grade based on its ID.
-- Used when fetching details for a specific grade (by ID).
SELECT id, enrollment_id, grade, semester FROM grades WHERE id = $1;

-- Used in handlers.go -> GetGrade (Authorization Check for Student)
-- Selects the student_id associated with the enrollment linked to a grade.
-- Used to verify if a student user is authorized to view this grade.
SELECT student_id FROM enrollments WHERE id = $1;

-- Used in handlers.go -> UpdateGrade
-- Checks if an enrollment with the given ID exists before updating a grade.
SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1);

-- Used in handlers.go -> UpdateGrade
-- Updates the details of an existing grade record based on its ID.
UPDATE grades SET enrollment_id = $1, grade = $2, semester = $3 WHERE id = $4;

-- Used in handlers.go -> DeleteGrade
-- Deletes a grade record from the grades table based on its ID.
DELETE FROM grades WHERE id = $1;

-- Used in handlers.go -> GetStudentTranscript
-- Selects the student's ID and name for the transcript header.
SELECT id, name FROM students WHERE id = $1;

-- Used in handlers.go -> GetStudentTranscript
-- Selects course and grade information for a student's transcript.
-- Joins enrollments, courses, and left joins grades to include courses without grades.
-- Orders by semester and course code.
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
	e.student_id = $1
ORDER BY
	g.semester, c.code;

-- Used in handlers.go -> CalculateGPA
-- Calculates the GPA for a student by summing (grade * credits) for graded courses
-- and dividing by the sum of credits for those courses. Handles division by zero.
SELECT
	SUM(g.grade * c.credits) / NULLIF(SUM(c.credits), 0) AS gpa
FROM
	enrollments e
JOIN
	courses c ON e.course_id = c.id
JOIN
	grades g ON e.id = g.enrollment_id
WHERE
	e.student_id = $1 AND g.grade IS NOT NULL;

-- Used in auth_handlers.go -> Login
-- Selects the password and date of birth for a user (student or faculty) based on ID.
-- Used to verify credentials during login.
-- Note: This query retrieves plain text passwords and date of birth for insecure verification.
-- Query structure varies slightly based on role ('students' or 'faculty' table).
-- Example for student:
-- SELECT password, date_of_birth FROM students WHERE id = $1;
-- Example for faculty:
-- SELECT password, date_of_birth FROM faculty WHERE id = $1;

