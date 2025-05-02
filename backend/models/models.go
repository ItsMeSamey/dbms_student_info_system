package models

import (
	"time"
)

// Student represents a student in the system
type Student struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Address     string    `json:"address"`
	Contact     string    `json:"contact"`
	Program     string    `json:"program"`
}

// Course represents a course offered
type Course struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Title   string `json:"title"`
	Credits int    `json:"credits"`
}

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	ID             int       `json:"id"`
	StudentID      int       `json:"student_id"`
	CourseID       int       `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
}

// Grade represents a grade received by a student in an enrollment
type Grade struct {
	ID           int     `json:"id"`
	EnrollmentID int     `json:"enrollment_id"`
	Grade        float64 `json:"grade"`
	Semester     string  `json:"semester"`
}

// StudentTranscript represents a student's transcript with course and grade information
type StudentTranscript struct {
	StudentID   int                `json:"student_id"`
	StudentName string             `json:"student_name"`
	Courses     []TranscriptCourse `json:"courses"`
}

// TranscriptCourse represents a course and the associated grade in a transcript
type TranscriptCourse struct {
	CourseCode  string  `json:"course_code"`
	CourseTitle string  `json:"course_title"`
	Credits     int     `json:"credits"`
	Grade       float64 `json:"grade"`
	Semester    string  `json:"semester"`
}
