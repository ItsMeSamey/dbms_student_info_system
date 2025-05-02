// src/components/Students/StudentDetails.tsx
import { useEffect, useState } from 'react';
import { Student, StudentTranscript, Grade, Enrollment } from '../../types/types'; // Import Enrollment
import { getStudent, getStudentTranscript, calculateStudentGPA, addGrade, updateGrade, deleteGrade, getEnrollments } from '../../api/api';

interface StudentDetailsProps {
  studentId: number;
  onBack: () => void;
}

function StudentDetails({ studentId, onBack }: StudentDetailsProps) {
  const [student, setStudent] = useState<Student | null>(null);
  const [transcript, setTranscript] = useState<StudentTranscript | null>(null);
  const [gpa, setGpa] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showAddGradeModal, setShowAddGradeModal] = useState(false);
  const [selectedEnrollmentId, setSelectedEnrollmentId] = useState<number | null>(null);
  const [gradeForm, setGradeForm] = useState<Grade>({ enrollment_id: 0, grade: undefined, semester: '' });
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]); // To select enrollment for grading
  const [isEditingGrade, setIsEditingGrade] = useState(false); // To track if we are editing or adding a grade
  const [currentGradeId, setCurrentGradeId] = useState<number | null>(null); // To store the ID of the grade being edited


  useEffect(() => {
    fetchStudentDetails();
    fetchStudentTranscript();
    fetchStudentGPA();
    fetchStudentEnrollments(); // Fetch enrollments for grading
  }, [studentId]);

  const fetchStudentDetails = async () => {
    try {
      const response = await getStudent(studentId);
      setStudent(response.data);
    } catch (err) {
      setError('Failed to fetch student details');
      console.error(err);
    }
  };

  const fetchStudentTranscript = async () => {
    try {
      const response = await getStudentTranscript(studentId);
      setTranscript(response.data);
    } catch (err) {
      setError('Failed to fetch student transcript');
      console.error(err);
    }
  };

  const fetchStudentGPA = async () => {
    try {
      const response = await calculateStudentGPA(studentId);
      setGpa(response.data.gpa);
    } catch (err) {
      setError('Failed to calculate GPA');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchStudentEnrollments = async () => {
    try {
      const response = await getEnrollments();
      // Filter enrollments for the current student
      const studentEnrollments = response.data.filter(enrollment => enrollment.student_id === studentId);
      setEnrollments(studentEnrollments);
    } catch (err) {
      setError('Failed to fetch student enrollments for grading');
      console.error(err);
    }
  }


  const handleAddGradeClick = (enrollmentId: number) => {
    setSelectedEnrollmentId(enrollmentId);
    setGradeForm({ enrollment_id: enrollmentId, grade: undefined, semester: '' }); // Reset form
    setIsEditingGrade(false); // We are adding a new grade
    setCurrentGradeId(null);
    setShowAddGradeModal(true);
  };

  const handleEditGradeClick = (grade: Grade) => {
    setSelectedEnrollmentId(grade.enrollment_id);
    setGradeForm({ enrollment_id: grade.enrollment_id, grade: grade.grade, semester: grade.semester });
    setIsEditingGrade(true); // We are editing an existing grade
    setCurrentGradeId(grade.id!); // Store the grade ID
    setShowAddGradeModal(true);
  }

  const handleDeleteGrade = async (gradeId: number | undefined) => {
    if (gradeId === undefined) return;
    if (window.confirm('Are you sure you want to delete this grade?')) {
      try {
        await deleteGrade(gradeId);
        alert('Grade deleted successfully!');
        fetchStudentTranscript(); // Refresh transcript
        fetchStudentGPA(); // Recalculate GPA
      } catch (err) {
        alert('Failed to delete grade');
        console.error(err);
      }
    }
  }


  const handleGradeFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setGradeForm({
      ...gradeForm,
      [name]: name === 'grade' ? parseFloat(value) || undefined : value,
    });
  };

  const handleGradeFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedEnrollmentId === null) return;

    setLoading(true);
    setError(null);

    try {
      if (isEditingGrade && currentGradeId !== null) {
        // Update existing grade
        await updateGrade(currentGradeId, gradeForm);
        alert('Grade updated successfully!');
      } else {
        // Add new grade
        await addGrade(gradeForm);
        alert('Grade added successfully!');
      }
      setShowAddGradeModal(false);
      fetchStudentTranscript(); // Refresh transcript
      fetchStudentGPA(); // Recalculate GPA

    } catch (err: any) {
      setError(`Failed to save grade: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };


  if (loading) {
    return <div className="text-center text-gray-600">Loading student details...</div>;
  }

  if (error) {
    return <div className="text-center text-red-600">{error}</div>;
  }

  if (!student) {
    return <div className="text-center text-gray-600">Student not found.</div>;
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      <button onClick={onBack} className="mb-6 bg-gray-300 text-gray-100 px-4 py-2 rounded-md hover:bg-gray-400 transition duration-200 ease-in-out">
        Back to Students List
      </button>
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Student Details: {student.name}</h2>

      <div className="mb-6 p-4 bg-gray-50 rounded-md"> {/* Added background and padding */}
        <h3 className="text-xl font-semibold text-gray-700 mb-3">Personal Information</h3> {/* Styled heading */}
        <p className="mb-1"><strong className="font-medium text-gray-700">ID:</strong> {student.id}</p> {/* Styled text */}
        <p className="mb-1"><strong className="font-medium text-gray-700">Name:</strong> {student.name}</p>
        <p className="mb-1"><strong className="font-medium text-gray-700">Date of Birth:</strong> {student.date_of_birth ? new Date(student.date_of_birth).toLocaleDateString() : 'N/A'}</p>
        <p className="mb-1"><strong className="font-medium text-gray-700">Address:</strong> {student.address || 'N/A'}</p>
        <p className="mb-1"><strong className="font-medium text-gray-700">Contact:</strong> {student.contact || 'N/A'}</p>
        <p><strong className="font-medium text-gray-700">Program:</strong> {student.program || 'N/A'}</p>
      </div>

      <div className="mb-6 p-4 bg-gray-50 rounded-md">
        <h3 className="text-xl font-semibold text-gray-700 mb-3">Academic Information</h3>
        <p className="mb-4"><strong className="font-medium text-gray-700">Calculated GPA:</strong> {gpa !== null ? gpa.toFixed(2) : 'N/A'}</p>

        <h4 className="text-lg font-semibold text-gray-700 mt-4 mb-2">Transcript</h4>
        {transcript && transcript.courses.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full bg-white border border-gray-200 rounded-md overflow-hidden">
              <thead className="bg-gray-200">
                <tr>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Course Code</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Course Title</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Credits</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Semester</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Grade</th>
                  {/* Add column for actions if needed */}
                </tr>
              </thead>
              <tbody>
                {transcript.courses.map((course, index) => (
                  <tr key={index} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.course_code}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.course_title}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.credits}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.semester || 'N/A'}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.grade !== -1 ? course.grade.toFixed(2) : 'N/A'}</td>
                    {/* Actions column could go here */}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
            <p className="text-gray-600">No transcript data available.</p>
          )}

        <h4 className="text-lg font-semibold text-gray-700 mt-6 mb-2">Manage Grades</h4>
        {enrollments.length > 0 ? (
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="select_enrollment_grade">
              Select Enrollment to Add/Edit Grade:
            </label>
            <select
              id="select_enrollment_grade"
              value={selectedEnrollmentId || 0}
              onChange={(e) => handleAddGradeClick(parseInt(e.target.value, 10))}
              className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value={0}>-- Select Enrollment --</option>
              {enrollments.map(enrollment => (
                <option key={enrollment.id} value={enrollment.id}>
                  Enrollment ID: {enrollment.id} (Course ID: {enrollment.course_id}) {/* Display more info if needed */}
                </option>
              ))}
            </select>
          </div>
        ) : (
            <p className="text-gray-600">No enrollments found for this student to add grades.</p>
          )}


      </div>

      {/* Add/Edit Grade Modal */}
      {showAddGradeModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4"> {/* Added padding for mobile */}
          <div className="bg-white p-6 rounded-lg shadow-xl max-w-sm w-full"> {/* Increased shadow and max-width */}
            <h3 className="text-xl font-semibold text-gray-800 mb-4">{isEditingGrade ? 'Edit Grade' : 'Add Grade'}</h3> {/* Dynamic title */}
            <form onSubmit={handleGradeFormSubmit}>
              <div className="mb-4">
                <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="grade">
                  Grade (e.g., 4.0):
                </label>
                <input
                  type="number"
                  id="grade"
                  name="grade"
                  value={gradeForm.grade ?? ''} // Use ?? '' to handle undefined
                  onChange={handleGradeFormChange}
                  step="0.01"
                  className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>
              <div className="mb-6"> {/* Increased bottom margin */}
                <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="semester">
                  Semester:
                </label>
                <input
                  type="text"
                  id="semester"
                  name="semester"
                  value={gradeForm.semester}
                  onChange={handleGradeFormChange}
                  className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>
              {error && <p className="text-red-600 text-xs italic mb-4">{error}</p>}
              <div className="flex items-center justify-between">
                <button
                  type="submit"
                  className={`bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline transition duration-200 ease-in-out ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
                  disabled={loading}
                >
                  {loading ? 'Saving...' : 'Save Grade'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowAddGradeModal(false)}
                  className="inline-block align-baseline font-bold text-sm text-gray-600 hover:text-gray-800"
                  disabled={loading}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}

export default StudentDetails;

