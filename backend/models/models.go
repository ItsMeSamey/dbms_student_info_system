package models

import (
  "time"
)

// Database Models

type Student struct {
  ID          int       `json:"id,omitempty"`
  Name        string    `json:"name"`
  Password    string    `json:"password"`
  DateOfBirth time.Time `json:"date_of_birth"`
  Address     string    `json:"address,omitempty"`
  Contact     string    `json:"contact,omitempty"`
  Program     string    `json:"program,omitempty"`
}

type Faculty struct {
  ID          int       `json:"id,omitempty"`
  Name        string    `json:"name"`
  Password    string    `json:"password"`
  DateOfBirth time.Time `json:"date_of_birth"`
  Info        string    `json:"info,omitempty"`
}

type Course struct {
  ID      int     `json:"id,omitempty"`
  Code    string  `json:"code"`
  Title   string  `json:"title"`
  Credits float32 `json:"credits"`
}

type Enrollment struct {
  ID             int       `json:"id,omitempty"`
  StudentID      int       `json:"student_id"`
  CourseID       int       `json:"course_id"`
  EnrollmentDate time.Time `json:"enrollment_date"`
}

type Grade struct {
  ID           int      `json:"id,omitempty"`
  EnrollmentID int      `json:"enrollment_id"`
  Grade        *float64 `json:"grade"`
  Semester     int      `json:"semester"`
}

// API Models

type StudentTranscript struct {
  StudentID   int              `json:"student_id"`
  StudentName string           `json:"student_name"`
  Courses     []TranscriptCourse `json:"courses"`
}

type TranscriptCourse struct {
  EnrollmentID int      `json:"enrollment_id"`
  CourseCode   string   `json:"course_code"`
  CourseTitle  string   `json:"course_title"`
  Credits      float32  `json:"credits"` 
  GradeID      *int     `json:"grade_id,omitempty"`
  Grade        *float64 `json:"grade"`
  Semester     *int     `json:"semester"`
}

type LoginRequest struct {
  ID       int    `json:"id"`
  Password string `json:"password"`
  Role     string `json:"role"`
}

type AuthResponse struct {
  Token string `json:"token"`
  Role  string `json:"role"`
  ID    int    `json:"id"`
}

