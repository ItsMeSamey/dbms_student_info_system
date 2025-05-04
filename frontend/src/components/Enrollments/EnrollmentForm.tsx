import { useEffect, useState } from 'react';
import { Enrollment, Student, Course } from '../../types/types';
import { createEnrollment, getStudents, getCourses } from '../../api/api';

interface EnrollmentFormProps {
  onSuccess: () => void;
}

function EnrollmentForm({ onSuccess }: EnrollmentFormProps) {
  const [enrollment, setEnrollment] = useState<Enrollment>({
    student_id: 0,
    course_id: 0,
  });
  const [students, setStudents] = useState<Student[]>([]); // Needed to populate dropdown
  const [courses, setCourses] = useState<Course[]>([]); // Needed to populate dropdown
  const [loading, setLoading] = useState(false); // For form submission
  const [error, setError] = useState<string | null>(null);
  const [dataLoading, setDataLoading] = useState(true); // For initial data fetch

  useEffect(() => {
    fetchStudentsAndCourses();
  }, []);

  const fetchStudentsAndCourses = async () => {
    setDataLoading(true);
    setError(null); // Clear errors before fetching
    try {
      // Fetch students and courses to populate dropdowns
      // TODO: Backend getStudents/getCourses should be filtered for faculty if needed
      const [studentsRes, coursesRes] = await Promise.all([
        getStudents(), // Consider filtering this on backend for faculty
        getCourses(), // Consider filtering this on backend for faculty
      ]);
      setStudents(studentsRes.data);
      setCourses(coursesRes.data);
      setDataLoading(false);
    } catch (err: any) {
      setError(`Failed to load data for form: ${err.response?.data?.error || err.message}`);
      setDataLoading(false);
      console.error(err);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const { name, value } = e.target;
    setEnrollment({
      ...enrollment,
      [name]: parseInt(value, 10), // Parse selected value as integer ID
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true); // Show loading indicator for submission
    setError(null); // Clear errors before submission

    // Client-side validation
    if (enrollment.student_id === 0 || enrollment.course_id === 0) {
      setError("Please select a student and a course.");
      setLoading(false);
      return;
    }

    try {
      await createEnrollment(enrollment);
      alert('Enrollment created successfully!'); // Consider a better UI notification
      onSuccess(); // Call success callback to navigate back
    } catch (err: any) {
      setError(`Failed to create enrollment: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoading(false); // Hide loading indicator
    }
  };

  if (dataLoading) {
    return <div className="text-center text-gray-600">Loading data for enrollment form...</div>;
  }

  // Show error if initial data fetch failed
  if (error && !dataLoading) {
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
        {/* Show submission error if any */}
        {error && !dataLoading && <p className="text-red-600 text-xs italic mb-4">{error}</p>}
        <div className="flex items-center justify-between">
          <button
            type="submit"
            className={`bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline transition duration-200 ease-in-out ${loading || dataLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
            disabled={loading || dataLoading} // Disable while submitting or initial data loading
          >
            {loading ? 'Creating...' : 'Create Enrollment'}
          </button>
          <button
            type="button"
            onClick={onSuccess}
            className="inline-block align-baseline font-bold text-sm text-gray-600 hover:text-gray-800"
            disabled={loading || dataLoading} // Disable while submitting or initial data loading
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default EnrollmentForm;

