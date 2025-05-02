package handlers

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"backend/database"
	"backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/gofiber/fiber/v3"
)

func parseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

func CreateStudent(c fiber.Ctx) error {
	student := new(models.Student)

	if err := c.Bind().Body(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if student.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Student name is required"})
	}

	query := `INSERT INTO students (name, date_of_birth, address, contact, program) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := database.DB.QueryRow(context.Background(), query, student.Name, student.DateOfBirth, student.Address, student.Contact, student.Program).Scan(&student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create student", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(student)
}

func GetStudents(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id, name, date_of_birth, address, contact, program FROM students")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch students", "details": err.Error()})
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		student := models.Student{}
		var dateOfBirth sql.NullTime
		var address, contact, program sql.NullString
		if err := rows.Scan(&student.ID, &student.Name, &dateOfBirth, &address, &contact, &program); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan student data", "details": err.Error()})
		}
		if dateOfBirth.Valid {
			student.DateOfBirth = dateOfBirth.Time
		}
		if address.Valid {
			student.Address = address.String
		}
		if contact.Valid {
			student.Contact = contact.String
		}
		if program.Valid {
			student.Program = program.String
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error during row iteration", "details": err.Error()})
	}

	return c.JSON(students)
}

func GetStudent(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	student := models.Student{}
	var dateOfBirth sql.NullTime
	var address, contact, program sql.NullString

	query := `SELECT id, name, date_of_birth, address, contact, program FROM students WHERE id = $1`
	err = database.DB.QueryRow(context.Background(), query, id).Scan(&student.ID, &student.Name, &dateOfBirth, &address, &contact, &program)

	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch student", "details": err.Error()})
	}

	if dateOfBirth.Valid {
		student.DateOfBirth = dateOfBirth.Time
	}
	if address.Valid {
		student.Address = address.String
	}
	if contact.Valid {
		student.Contact = contact.String
	}
	if program.Valid {
		student.Program = program.String
	}

	return c.JSON(student)
}

func UpdateStudent(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	student := new(models.Student)
	if err := c.Bind().Body(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `UPDATE students SET name = $1, date_of_birth = $2, address = $3, contact = $4, program = $5 WHERE id = $6`
	result, err := database.DB.Exec(context.Background(), query, student.Name, student.DateOfBirth, student.Address, student.Contact, student.Program, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update student", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student updated successfully"})
}

func DeleteStudent(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	query := `DELETE FROM students WHERE id = $1`
	result, err := database.DB.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete student", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student deleted successfully"})
}

func CreateCourse(c fiber.Ctx) error {
	course := new(models.Course)

	if err := c.Bind().Body(course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if course.Code == "" || course.Title == "" || course.Credits <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Course code, title, and positive credits are required"})
	}

	query := `INSERT INTO courses (code, title, credits) VALUES ($1, $2, $3) RETURNING id`
	err := database.DB.QueryRow(context.Background(), query, course.Code, course.Title, course.Credits).Scan(&course.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create course", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(course)
}

func GetCourses(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id, code, title, credits FROM courses")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch courses", "details": err.Error()})
	}
	defer rows.Close()

	courses := []models.Course{}
	for rows.Next() {
		course := models.Course{}
		if err := rows.Scan(&course.ID, &course.Code, &course.Title, &course.Credits); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan course data", "details": err.Error()})
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error during row iteration", "details": err.Error()})
	}

	return c.JSON(courses)
}

func GetCourse(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	course := models.Course{}
	query := `SELECT id, code, title, credits FROM courses WHERE id = $1`
	err = database.DB.QueryRow(context.Background(), query, id).Scan(&course.ID, &course.Code, &course.Title, &course.Credits)

	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch course", "details": err.Error()})
	}

	return c.JSON(course)
}

func UpdateCourse(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	course := new(models.Course)
	if err := c.Bind().Body(course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if course.Code == "" || course.Title == "" || course.Credits <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Course code, title, and positive credits are required"})
	}

	query := `UPDATE courses SET code = $1, title = $2, credits = $3 WHERE id = $4`
	result, err := database.DB.Exec(context.Background(), query, course.Code, course.Title, course.Credits, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update course", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course updated successfully"})
}

func DeleteCourse(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	query := `DELETE FROM courses WHERE id = $1`
	result, err := database.DB.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete course", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course deleted successfully"})
}

func EnrollStudent(c fiber.Ctx) error {
	enrollment := new(models.Enrollment)

	if err := c.Bind().Body(enrollment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if enrollment.StudentID == 0 || enrollment.CourseID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Student ID and Course ID are required"})
	}

	var studentExists bool
	err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", enrollment.StudentID).Scan(&studentExists)
	if err != nil || !studentExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	var courseExists bool
	err = database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", enrollment.CourseID).Scan(&courseExists)
	if err != nil || !courseExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	query := `INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2) RETURNING id, enrollment_date`
	err = database.DB.QueryRow(context.Background(), query, enrollment.StudentID, enrollment.CourseID).Scan(&enrollment.ID, &enrollment.EnrollmentDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create enrollment", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(enrollment)
}

func GetEnrollments(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id, student_id, course_id, enrollment_date FROM enrollments")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch enrollments", "details": err.Error()})
	}
	defer rows.Close()

	enrollments := []models.Enrollment{}
	for rows.Next() {
		enrollment := models.Enrollment{}
		if err := rows.Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.EnrollmentDate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan enrollment data", "details": err.Error()})
		}
		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error during row iteration", "details": err.Error()})
	}

	return c.JSON(enrollments)
}

func GetEnrollment(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid enrollment ID"})
	}

	enrollment := models.Enrollment{}
	query := `SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE id = $1`
	err = database.DB.QueryRow(context.Background(), query, id).Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.EnrollmentDate)

	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Enrollment not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch enrollment", "details": err.Error()})
	}

	return c.JSON(enrollment)
}

func DeleteEnrollment(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid enrollment ID"})
	}

	query := `DELETE FROM enrollments WHERE id = $1`
	result, err := database.DB.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete enrollment", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Enrollment not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Enrollment deleted successfully"})
}

func AddGrade(c fiber.Ctx) error {
	grade := new(models.Grade)

	if err := c.Bind().Body(grade); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if grade.EnrollmentID == 0 || grade.Semester == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Enrollment ID and Semester are required"})
	}

	var enrollmentExists bool
	err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1)", grade.EnrollmentID).Scan(&enrollmentExists)
	if err != nil || !enrollmentExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid enrollment ID"})
	}

	query := `INSERT INTO grades (enrollment_id, grade, semester) VALUES ($1, $2, $3) RETURNING id`
	err = database.DB.QueryRow(context.Background(), query, grade.EnrollmentID, grade.Grade, grade.Semester).Scan(&grade.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add grade", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(grade)
}

func GetGrades(c fiber.Ctx) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id, enrollment_id, grade, semester FROM grades")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch grades", "details": err.Error()})
	}
	defer rows.Close()

	grades := []models.Grade{}
	for rows.Next() {
		grade := models.Grade{}
		var gradeValue sql.NullFloat64
		if err := rows.Scan(&grade.ID, &grade.EnrollmentID, &gradeValue, &grade.Semester); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan grade data", "details": err.Error()})
		}
		if gradeValue.Valid {
			grade.Grade = gradeValue.Float64
		}

		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error during row iteration", "details": err.Error()})
	}

	return c.JSON(grades)
}

func GetGrade(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid grade ID"})
	}

	grade := models.Grade{}
	var gradeValue sql.NullFloat64
	query := `SELECT id, enrollment_id, grade, semester FROM grades WHERE id = $1`
	err = database.DB.QueryRow(context.Background(), query, id).Scan(&grade.ID, &grade.EnrollmentID, &gradeValue, &grade.Semester)

	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Grade not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch grade", "details": err.Error()})
	}

	if gradeValue.Valid {
		grade.Grade = gradeValue.Float64
	}

	return c.JSON(grade)
}

func UpdateGrade(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid grade ID"})
	}

	grade := new(models.Grade)
	if err := c.Bind().Body(grade); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if grade.EnrollmentID == 0 || grade.Semester == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Enrollment ID and Semester are required"})
	}

	var enrollmentExists bool
	err = database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1)", grade.EnrollmentID).Scan(&enrollmentExists)
	if err != nil || !enrollmentExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid enrollment ID"})
	}

	query := `UPDATE grades SET enrollment_id = $1, grade = $2, semester = $3 WHERE id = $4`
	result, err := database.DB.Exec(context.Background(), query, grade.EnrollmentID, grade.Grade, grade.Semester, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update grade", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Grade not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Grade updated successfully"})
}

func DeleteGrade(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid grade ID"})
	}

	query := `DELETE FROM grades WHERE id = $1`
	result, err := database.DB.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete grade", "details": err.Error()})
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Grade not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Grade deleted successfully"})
}

func GetStudentTranscript(c fiber.Ctx) error {
	studentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	student := models.Student{}
	var dateOfBirth sql.NullTime
	var address, contact, program sql.NullString
	studentQuery := `SELECT id, name, date_of_birth, address, contact, program FROM students WHERE id = $1`
	err = database.DB.QueryRow(context.Background(), studentQuery, studentID).Scan(&student.ID, &student.Name, &dateOfBirth, &address, &contact, &program)
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch student details", "details": err.Error()})
	}

	transcriptQuery := `
		SELECT
			c.code,
			c.title,
			c.credits,
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
			g.semester, c.code
	`
	rows, err := database.DB.Query(context.Background(), transcriptQuery, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch transcript data", "details": err.Error()})
	}
	defer rows.Close()

	transcriptCourses := []models.TranscriptCourse{}
	for rows.Next() {
		course := models.TranscriptCourse{}
		var grade sql.NullFloat64
		var semester sql.NullString

		if err := rows.Scan(&course.CourseCode, &course.CourseTitle, &course.Credits, &grade, &semester); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan transcript course data", "details": err.Error()})
		}

		if grade.Valid {
			course.Grade = grade.Float64
		} else {
			course.Grade = -1
		}
		if semester.Valid {
			course.Semester = semester.String
		} else {
			course.Semester = "N/A"
		}

		transcriptCourses = append(transcriptCourses, course)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error during transcript row iteration", "details": err.Error()})
	}

	transcript := models.StudentTranscript{
		StudentID: student.ID,
		StudentName: student.Name,
		Courses: transcriptCourses,
	}

	return c.JSON(transcript)
}

func CalculateGPA(c fiber.Ctx) error {
	studentID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	query := `
		SELECT
			SUM(g.grade * c.credits) / NULLIF(SUM(c.credits), 0) AS gpa
		FROM
			enrollments e
		JOIN
			courses c ON e.course_id = c.id
		JOIN
			grades g ON e.id = g.enrollment_id
		WHERE
			e.student_id = $1 AND g.grade IS NOT NULL
	`
	var gpa sql.NullFloat64
	err = database.DB.QueryRow(context.Background(), query, studentID).Scan(&gpa)

	if err == pgx.ErrNoRows || !gpa.Valid {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"student_id": studentID, "gpa": 0.0, "message": "No grades available to calculate GPA"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate GPA", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"student_id": studentID, "gpa": gpa.Float64})
}
