package routes

import (
  "backend/handlers"
  "backend/middleware"

  "github.com/gofiber/fiber/v3"
)

func SetupRoutes(app fiber.Router) {
  app.Post("/login", handlers.Login)

  app.Use(middleware.AuthRequired)

  studentGroup := app.Group("/students")
  studentGroup.Post("/", middleware.FacultyOnly, handlers.CreateStudent)
  studentGroup.Put("/:id", middleware.FacultyOnly, handlers.UpdateStudent)
  studentGroup.Delete("/:id", middleware.FacultyOnly, handlers.DeleteStudent)

  studentGroup.Get("/", handlers.GetStudents)
  studentGroup.Get("/:id", handlers.GetStudent)
  studentGroup.Get("/:id/transcript", handlers.GetStudentTranscript)
  studentGroup.Get("/:id/gpa", handlers.CalculateGPA)

  courseGroup := app.Group("/courses")
  courseGroup.Get("/", handlers.GetCourses)
  courseGroup.Get("/:id", handlers.GetCourse)
  courseGroup.Post("/", middleware.FacultyOnly, handlers.CreateCourse)
  courseGroup.Put("/:id", middleware.FacultyOnly, handlers.UpdateCourse)
  courseGroup.Delete("/:id", middleware.FacultyOnly, handlers.DeleteCourse)

  enrollmentGroup := app.Group("/enrollments")
  enrollmentGroup.Post("/", middleware.FacultyOnly, handlers.EnrollStudent)
  enrollmentGroup.Delete("/:id", middleware.FacultyOnly, handlers.DeleteEnrollment)
  enrollmentGroup.Get("/", handlers.GetEnrollments)
  enrollmentGroup.Get("/:id", handlers.GetEnrollment)

  gradeGroup := app.Group("/grades")
  gradeGroup.Post("/", middleware.FacultyOnly, handlers.AddGrade)
  gradeGroup.Put("/:id", middleware.FacultyOnly, handlers.UpdateGrade)
  gradeGroup.Delete("/:id", middleware.FacultyOnly, handlers.DeleteGrade)
  gradeGroup.Get("/", handlers.GetGrades)
  gradeGroup.Get("/:id", handlers.GetGrade)
}

