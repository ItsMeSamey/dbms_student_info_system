-- Used in handlers.go -> CreateStudent
-- Calls the PL/SQL function to create a new student.
SELECT create_student($1, $2, $3, $4, $5, $6);

-- Used in handlers.go -> GetStudents
-- Calls the PL/SQL function to get student(s) based on user role and ID.
SELECT id, name, password, date_of_birth, address, contact, program FROM get_students($1, $2);

-- Used in handlers.go -> GetStudent
-- Calls the PL/SQL function to get a single student by ID, with authorization check.
SELECT id, name, password, date_of_birth, address, contact, program FROM get_student_by_id($1, $2, $3);

-- Used in handlers.go -> UpdateStudent
-- Calls the PL/SQL procedure to update a student, with authorization check.
CALL update_student($1, $2, $3, $4, $5, $6, $7, $8);

-- Used in handlers.go -> DeleteStudent
-- Calls the PL/SQL procedure to delete a student, with authorization check.
CALL delete_student($1, $2, $3);

-- Used in handlers.go -> CreateCourse
-- Calls the PL/SQL function to create a new course, with authorization check.
SELECT create_course($1, $2, $3, $4, $5);

-- Used in handlers.go -> GetCourses
-- Calls the PL/SQL function to get all courses.
SELECT id, code, title, credits FROM get_all_courses();

-- Used in handlers.go -> GetCourse
-- Calls the PL/SQL function to get a single course by ID.
SELECT id, code, title, credits FROM get_course_by_id($1);

-- Used in handlers.go -> UpdateCourse
-- Calls the PL/SQL procedure to update a course, with authorization check.
CALL update_course($1, $2, $3, $4, $5, $6);

-- Used in handlers.go -> DeleteCourse
-- Calls the PL/SQL procedure to delete a course, with authorization check.
CALL delete_course($1, $2, $3);

-- Used in handlers.go -> EnrollStudent
-- Calls the PL/SQL function to create a new enrollment, with authorization check.
-- The function returns the new enrollment's ID and date.
SELECT enrollment_id, enrollment_date FROM create_enrollment($1, $2, $3, $4);

-- Used in handlers.go -> GetEnrollments
-- Calls the PL/SQL function to get enrollments based on user role and optional student ID filter.
SELECT id, student_id, course_id, enrollment_date FROM get_enrollments($1, $2, $3);

-- Used in handlers.go -> GetEnrollment
-- Calls the PL/SQL function to get a single enrollment by ID, with authorization check.
SELECT id, student_id, course_id, enrollment_date FROM get_enrollment_by_id($1, $2, $3);

-- Used in handlers.go -> DeleteEnrollment
-- Calls the PL/SQL procedure to delete an enrollment, with authorization check.
CALL delete_enrollment($1, $2, $3);

-- Used in handlers.go -> AddGrade
-- Calls the PL/SQL function to add a grade, with authorization check.
SELECT add_grade($1, $2, $3, $4, $5);

-- Used in handlers.go -> GetGrades
-- Calls the PL/SQL function to get grades based on user role.
SELECT id, enrollment_id, grade, semester FROM get_all_grades($1, $2);

-- Used in handlers.go -> GetGrade
-- Calls the PL/SQL function to get a single grade by ID, with authorization check.
SELECT id, enrollment_id, grade, semester FROM get_grade_by_id($1, $2, $3);

-- Used in handlers.go -> UpdateGrade
-- Calls the PL/SQL procedure to update a grade, with authorization check.
CALL update_grade($1, $2, $3, $4, $5, $6);

-- Used in handlers.go -> DeleteGrade
-- Calls the PL/SQL procedure to delete a grade, with authorization check.
CALL delete_grade($1, $2, $3);

-- Used in handlers.go -> GetStudentTranscript
-- Calls the PL/SQL function to get a student's transcript, with authorization check.
SELECT enrollment_id, course_code, course_title, credits, grade_id, grade, semester FROM get_student_transcript($1, $2, $3);

-- Used in handlers.go -> CalculateGPA
-- Calls the PL/SQL function to calculate a student's GPA, with authorization check.
SELECT calculate_student_gpa($1, $2, $3);

-- Used in auth_handlers.go -> Login
-- Calls the PL/SQL function to authenticate user credentials.
SELECT user_id, user_role FROM authenticate_user($1, $2, $3);

