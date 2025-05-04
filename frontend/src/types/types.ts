export interface Student {
  id?: number;
  name: string;
  password?: string;
  date_of_birth: string;
  address?: string;
  contact?: string;
  program: string;
}

export interface Faculty {
  id?: number;
  name: string;
  password?: string;
  date_of_birth: string;
  info?: string;
}

export interface Course {
  id?: number;
  code: string;
  title: string;
  credits: number;
}

export interface Enrollment {
  id?: number;
  student_id: number;
  course_id: number;
  enrollment_date?: string;
}

export interface Grade {
  id?: number;
  enrollment_id: number;
  grade?: number;
  semester: string;
}

export interface TranscriptCourse {
  course_code: string;
  course_title: string;
  credits: number;
  grade?: number;
  semester?: string;
}

export interface StudentTranscript {
  student_id: number;
  student_name: string;
  courses: TranscriptCourse[];
}

export interface LoginRequest {
  id: number;
  password: string;
  role: 'student' | 'faculty';
}

export interface AuthResponse {
  token: string;
  role: 'student' | 'faculty';
  id: number;
}

