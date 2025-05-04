import { useEffect, useState } from 'react';
import { Enrollment, Student, Course } from '../../types/types';
import { getEnrollments, deleteEnrollment, getStudents, getCourses } from '../../api/api';

interface EnrollmentListProps {
  onAddEnrollment: () => void;
  userRole: 'student' | 'faculty';
  userId: number;
}

function EnrollmentList({ onAddEnrollment, userRole, userId }: EnrollmentListProps) {
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [students, setStudents] = useState<Student[]>([]); // Needed to map student_id to name
  const [courses, setCourses] = useState<Course[]>([]); // Needed to map course_id to title
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchData();
  }, [userRole, userId]); // Refetch if user or role changes

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      // Fetch enrollments. Backend handles filtering for student role.
      // Faculty view needs student and course data to display names/titles.
      // TODO: Backend getStudents/getCourses should be filtered for faculty if needed
      const [enrollmentsRes, studentsRes, coursesRes] = await Promise.all([
        getEnrollments(), // Fetch enrollments (backend filters for student, gets all for faculty)
        getStudents(), // Fetch all students (consider filtering this on backend for faculty)
        getCourses(), // Fetch all courses (consider filtering this on backend for faculty)
      ]);
      setEnrollments(enrollmentsRes.data);
      setStudents(studentsRes.data);
      setCourses(coursesRes.data);
    } catch (err: any) {
      setError(`Failed to fetch data: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const getStudentName = (studentId: number) => {
    const student = students.find((s) => s.id === studentId);
    return student ? student.name : 'Unknown Student';
  };

  const getCourseTitle = (courseId: number) => {
    const course = courses.find((c) => c.id === courseId);
    return course ? course.title : 'Unknown Course';
  };


  const handleDelete = async (id: number | undefined) => {
    if (id === undefined) return;
    if (window.confirm('Are you sure you want to delete this enrollment?')) {
      setLoading(true); // Show loading indicator for delete
      setError(null);
      try {
        await deleteEnrollment(id);
        alert('Enrollment deleted successfully!'); // Consider a better UI notification
        fetchData(); // Refresh list after deletion
      } catch (err: any) {
        setError(`Failed to delete enrollment: ${err.response?.data?.error || err.message}`);
        alert('Failed to delete enrollment'); // Use a better UI notification
        console.error(err);
      } finally {
        setLoading(false); // Hide loading indicator
      }
    }
  };

  if (loading) {
    return <div className="text-center text-gray-600">Loading enrollments...</div>;
  }

  if (error) {
    return <div className="text-center text-red-600">{error}</div>;
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Enrollments</h2>
      {userRole === 'faculty' && ( // Only faculty can add enrollments
        <button
          onClick={onAddEnrollment}
          className="mb-6 bg-green-600 text-white px-6 py-2 rounded-md hover:bg-green-700 transition duration-200 ease-in-out shadow-md"
        >
          Add New Enrollment
        </button>
      )}
      {enrollments.length === 0 ? (
        <p className="text-gray-600">No enrollments found.</p>
      ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full bg-white border border-gray-200 rounded-md overflow-hidden">
              <thead className="bg-gray-200">
                <tr>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">ID</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Student</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Course</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Enrollment Date</th>
                  {userRole === 'faculty' && ( // Only faculty sees actions column
                    <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Actions</th>
                  )}
                </tr>
              </thead>
              <tbody>
                {enrollments.map((enrollment, index) => (
                  <tr key={enrollment.id} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{enrollment.id}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{getStudentName(enrollment.student_id)}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{getCourseTitle(enrollment.course_id)}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{enrollment.enrollment_date ? new Date(enrollment.enrollment_date).toLocaleDateString() : 'N/A'}</td>
                    {userRole === 'faculty' && ( // Only faculty sees action buttons
                      <td className="py-3 px-4 border-b text-sm text-gray-700">
                        <button
                          onClick={() => handleDelete(enrollment.id)}
                          className="bg-red-600 text-white px-3 py-1 rounded-md hover:bg-red-700 transition duration-200 ease-in-out text-xs"
                        >
                          Delete
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
    </div>
  );
}

export default EnrollmentList;

