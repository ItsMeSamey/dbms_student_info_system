package handlers

import (
  "context"
  "errors"
  "strconv"
  "time"

  "backend/database"
  "backend/models"

  "github.com/gofiber/fiber/v3"
  "github.com/jackc/pgx/v5/pgconn"
)

func sendInternalServerError(c fiber.Ctx, err error) error {
  println("Internal Server Error:", err.Error())
  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An internal server error occurred."})
}

func sendBadRequestError(c fiber.Ctx, message string) error {
  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
}

func handleDatabaseError(c fiber.Ctx, err error) error {
  println("Database Error:", err.Error())
  if pgErr, ok := err.(*pgconn.PgError); ok {
    switch pgErr.Code {
    case "P0001":
      switch pgErr.Message {
      case "Invalid credentials":
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
      case "Student not found", "Course not found", "Enrollment not found", "Grade not found":
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": pgErr.Message})
      case "Access denied. Invalid user role.",
        "Access denied. Students can only view their own details.",
        "Access denied. Only faculty can update student details.",
        "Access denied. Only faculty can create courses.",
        "Access denied. Only faculty can update courses.",
        "Access denied. Only faculty can delete courses.",
        "Access denied. Only faculty can create enrollments.",
        "Access denied. Students can only view their own enrollments.",
        "Access denied. Only faculty can delete enrollments.",
        "Access denied. Only faculty can add grades.",
        "Access denied. Students can only view grades for their own enrollments.",
        "Access denied. Only faculty can update grades.",
        "Access denied. Only faculty can delete grades.",
        "Access denied. Students can only view their own transcript.",
        "Access denied. Students can only calculate their own GPA.":
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": pgErr.Message})
      case "Student name is required", "Student date of birth is required",
        "Course code is required", "Course title is required", "Positive credits are required",
        "Student ID and Course ID are required", "Invalid student ID or course ID",
        "Enrollment ID and Semester are required", "Invalid enrollment ID":
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": pgErr.Message})
      case "Course with code already exists", "Student is already enrolled in this course", "Grade for this enrollment and semester already exists":
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": pgErr.Message})
      default:
        return sendInternalServerError(c, errors.New("Database operation failed"))
      }
    case "23505":
      return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Duplicate entry violates unique constraint"})
    default:
      return sendInternalServerError(c, errors.New("Database error"))
    }
  }
  return sendInternalServerError(c, errors.New("An unexpected error occurred"))
}

func CreateStudent(c fiber.Ctx) error {
  student := new(models.Student)

  if err := c.Bind().JSON(student); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  if student.Name == "" {
    return sendBadRequestError(c, "Student name is required")
  }
  if student.DateOfBirth.IsZero() {
    return sendBadRequestError(c, "Student date of birth is required")
  }

  var newStudentID int
  query := `SELECT create_student($1, $2, $3, $4, $5, $6)`
  err := database.DB.QueryRow(context.Background(), query,
    student.Name,
    student.Password,
    student.DateOfBirth,
    student.Address,
    student.Contact,
    student.Program,
  ).Scan(&newStudentID)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)
  createdStudent := models.Student{}
  getStudentQuery := `SELECT id, name, date_of_birth, address, contact, program FROM get_student_by_id($1, $2, $3)`
  err = database.DB.QueryRow(context.Background(), getStudentQuery, newStudentID, userID, userRole).Scan(
    &createdStudent.ID,
    &createdStudent.Name,
    &createdStudent.DateOfBirth,
    &createdStudent.Address,
    &createdStudent.Contact,
    &createdStudent.Program,
  )
  if err != nil {
    return sendInternalServerError(c, errors.New("Failed to retrieve created student"))
  }

  createdStudent.Password = ""
  return c.Status(fiber.StatusCreated).JSON(createdStudent)
}

func GetStudents(c fiber.Ctx) error {
  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  query := `SELECT id, name, password, date_of_birth, address, contact, program FROM get_students($1, $2)`
  rows, err := database.DB.Query(context.Background(), query, userID, userRole)
  if err != nil {
    return handleDatabaseError(c, err)
  }
  defer rows.Close()

  students := []models.Student{}
  for rows.Next() {
    student := models.Student{}
    var password string
    if err := rows.Scan(
      &student.ID,
      &student.Name,
      &password,
      &student.DateOfBirth,
      &student.Address,
      &student.Contact,
      &student.Program,
    ); err != nil {
      return sendInternalServerError(c, err)
    }
    student.Password = ""
    students = append(students, student)
  }

  if err := rows.Err(); err != nil {
    return sendInternalServerError(c, err)
  }

  return c.JSON(students)
}

func GetStudent(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  student := models.Student{}
  query := `SELECT id, name, password, date_of_birth, address, contact, program FROM get_student_by_id($1, $2, $3)`
  var password string
  err = database.DB.QueryRow(context.Background(), query, id, userID, userRole).Scan(
    &student.ID,
    &student.Name,
    &password,
    &student.DateOfBirth,
    &student.Address,
    &student.Contact,
    &student.Program,
  )

  if err != nil {
    return handleDatabaseError(c, err)
  }

  student.Password = ""
  return c.JSON(student)
}

func UpdateStudent(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  student := new(models.Student)
  if err := c.Bind().JSON(student); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  if student.Name == "" {
    return sendBadRequestError(c, "Student name is required")
  }
  if student.DateOfBirth.IsZero() {
    return sendBadRequestError(c, "Student date of birth is required")
  }

  query := `CALL update_student($1, $2, $3, $4, $5, $6, $7, $8)`
  _, err = database.DB.Exec(context.Background(), query,
    id,
    student.Name,
    student.DateOfBirth,
    student.Address,
    student.Contact,
    student.Program,
    userID,
    userRole,
  )

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student updated successfully"})
}

func DeleteStudent(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid student ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  query := `CALL delete_student($1, $2, $3)`
  _, err = database.DB.Exec(context.Background(), query, id, userID, userRole)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Student deleted successfully"})
}

func CreateCourse(c fiber.Ctx) error {
  course := new(models.Course)

  if err := c.Bind().JSON(course); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  if course.Code == "" || course.Title == "" || course.Credits <= 0 {
    return sendBadRequestError(c, "Course code, title, and positive credits are required")
  }

  var newCourseID int
  query := `SELECT create_course($1, $2, $3, $4, $5)`
  err := database.DB.QueryRow(context.Background(), query,
    course.Code,
    course.Title,
    course.Credits,
    userID,
    userRole,
  ).Scan(&newCourseID)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  createdCourse := models.Course{}
  getCourseQuery := `SELECT id, code, title, credits FROM get_course_by_id($1)`
  err = database.DB.QueryRow(context.Background(), getCourseQuery, newCourseID).Scan(
    &createdCourse.ID,
    &createdCourse.Code,
    &createdCourse.Title,
    &createdCourse.Credits,
  )
  if err != nil {
    return sendInternalServerError(c, errors.New("Failed to retrieve created course"))
  }

  return c.Status(fiber.StatusCreated).JSON(createdCourse)
}

func GetCourses(c fiber.Ctx) error {
  query := `SELECT id, code, title, credits FROM get_all_courses()`
  rows, err := database.DB.Query(context.Background(), query)
  if err != nil {
    return handleDatabaseError(c, err)
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

  course := models.Course{}
  query := `SELECT id, code, title, credits FROM get_course_by_id($1)`
  err = database.DB.QueryRow(context.Background(), query, id).Scan(&course.ID, &course.Code, &course.Title, &course.Credits)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.JSON(course)
}

func UpdateCourse(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid course ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  course := new(models.Course)
  if err := c.Bind().JSON(course); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  if course.Code == "" || course.Title == "" || course.Credits <= 0 {
    return sendBadRequestError(c, "Course code, title, and positive credits are required")
  }

  query := `CALL update_course($1, $2, $3, $4, $5, $6)`
  _, err = database.DB.Exec(context.Background(), query,
    id,
    course.Code,
    course.Title,
    course.Credits,
    userID,
    userRole,
  )

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course updated successfully"})
}

func DeleteCourse(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid course ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  query := `CALL delete_course($1, $2, $3)`
  _, err = database.DB.Exec(context.Background(), query, id, userID, userRole)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Course deleted successfully"})
}

func EnrollStudent(c fiber.Ctx) error {
  enrollment := new(models.Enrollment)

  if err := c.Bind().JSON(enrollment); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  if enrollment.StudentID == 0 || enrollment.CourseID == 0 {
    return sendBadRequestError(c, "Student ID and Course ID are required")
  }

  var newEnrollmentID int
  var enrollmentDate time.Time
  query := `SELECT v_id, v_date FROM create_enrollment($1, $2, $3, $4)`
  err := database.DB.QueryRow(context.Background(), query,
    enrollment.StudentID,
    enrollment.CourseID,
    userID,
    userRole,
  ).Scan(&newEnrollmentID, &enrollmentDate)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  createdEnrollment := models.Enrollment{
    ID: newEnrollmentID,
    StudentID: enrollment.StudentID,
    CourseID: enrollment.CourseID,
    EnrollmentDate: enrollmentDate,
  }

  return c.Status(fiber.StatusCreated).JSON(createdEnrollment)
}

func GetEnrollments(c fiber.Ctx) error {
  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  studentIDParam := c.Query("student_id")
  var filterStudentID *int
  if studentIDParam != "" {
    id, err := strconv.Atoi(studentIDParam)
    if err != nil {
      return sendBadRequestError(c, "Invalid student_id query parameter")
    }
    filterStudentID = &id
  }

  query := `SELECT id, student_id, course_id, enrollment_date FROM get_enrollments($1, $2, $3)`
  rows, err := database.DB.Query(context.Background(), query, userID, userRole, filterStudentID)

  if err != nil {
    return handleDatabaseError(c, err)
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
  query := `SELECT id, student_id, course_id, enrollment_date FROM get_enrollment_by_id($1, $2, $3)`
  err = database.DB.QueryRow(context.Background(), query, id, userID, userRole).Scan(
    &enrollment.ID,
    &enrollment.StudentID,
    &enrollment.CourseID,
    &enrollment.EnrollmentDate,
  )

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.JSON(enrollment)
}

func DeleteEnrollment(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid enrollment ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  query := `CALL delete_enrollment($1, $2, $3)`
  _, err = database.DB.Exec(context.Background(), query, id, userID, userRole)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Enrollment deleted successfully"})
}

func AddGrade(c fiber.Ctx) error {
  grade := new(models.Grade)

  if err := c.Bind().JSON(grade); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  if grade.EnrollmentID == 0 || grade.Semester == 0 {
    return sendBadRequestError(c, "Enrollment ID and Semester are required")
  }

  var newGradeID int
  query := `SELECT add_grade($1, $2, $3, $4, $5)`
  err := database.DB.QueryRow(context.Background(), query,
    grade.EnrollmentID,
    grade.Grade,
    grade.Semester,
    userID,
    userRole,
  ).Scan(&newGradeID)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  createdGrade := models.Grade{}
  getGradeQuery := `SELECT id, enrollment_id, grade, semester FROM get_grade_by_id($1, $2, $3)`
  err = database.DB.QueryRow(context.Background(), getGradeQuery, newGradeID, userID, userRole).Scan(
    &createdGrade.ID,
    &createdGrade.EnrollmentID,
    &createdGrade.Grade,
    &createdGrade.Semester,
  )
  if err != nil {
    return sendInternalServerError(c, errors.New("Failed to retrieve created grade"))
  }

  return c.Status(fiber.StatusCreated).JSON(createdGrade)
}

func GetGrades(c fiber.Ctx) error {
  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  query := `SELECT id, enrollment_id, grade, semester FROM get_all_grades($1, $2)`
  rows, err := database.DB.Query(context.Background(), query, userID, userRole)

  if err != nil {
    return handleDatabaseError(c, err)
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
  query := `SELECT id, enrollment_id, grade, semester FROM get_grade_by_id($1, $2, $3)`
  err = database.DB.QueryRow(context.Background(), query, id, userID, userRole).Scan(
    &grade.ID,
    &grade.EnrollmentID,
    &grade.Grade,
    &grade.Semester,
  )

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.JSON(grade)
}

func UpdateGrade(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid grade ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  grade := new(models.Grade)
  if err := c.Bind().JSON(grade); err != nil {
    return sendBadRequestError(c, "Invalid request body")
  }

  if grade.EnrollmentID == 0 || grade.Semester == 0 {
    return sendBadRequestError(c, "Enrollment ID and Semester are required")
  }

  query := `CALL update_grade($1, $2, $3, $4, $5, $6)`
  _, err = database.DB.Exec(context.Background(), query,
    id,
    grade.EnrollmentID,
    grade.Grade,
    grade.Semester,
    userID,
    userRole,
  )

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Grade updated successfully"})
}

func DeleteGrade(c fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return sendBadRequestError(c, "Invalid grade ID")
  }

  userRole := c.Locals("userRole").(string)
  userID := c.Locals("userID").(int)

  query := `CALL delete_grade($1, $2, $3)`
  _, err = database.DB.Exec(context.Background(), query, id, userID, userRole)

  if err != nil {
    return handleDatabaseError(c, err)
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

  query := `SELECT enrollment_id, course_code, course_title, credits, grade_id, grade, semester FROM get_student_transcript($1, $2, $3)`
  rows, err := database.DB.Query(context.Background(), query, studentID, userID, userRole)
  if err != nil {
    return handleDatabaseError(c, err)
  }
  defer rows.Close()

  transcriptCourses := []models.TranscriptCourse{}
  for rows.Next() {
    course := models.TranscriptCourse{}
    var gradeID *int
    var grade *float64
    var semester *int

    err := rows.Scan(
      &course.EnrollmentID,
      &course.CourseCode,
      &course.CourseTitle,
      &course.Credits,
      &gradeID,
      &grade,
      &semester,
    )
    if err != nil {
      return sendInternalServerError(c, err)
    }

    course.GradeID = gradeID
    course.Grade = grade
    course.Semester = semester

    transcriptCourses = append(transcriptCourses, course)
  }

  if err := rows.Err(); err != nil {
    return sendInternalServerError(c, err)
  }

  student := models.Student{}
  getStudentQuery := `SELECT id, name FROM get_student_by_id($1, $2, $3)`
  err = database.DB.QueryRow(context.Background(), getStudentQuery, studentID, userID, userRole).Scan(&student.ID, &student.Name)
  if err != nil {
    return handleDatabaseError(c, err)
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

  var gpa float64
  query := `SELECT calculate_student_gpa($1, $2, $3)`
  err = database.DB.QueryRow(context.Background(), query, studentID, userID, userRole).Scan(&gpa)

  if err != nil {
    return handleDatabaseError(c, err)
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"student_id": studentID, "gpa": gpa})
}

