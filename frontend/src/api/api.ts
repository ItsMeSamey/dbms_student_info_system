import axios from 'axios';
import { Student, Course, Enrollment, Grade, StudentTranscript, LoginRequest, AuthResponse } from '../types/types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

const api = axios.create({
	baseURL: API_URL,
	headers: {
		'Content-Type': 'application/json',
	},
});

export const setAuthToken = (token: string | null) => {
	if (token) {
		api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
	} else {
		delete api.defaults.headers.common['Authorization'];
	}
};

export const login = (credentials: LoginRequest) => api.post<AuthResponse>('/login', credentials);

export const getStudents = () => api.get<Student[]>('/students');
export const getStudent = (id: number) => api.get<Student>(`/students/${id}`);
export const createStudent = (student: Student) => api.post<Student>('/students', student);
export const updateStudent = (id: number, student: Student) => api.put<void>(`/students/${id}`, student);
export const deleteStudent = (id: number) => api.delete<void>(`/students/${id}`);
export const getStudentTranscript = (id: number) => api.get<StudentTranscript>(`/students/${id}/transcript`);
export const calculateStudentGPA = (id: number) => api.get<{ student_id: number; gpa: number; message?: string }>(`/students/${id}/gpa`);

export const getCourses = () => api.get<Course[]>('/courses');
export const getCourse = (id: number) => api.get<Course>(`/courses/${id}`);
export const createCourse = (course: Course) => api.post<Course>('/courses', course);
export const updateCourse = (id: number, course: Course) => api.put<void>(`/courses/${id}`, course);
export const deleteCourse = (id: number) => api.delete<void>(`/courses/${id}`);


export const getEnrollments = (studentId?: number) => {
	const params = studentId !== undefined ? { student_id: studentId } : {};
	return api.get<Enrollment[]>('/enrollments', { params });
};
export const getEnrollment = (id: number) => api.get<Enrollment>(`/enrollments/${id}`);
export const createEnrollment = (enrollment: Enrollment) => api.post<Enrollment>('/enrollments', enrollment);
export const deleteEnrollment = (id: number) => api.delete<void>(`/enrollments/${id}`);

export const getGrades = () => api.get<Grade[]>('/grades');
export const getGrade = (id: number) => api.get<Grade>(`/grades/${id}`);
export const addGrade = (grade: Grade) => api.post<Grade>('/grades', grade);
export const updateGrade = (id: number, grade: Grade) => api.put<void>(`/grades/${id}`, grade);
export const deleteGrade = (id: number) => api.delete<void>(`/grades/${id}`);

export default api;

