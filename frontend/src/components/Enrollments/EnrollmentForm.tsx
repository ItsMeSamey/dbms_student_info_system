// src/components/Enrollments/EnrollmentForm.tsx
import { useEffect, useState } from 'react';
import { Enrollment, Student, Course } from '../../types/types';
import { createEnrollment, getStudents, getCourses } from '../../api/api';

interface EnrollmentFormProps {
  onSuccess: () => void;
  // Add onError handling if needed
}

function EnrollmentForm({ onSuccess }: EnrollmentFormProps) {
  const [enrollment, setEnrollment] = useState<Enrollment>({
    student_id: 0,
    course_id: 0,
  });
  const [students, setStudents] = useState<Student[]>([]);
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [dataLoading, setDataLoading] = useState(true);

  useEffect(() => {
    fetchStudentsAndCourses();
  }, []);

  const fetchStudentsAndCourses = async () => {
    try {
      const [studentsRes, coursesRes] = await Promise.all([
        getStudents(),
        getCourses(),
      ]);
      setStudents(studentsRes.data);
      setCourses(coursesRes.data);
      setDataLoading(false);
    } catch (err) {
      setError('Failed to fetch students and courses');
      setDataLoading(false);
      console.error(err);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const { name, value } = e.target;
    setEnrollment({
      ...enrollment,
      [name]: parseInt(value, 10),
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    if (enrollment.student_id === 0 || enrollment.course_id === 0) {
      setError("Please select a student and a course.");
      setLoading(false);
      return;
    }

    try {
      await createEnrollment(enrollment);
      alert('Enrollment created successfully!');
      onSuccess(); // Navigate back to the list or another page
    } catch (err: any) {
      setError(`Failed to create enrollment: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (dataLoading) {
    return <div className="text-center text-gray-600">Loading data for enrollment form...</div>;
  }

  if (error && dataLoading) {
    return <div className="text-center text-red-600">{error}</div>;
  }


  return (
    <div className="bg-white p-8 rounded-lg shadow-xl max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">Add New Enrollment</h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="student_id">
            Select Student:
          </label>
          <select
            id="student_id"
            name="student_id"
            value={enrollment.student_id}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          >
            <option value={0}>-- Select Student --</option>
            {students.map((student) => (
              <option key={student.id} value={student.id}>
                {student.name} (ID: {student.id})
              </option>
            ))}
          </select>
        </div>
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="course_id">
            Select Course:
          </label>
          <select
            id="course_id"
            name="course_id"
            value={enrollment.course_id}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          >
            <option value={0}>-- Select Course --</option>
            {courses.map((course) => (
              <option key={course.id} value={course.id}>
                {course.title} (Code: {course.code})
              </option>
            ))}
          </select>
        </div>
        {error && <p className="text-red-600 text-xs italic mb-4">{error}</p>}
        <div className="flex items-center justify-between">
          <button
            type="submit"
            className={`bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline transition duration-200 ease-in-out ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            disabled={loading}
          >
            {loading ? 'Saving...' : 'Create Enrollment'}
          </button>
          <button
            type="button"
            onClick={onSuccess}
            className="inline-block align-baseline font-bold text-sm text-gray-600 hover:text-gray-800"
            disabled={loading}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default EnrollmentForm;

