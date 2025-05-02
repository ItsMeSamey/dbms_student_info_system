package routes

import (
	"backend/handlers"

	"github.com/gofiber/fiber/v3"
)

// SetupRoutes defines the API routes for the application
func SetupRoutes(app fiber.Router) {
	// Student routes
	studentGroup := app.Group("/students")
	studentGroup.Post("/", handlers.CreateStudent)
	studentGroup.Get("/", handlers.GetStudents)
	studentGroup.Get("/:id", handlers.GetStudent)
	studentGroup.Put("/:id", handlers.UpdateStudent)
	studentGroup.Delete("/:id", handlers.DeleteStudent)
	studentGroup.Get("/:id/transcript", handlers.GetStudentTranscript)
	studentGroup.Get("/:id/gpa", handlers.CalculateGPA)

	// Course routes
	courseGroup := app.Group("/courses")
	courseGroup.Post("/", handlers.CreateCourse)
	courseGroup.Get("/", handlers.GetCourses)
	courseGroup.Get("/:id", handlers.GetCourse)
	courseGroup.Put("/:id", handlers.UpdateCourse)
	courseGroup.Delete("/:id", handlers.DeleteCourse)

	// Enrollment routes
	enrollmentGroup := app.Group("/enrollments")
	enrollmentGroup.Post("/", handlers.EnrollStudent)
	enrollmentGroup.Get("/", handlers.GetEnrollments)
	enrollmentGroup.Get("/:id", handlers.GetEnrollment)
	enrollmentGroup.Delete("/:id", handlers.DeleteEnrollment)

	// Grade routes
	gradeGroup := app.Group("/grades")
	gradeGroup.Post("/", handlers.AddGrade)
	gradeGroup.Get("/", handlers.GetGrades)
	gradeGroup.Get("/:id", handlers.GetGrade)
	gradeGroup.Put("/:id", handlers.UpdateGrade)
	gradeGroup.Delete("/:id", handlers.DeleteGrade)
}

