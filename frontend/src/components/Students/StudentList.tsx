import { useEffect, useState } from 'react';
import { Student } from '../../types/types';
import { getStudents, deleteStudent } from '../../api/api';

interface StudentListProps {
  onViewDetails: (id: number) => void;
  onAddStudent: () => void;
}

function StudentList({ onViewDetails, onAddStudent }: StudentListProps) {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStudents();
  }, []);

  const fetchStudents = async () => {
    try {
      const response = await getStudents();
      setStudents(response.data);
      setLoading(false);
    } catch (err) {
      setError('Failed to fetch students');
      setLoading(false);
      console.error(err);
    }
  };

  const handleDelete = async (id: number | undefined) => {
    if (id === undefined) return;
    if (window.confirm('Are you sure you want to delete this student?')) {
      try {
        await deleteStudent(id);
        fetchStudents();
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

  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Students</h2>
      <button
        onClick={onAddStudent}
        className="mb-6 bg-green-600 text-white px-6 py-2 rounded-md hover:bg-green-700 transition duration-200 ease-in-out shadow-md"
      >
        Add New Student
      </button>
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
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Actions</th>
                </tr>
              </thead>
              <tbody>
                {students.map((student, index) => (
                  <tr key={student.id} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{student.id}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{student.name}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{student.program || 'N/A'}</td>
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

