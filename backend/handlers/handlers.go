package handlers

import (
  "context"
  "errors"
  "strconv"

  "backend/database"
  "backend/models"

  "github.com/gofiber/fiber/v3"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgconn"
)

// Helper function to send generic internal server error
func sendInternalServerError(c fiber.Ctx, err error) error {
  println("Internal Server Error:", err.Error())
  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An internal server error occurred."})
}

// Helper function to send generic bad request error
func sendBadRequestError(c fiber.Ctx, message string) error {
  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
}

// Helper function to send not found error
func sendNotFoundError(c fiber.Ctx, message string) error {
  return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": message})
}

// Helper function to send forbidden error
func sendForbiddenError(c fiber.Ctx, message string) error {
  return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": message})
}


func CreateStudent(c fiber.Ctx) error {
  student := new(models.Student)

  if err := c.Bind().JSON(student); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Basic validation (more robust validation should be implemented)
  if student.Name == "" {
    return sendBadRequestError(c, "Student name is required")
  }
  if student.DateOfBirth.IsZero() {
    return sendBadRequestError(c, "Student date of birth is required")
  }
  // WARNING: Password is being stored in plain text. Implement hashing!
  if student.Password == "" {
    // Consider requiring a password or using a secure initial setup flow
    // For now, allowing empty password based on schema default, but this is insecure.
  }


  query := `INSERT INTO students (name, password, date_of_birth, address, contact, program) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
  err := database.DB.QueryRow(context.Background(), query, student.Name, student.Password, student.DateOfBirth, student.Address, student.Contact, student.Program).Scan(&student.ID)
  if err != nil {
    return sendInternalServerError(c, err)
  }

  // Avoid returning password in the response
  student.Password = ""
  return c.Status(fiber.StatusCreated).JSON(student)
}

func GetStudents(c fiber.Ctx) error {
  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  // Faculty can see all students, students can only see themselves (handled by backend query)
  if userRole == "student" {
    student := models.Student{}
    query := `SELECT id, name, date_of_birth, address, contact, program FROM students WHERE id = $1`
    err := database.DB.QueryRow(context.Background(), query, userID).Scan(
      &student.ID,
      &student.Name,
      &student.DateOfBirth,
      &student.Address,
      &student.Contact,
      &student.Program,
    )

    if err == pgx.ErrNoRows {
      return sendNotFoundError(c, "Student not found")
    } else if err != nil {
      return sendInternalServerError(c, err)
    }

    student.Password = "" // Ensure password is not returned
    return c.JSON([]models.Student{student}) // Return as a list for consistency with faculty view
  } else { // Faculty
    rows, err := database.DB.Query(context.Background(), "SELECT id, name, date_of_birth, address, contact, program FROM students")
    if err != nil {
      return sendInternalServerError(c, err)
    }
    defer rows.Close()

    students := []models.Student{}
    for rows.Next() {
      student := models.Student{}
      if err := rows.Scan(
        &student.ID,
        &student.Name,
        &student.DateOfBirth,
        &student.Address,
        &student.Contact,
        &student.Program,
      ); err != nil {
        return sendInternalServerError(c, err)
      }
      student.Password = "" // Ensure password is not returned
      students = append(students, student)
    }

    if err := rows.Err(); err != nil {
      return sendInternalServerError(c, err)
    }

    return c.JSON(students)
  }
}

func GetStudent(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  // Student can only view their own details
  if userRole == "student" && id != userID {
    return sendForbiddenError(c, "Access denied. Students can only view their own details.")
  }

  student := models.Student{}
  query := `SELECT id, name, date_of_birth, address, contact, program FROM students WHERE id = $1`
  err = database.DB.QueryRow(context.Background(), query, id).Scan(
    &student.ID,
    &student.Name,
    &student.DateOfBirth,
    &student.Address,
    &student.Contact,
    &student.Program,
  )

  if err == pgx.ErrNoRows {
    return sendNotFoundError(c, "Student not found")
  } else if err != nil {
    return sendInternalServerError(c, err)
  }

  student.Password = "" // Ensure password is not returned
  return c.JSON(student)
}

func UpdateStudent(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  // Authorization check (FacultyOnly middleware is applied, but could add more granular checks here)
  // For now, assuming FacultyOnly is sufficient based on routes.

  student := new(models.Student)
  if err := c.Bind().JSON(student); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Note: This handler does not allow updating the password.
  // A separate, secure endpoint should be used for password changes.
  query := `UPDATE students SET name = $1, date_of_birth = $2, address = $3, contact = $4, program = $5 WHERE id = $6`
  result, err := database.DB.Exec(context.Background(), query, student.Name, student.DateOfBirth, student.Address, student.Contact, student.Program, id)
  if err != nil {
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Student not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student updated successfully"})
}

func DeleteStudent(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  // Authorization check (FacultyOnly middleware is applied, but could add more granular checks here)
  // For now, assuming FacultyOnly is sufficient based on routes.

  query := `DELETE FROM students WHERE id = $1`
  result, err := database.DB.Exec(context.Background(), query, id)
  if err != nil {
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Student not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student deleted successfully"})
}

func CreateCourse(c fiber.Ctx) error {
  course := new(models.Course)

  if err := c.Bind().JSON(course); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Basic validation (more robust validation should be implemented)
  if course.Code == "" || course.Title == "" || course.Credits <= 0 {
    return sendBadRequestError(c, "Course code, title, and positive credits are required")
  }

  query := `INSERT INTO courses (code, title, credits) VALUES ($1, $2, $3) RETURNING id`
  err := database.DB.QueryRow(context.Background(), query, course.Code, course.Title, course.Credits).Scan(&course.ID)
  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Course with this code already exists"})
    }
    return sendInternalServerError(c, err)
  }

  return c.Status(fiber.StatusCreated).JSON(course)
}

func GetCourses(c fiber.Ctx) error {
  // Authorization check: Determine if all users can see all courses or if filtering is needed.
  // For now, assuming all authenticated users can see all courses based on routes.
  rows, err := database.DB.Query(context.Background(), "SELECT id, code, title, credits FROM courses")
  if err != nil {
    return sendInternalServerError(c, err)
  }
  defer rows.Close()

  courses := []models.Course{}
  for rows.Next() {
    course := models.Course{}
    if err := rows.Scan(&course.ID, &course.Code, &course.Title, &course.Credits); err != nil {
      return sendInternalServerError(c, err)
    }
    courses = append(courses, course)
  }

  if err := rows.Err(); err != nil {
    return sendInternalServerError(c, err)
  }

  return c.JSON(courses)
}

func GetCourse(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid course ID")
  }

  // Authorization check: Determine if all users can see any course or if filtering is needed.
  // For now, assuming all authenticated users can see any course based on routes.

  course := models.Course{}
  query := `SELECT id, code, title, credits FROM courses WHERE id = $1`
  err = database.DB.QueryRow(context.Background(), query, id).Scan(&course.ID, &course.Code, &course.Title, &course.Credits)

  if err == pgx.ErrNoRows {
    return sendNotFoundError(c, "Course not found")
  } else if err != nil {
    return sendInternalServerError(c, err)
  }

  return c.JSON(course)
}

func UpdateCourse(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid course ID")
  }

  // Authorization check (FacultyOnly middleware is applied, but could add more granular checks here)

  course := new(models.Course)
  if err := c.Bind().JSON(course); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Basic validation
  if course.Code == "" || course.Title == "" || course.Credits <= 0 {
    return sendBadRequestError(c, "Course code, title, and positive credits are required")
  }

  query := `UPDATE courses SET code = $1, title = $2, credits = $3 WHERE id = $4`
  result, err := database.DB.Exec(context.Background(), query, course.Code, course.Title, course.Credits, id)
  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Course with this code already exists"})
    }
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Course not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course updated successfully"})
}

func DeleteCourse(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid course ID")
  }

  // Authorization check (FacultyOnly middleware is applied, but could add more granular checks here)

  query := `DELETE FROM courses WHERE id = $1`
  result, err := database.DB.Exec(context.Background(), query, id)
  if err != nil {
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Course not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course deleted successfully"})
}

func EnrollStudent(c fiber.Ctx) error {
  enrollment := new(models.Enrollment)

  if err := c.Bind().JSON(enrollment); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Basic validation
  if enrollment.StudentID == 0 || enrollment.CourseID == 0 {
    return sendBadRequestError(c, "Student ID and Course ID are required")
  }

  // Verify student and course exist
  var studentExists bool
  err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM students WHERE id = $1)", enrollment.StudentID).Scan(&studentExists)
  if err != nil || !studentExists {
    // Return generic bad request for invalid ID, avoid confirming existence of specific IDs
    return sendBadRequestError(c, "Invalid student ID or course ID")
  }

  var courseExists bool
  err = database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)", enrollment.CourseID).Scan(&courseExists)
  if err != nil || !courseExists {
    return sendBadRequestError(c, "Invalid student ID or course ID")
  }

  query := `INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2) RETURNING id, enrollment_date`
  err = database.DB.QueryRow(context.Background(), query, enrollment.StudentID, enrollment.CourseID).Scan(&enrollment.ID, &enrollment.EnrollmentDate)

  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Student is already enrolled in this course"})
    }
    return sendInternalServerError(c, err)
  }

  return c.Status(fiber.StatusCreated).JSON(enrollment)
}

func GetEnrollments(c fiber.Ctx) error {
  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  // Allow filtering by student_id for faculty
  studentIDParam := c.Query("student_id")
  var filterStudentID *int
  if studentIDParam != "" {
    id, err := strconv.Atoi(studentIDParam)
    if err != nil {
      return sendBadRequestError(c, "Invalid student_id query parameter")
    }
    filterStudentID = &id
  }


  var rows pgx.Rows
  var err error
  var query string
  var args []any

  if userRole == "student" {
    // Students can only see their own enrollments
    query = `SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE student_id = $1`
    args = []any{userID}
  } else { // Faculty
    if filterStudentID != nil {
      // Faculty requesting enrollments for a specific student
      // TODO: Add more granular authorization check here if faculty should only see enrollments for students they manage/teach
      query = `SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE student_id = $1`
      args = []any{*filterStudentID}
    } else {
      // Faculty requesting all enrollments
      // TODO: Consider if faculty should see ALL enrollments or only those related to their courses/students
      query = `SELECT id, student_id, course_id, enrollment_date FROM enrollments`
      args = []any{}
    }
  }

  rows, err = database.DB.Query(context.Background(), query, args...)
  if err != nil {
    return sendInternalServerError(c, err)
  }
  defer rows.Close()

  enrollments := []models.Enrollment{}
  for rows.Next() {
    enrollment := models.Enrollment{}
    if err := rows.Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.EnrollmentDate); err != nil {
      return sendInternalServerError(c, err)
    }
    enrollments = append(enrollments, enrollment)
  }

  if err := rows.Err(); err != nil {
    return sendInternalServerError(c, err)
  }

  return c.JSON(enrollments)
}

func GetEnrollment(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid enrollment ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  enrollment := models.Enrollment{}
  query := `SELECT id, student_id, course_id, enrollment_date FROM enrollments WHERE id = $1`
  err = database.DB.QueryRow(context.Background(), query, id).Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.EnrollmentDate)

  if err == pgx.ErrNoRows {
    return sendNotFoundError(c, "Enrollment not found")
  } else if err != nil {
    return sendInternalServerError(c, err)
  }

  // Student can only view their own enrollments
  if userRole == "student" && enrollment.StudentID != userID {
    return sendForbiddenError(c, "Access denied. Students can only view their own enrollments.")
  }

  // TODO: Add authorization check for faculty if needed (e.g., can only view enrollments for students they manage)

  return c.JSON(enrollment)
}

func DeleteEnrollment(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid enrollment ID")
  }

  // Authorization check (FacultyOnly middleware is applied)
  // TODO: Add more granular authorization check here if faculty should only delete enrollments they manage

  query := `DELETE FROM enrollments WHERE id = $1`
  result, err := database.DB.Exec(context.Background(), query, id)
  if err != nil {
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Enrollment not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Enrollment deleted successfully"})
}

func AddGrade(c fiber.Ctx) error {
  grade := new(models.Grade)

  if err := c.Bind().JSON(grade); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Basic validation
  if grade.EnrollmentID == 0 || grade.Semester == 0 { // Changed check to 0 for int
    return sendBadRequestError(c, "Enrollment ID and Semester are required")
  }
  // Grade can be NULL, so no check for grade.Grade

  // Verify enrollment exists
  var enrollmentExists bool
  err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1)", grade.EnrollmentID).Scan(&enrollmentExists)
  if err != nil || !enrollmentExists {
    // Return generic bad request for invalid ID
    return sendBadRequestError(c, "Invalid enrollment ID")
  }

  // TODO: Add authorization check here - ensure faculty user is authorized to add a grade for this enrollment

  query := `INSERT INTO grades (enrollment_id, grade, semester) VALUES ($1, $2, $3) RETURNING id`
  err = database.DB.QueryRow(context.Background(), query, grade.EnrollmentID, grade.Grade, grade.Semester).Scan(&grade.ID)
  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Grade for this enrollment and semester already exists"})
    }
    return sendInternalServerError(c, err)
  }

  return c.Status(fiber.StatusCreated).JSON(grade)
}

func GetGrades(c fiber.Ctx) error {
  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  var rows pgx.Rows
  var err error
  var query string
  var args []any

  if userRole == "student" {
    // Students can only see grades for their own enrollments
    query = `
    SELECT
      g.id, g.enrollment_id, g.grade, g.semester
    FROM
      grades g
    JOIN
      enrollments e ON g.enrollment_id = e.id
    WHERE
      e.student_id = $1
    `
    args = []any{userID}
  } else { // Faculty
    // TODO: Consider if faculty should see ALL grades or only those related to their courses/students
    query = `SELECT id, enrollment_id, grade, semester FROM grades`
    args = []any{}
  }

  rows, err = database.DB.Query(context.Background(), query, args...)
  if err != nil {
    return sendInternalServerError(c, err)
  }
  defer rows.Close()

  grades := []models.Grade{}
  for rows.Next() {
    grade := models.Grade{}
    if err := rows.Scan(&grade.ID, &grade.EnrollmentID, &grade.Grade, &grade.Semester); err != nil {
      return sendInternalServerError(c, err)
    }
    grades = append(grades, grade)
  }

  if err := rows.Err(); err != nil {
    return sendInternalServerError(c, err)
  }

  return c.JSON(grades)
}

func GetGrade(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid grade ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  grade := models.Grade{}
  query := `SELECT id, enrollment_id, grade, semester FROM grades WHERE id = $1`
  err = database.DB.QueryRow(context.Background(), query, id).Scan(&grade.ID, &grade.EnrollmentID, &grade.Grade, &grade.Semester)

  if err == pgx.ErrNoRows {
    return sendNotFoundError(c, "Grade not found")
  } else if err != nil {
    return sendInternalServerError(c, err)
  }

  if userRole == "student" {
    // Students can only view grades for their own enrollments
    var studentID int
    enrollmentQuery := `SELECT student_id FROM enrollments WHERE id = $1`
    err = database.DB.QueryRow(context.Background(), enrollmentQuery, grade.EnrollmentID).Scan(&studentID)
    if err != nil {
      // If enrollment not found for the grade's enrollment_id, something is wrong
      return sendInternalServerError(c, errors.New("Failed to verify enrollment for grade: " + err.Error()))
    }
    if studentID != userID {
      return sendForbiddenError(c, "Access denied. Students can only view grades for their own enrollments.")
    }
  }

  // TODO: Add authorization check for faculty if needed (e.g., can only view grades for students they manage)

  return c.JSON(grade)
}

func UpdateGrade(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid grade ID")
  }

  // Authorization check (FacultyOnly middleware is applied)
  // TODO: Add more granular authorization check here - ensure faculty user is authorized to update this specific grade

  grade := new(models.Grade)
  if err := c.Bind().JSON(grade); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  // Basic validation
  if grade.EnrollmentID == 0 || grade.Semester == 0 { // Changed check to 0 for int
    return sendBadRequestError(c, "Enrollment ID and Semester are required")
  }
  // Grade can be NULL, so no check for grade.Grade

  // Verify enrollment exists (optional, as grade references enrollment, but good practice)
  var enrollmentExists bool
  err = database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1)", grade.EnrollmentID).Scan(&enrollmentExists)
  if err != nil || !enrollmentExists {
    // Return generic bad request for invalid ID
    return sendBadRequestError(c, "Invalid enrollment ID")
  }


  query := `UPDATE grades SET enrollment_id = $1, grade = $2, semester = $3 WHERE id = $4`
  result, err := database.DB.Exec(context.Background(), query, grade.EnrollmentID, grade.Grade, grade.Semester, id)
  if err != nil {
    if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
      return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Grade for this enrollment and semester already exists"})
    }
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Grade not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Grade updated successfully"})
}

func DeleteGrade(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid grade ID")
  }

  // Authorization check (FacultyOnly middleware is applied)
  // TODO: Add more granular authorization check here - ensure faculty user is authorized to delete this specific grade

  query := `DELETE FROM grades WHERE id = $1`
  result, err := database.DB.Exec(context.Background(), query, id)
  if err != nil {
    return sendInternalServerError(c, err)
  }

  rowsAffected := result.RowsAffected()

  if rowsAffected == 0 {
    return sendNotFoundError(c, "Grade not found")
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Grade deleted successfully"})
}

func GetStudentTranscript(c fiber.Ctx) error {
  studentID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  // Student can only view their own transcript
  if userRole == "student" && studentID != userID {
    return sendForbiddenError(c, "Access denied. Students can only view their own transcript.")
  }

  // TODO: Add authorization check for faculty if needed (e.g., can only view transcripts for students they manage)

  student := models.Student{}
  studentQuery := `SELECT id, name FROM students WHERE id = $1`
  err = database.DB.QueryRow(context.Background(), studentQuery, studentID).Scan(&student.ID, &student.Name)
  if err == pgx.ErrNoRows {
    return sendNotFoundError(c, "Student not found")
  } else if err != nil {
    return sendInternalServerError(c, errors.New("Failed to fetch student details for transcript: " + err.Error()))
  }

  transcriptQuery := `
  SELECT
    e.id AS enrollment_id, -- Added enrollment_id
    c.code,
    c.title,
    c.credits,
    g.id AS grade_id, -- Added grade_id
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
    g.semester, c.code -- Order by semester (int)
  `
  rows, err := database.DB.Query(context.Background(), transcriptQuery, studentID)
  if err != nil {
    return sendInternalServerError(c, errors.New("Failed to fetch transcript data: " + err.Error()))
  }
  defer rows.Close()

  transcriptCourses := []models.TranscriptCourse{}
  for rows.Next() {
    course := models.TranscriptCourse{}
    // Scan into the new fields
    err := rows.Scan(&course.EnrollmentID, &course.CourseCode, &course.CourseTitle, &course.Credits, &course.GradeID, &course.Grade, &course.Semester)
    if err != nil {
      return sendInternalServerError(c, errors.New("Failed to scan transcript course data: " + err.Error()))
    }
    transcriptCourses = append(transcriptCourses, course)
  }

  if err := rows.Err(); err != nil {
    return sendInternalServerError(c, errors.New("Error during transcript row iteration: " + err.Error()))
  }

  transcript := models.StudentTranscript{
    StudentID:   student.ID,
    StudentName: student.Name,
    Courses:     transcriptCourses,
  }

  return c.JSON(transcript)
}

func CalculateGPA(c fiber.Ctx) error {
  studentID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  // Student can only calculate their own GPA
  if userRole == "student" && studentID != userID {
    return sendForbiddenError(c, "Access denied. Students can only calculate their own GPA.")
  }

  // TODO: Add authorization check for faculty if needed

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
  var gpa *float64
  err = database.DB.QueryRow(context.Background(), query, studentID).Scan(&gpa)

  if err == pgx.ErrNoRows || gpa == nil {
    // Return 0.0 GPA if no grades or no graded courses
    return c.Status(fiber.StatusOK).JSON(fiber.Map{"student_id": studentID, "gpa": 0.0, "message": "No grades available to calculate GPA"})
  } else if err != nil {
    return sendInternalServerError(c, errors.New("Failed to calculate GPA: " + err.Error()))
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"student_id": studentID, "gpa": *gpa})
}
