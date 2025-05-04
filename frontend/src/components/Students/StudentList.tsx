import { useEffect, useState } from 'react';
import { Student } from '../../types/types';
import { getStudents, deleteStudent } from '../../api/api';

interface StudentListProps {
  onViewDetails: (id: number) => void;
  onAddStudent: () => void;
  userRole: 'student' | 'faculty';
  userId: number;
}

function StudentList({ onViewDetails, onAddStudent, userRole, userId }: StudentListProps) {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStudents();
  }, [userRole, userId]); // Refetch if user or role changes

  const fetchStudents = async () => {
    setLoading(true);
    setError(null);
    try {
      // Backend handles filtering for student role, so frontend just calls getStudents
      const response = await getStudents();
      setStudents(response.data);
    } catch (err) {
      setError('Failed to fetch students');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number | undefined) => {
    if (id === undefined) return;
    if (window.confirm('Are you sure you want to delete this student?')) {
      try {
        await deleteStudent(id);
        fetchStudents(); // Refresh list after deletion
      } catch (err) {
        alert('Failed to delete student');
        console.error(err);
      }
    }
  };

  if (loading) {
    return <div className="text-center text-gray-600">Loading students...</div>;
  }

  if (error) {
    return <div className="text-center text-red-600">{error}</div>;
  }

  // Students only see their own details, handled by backend.
  // This component is primarily for the faculty view.
  if (userRole === 'student') {
    // Redirect student to their details page if they land here somehow
    useEffect(() => {
      onViewDetails(userId);
    }, [userId, onViewDetails]);
    return <div className="text-center text-gray-600">Redirecting to your details...</div>;
  }


  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Students</h2>
      {userRole === 'faculty' && ( // Only faculty can add students
        <button
          onClick={onAddStudent}
          className="mb-6 bg-green-600 text-white px-6 py-2 rounded-md hover:bg-green-700 transition duration-200 ease-in-out shadow-md"
        >
          Add New Student
        </button>
      )}
      {students.length === 0 ? (
        <p className="text-gray-600">No students found.</p>
      ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full bg-white border border-gray-200 rounded-md overflow-hidden">
              <thead className="bg-gray-200">
                <tr>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">ID</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Name</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Program</th>
                  {userRole === 'faculty' && ( // Only faculty sees actions column
                    <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Actions</th>
                  )}
                </tr>
              </thead>
              <tbody>
                {students.map((student, index) => (
                  <tr key={student.id} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{student.id}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{student.name}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{student.program || 'N/A'}</td>
                    {userRole === 'faculty' && ( // Only faculty sees action buttons
                      <td className="py-3 px-4 border-b text-sm text-gray-700">
                        <button
                          onClick={() => onViewDetails(student.id!)}
                          className="mr-2 bg-blue-600 text-white px-3 py-1 rounded-md hover:bg-blue-700 transition duration-200 ease-in-out text-xs"
                        >
                          View Details
                        </button>
                        <button
                          onClick={() => handleDelete(student.id)}
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

export default StudentList;

