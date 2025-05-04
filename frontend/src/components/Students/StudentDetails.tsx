import { useEffect, useState } from 'react';
import { Student, StudentTranscript, Grade, Enrollment } from '../../types/types';
import { getStudent, getStudentTranscript, calculateStudentGPA, addGrade, updateGrade, deleteGrade, getEnrollments } from '../../api/api';

interface StudentDetailsProps {
  studentId: number;
  onBack: () => void;
  userRole: 'student' | 'faculty';
  userId: number; // Current logged-in user's ID
}

function StudentDetails({ studentId, onBack, userRole, userId }: StudentDetailsProps) {
  const [student, setStudent] = useState<Student | null>(null);
  const [transcript, setTranscript] = useState<StudentTranscript | null>(null);
  const [gpa, setGpa] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showAddGradeModal, setShowAddGradeModal] = useState(false);
  const [selectedEnrollmentId, setSelectedEnrollmentId] = useState<number | null>(null);
  const [gradeForm, setGradeForm] = useState<Grade>({ enrollment_id: 0, grade: undefined, semester: '' });
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [isEditingGrade, setIsEditingGrade] = useState(false);
  const [currentGradeId, setCurrentGradeId] = useState<number | null>(null);


  useEffect(() => {
    // Ensure student can only fetch their own details
    if (userRole === 'student' && studentId !== userId) {
      setError("Access denied. You can only view your own details.");
      setLoading(false);
      return;
    }
    fetchStudentDetails();
    fetchStudentTranscript();
    fetchStudentGPA();
    if (userRole === 'faculty') { // Only faculty needs to manage grades
      fetchStudentEnrollments();
    }
  }, [studentId, userRole, userId]); // Depend on user/role to refetch

  const fetchStudentDetails = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await getStudent(studentId);
      setStudent(response.data);
    } catch (err) {
      setError('Failed to fetch student details');
      console.error(err);
    } finally {
      // Don't set loading to false here, wait for all fetches
    }
  };

  const fetchStudentTranscript = async () => {
    setError(null); // Clear previous errors
    try {
      const response = await getStudentTranscript(studentId);
      setTranscript(response.data);
    } catch (err) {
      setError('Failed to fetch student transcript');
      console.error(err);
    } finally {
      // Don't set loading to false here, wait for all fetches
    }
  };

  const fetchStudentGPA = async () => {
    setError(null); // Clear previous errors
    try {
      const response = await calculateStudentGPA(studentId);
      setGpa(response.data.gpa);
    } catch (err) {
      setError('Failed to calculate GPA');
      console.error(err);
    } finally {
      // Set loading to false after the last fetch completes
      setLoading(false);
    }
  };

  const fetchStudentEnrollments = async () => {
    setError(null); // Clear previous errors
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
    setGradeForm({ enrollment_id: enrollmentId, grade: undefined, semester: '' });
    setIsEditingGrade(false);
    setCurrentGradeId(null);
    setShowAddGradeModal(true);
  };

  const handleEditGradeClick = (grade: Grade) => {
    setSelectedEnrollmentId(grade.enrollment_id);
    setGradeForm({ enrollment_id: grade.enrollment_id, grade: grade.grade, semester: grade.semester });
    setIsEditingGrade(true);
    setCurrentGradeId(grade.id!);
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

    setLoading(true); // Set loading for the save operation
    setError(null);

    try {
      if (isEditingGrade && currentGradeId !== null) {
        await updateGrade(currentGradeId, gradeForm);
        alert('Grade updated successfully!');
      } else {
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
      setLoading(false); // Unset loading after save operation
    }
  };


  if (loading) {
    return <div className="text-center text-gray-600">Loading student details...</div>;
  }

  if (error && !student) { // Only show fatal error if student details couldn't be loaded
    return <div className="text-center text-red-600">{error}</div>;
  }

  if (!student) {
    return <div className="text-center text-gray-600">Student not found.</div>;
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      {userRole === 'faculty' && ( // Only faculty sees the back button to the list
        <button onClick={onBack} className="mb-6 bg-gray-300 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-400 transition duration-200 ease-in-out">
          Back to Students List
        </button>
      )}
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Student Details: {student.name}</h2>

      <div className="mb-6 p-4 bg-gray-50 rounded-md">
        <h3 className="text-xl font-semibold text-gray-700 mb-3">Personal Information</h3>
        <p className="mb-1"><strong className="font-medium text-gray-700">ID:</strong> {student.id}</p>
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
                  {userRole === 'faculty' && ( // Only faculty sees actions column
                    <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Actions</th>
                  )}
                </tr>
              </thead>
              <tbody>
                {transcript.courses.map((course, index) => (
                  <tr key={index} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.course_code}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.course_title}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.credits}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.semester || 'N/A'}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.grade?.toFixed(2) ?? 'N/A'}</td>
                    {userRole === 'faculty' && ( // Only faculty sees action buttons
                      <td className="py-3 px-4 border-b text-sm text-gray-700">
                        {/* Find the grade ID for this course and student's enrollment */}
                        {/* This requires finding the enrollment first, which is complex here.
                             Assuming for simplicity that transcriptCourse might contain grade ID or we fetch grades separately if needed for edit/delete */}
                        {/* For a proper implementation, the transcript API should ideally return the grade ID */}
                        {/* For now, we'll make a simplified assumption or require faculty to go to the grades list */}
                        {/* A better approach is to fetch grades separately for the student and match them */}
                        {/* Let's add a placeholder for now */}
                        {/* <button
                            onClick={() => handleEditGradeClick(grade)} // Need the actual grade object
                            className="mr-2 bg-yellow-600 text-white px-3 py-1 rounded-md hover:bg-yellow-700 transition duration-200 ease-in-out text-xs"
                          >
                            Edit Grade
                          </button>
                          <button
                            onClick={() => handleDeleteGrade(grade.id)} // Need the actual grade ID
                            className="bg-red-600 text-white px-3 py-1 rounded-md hover:bg-red-700 transition duration-200 ease-in-out text-xs"
                          >
                            Delete Grade
                          </button> */}
                        <span className="text-gray-500">Manage via Grades list</span>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
            <p className="text-gray-600">No transcript data available.</p>
          )}

        {userRole === 'faculty' && ( // Only faculty can manage grades
          <>
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
                      Enrollment ID: {enrollment.id} (Course ID: {enrollment.course_id})
                    </option>
                  ))}
                </select>
              </div>
            ) : (
                <p className="text-gray-600">No enrollments found for this student to add grades.</p>
              )}
          </>
        )}


      </div>


      {showAddGradeModal && userRole === 'faculty' && ( // Only faculty sees the modal
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
          <div className="bg-white p-6 rounded-lg shadow-xl max-w-sm w-full">
            <h3 className="text-xl font-semibold text-gray-800 mb-4">{isEditingGrade ? 'Edit Grade' : 'Add Grade'}</h3>
            {error && <p className="text-red-600 text-xs italic mb-4">{error}</p>} {/* Show error in modal */}
            <form onSubmit={handleGradeFormSubmit}>
              <div className="mb-4">
                <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="grade">
                  Grade (e.g., 4.0):
                </label>
                <input
                  type="number"
                  id="grade"
                  name="grade"
                  value={gradeForm.grade ?? ''}
                  onChange={handleGradeFormChange}
                  step="0.01"
                  className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>
              <div className="mb-6">
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

