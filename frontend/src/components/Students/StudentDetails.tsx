import { useEffect, useState } from 'react';
import { Student, StudentTranscript, Grade,  TranscriptCourse } from '../../types/types';
import { getStudent, getStudentTranscript, calculateStudentGPA, addGrade, updateGrade, deleteGrade } from '../../api/api';

interface StudentDetailsProps {
  studentId: number;
  onBack: () => void;
  userRole: 'student' | 'faculty';
  userId: number;
}

function StudentDetails({ studentId, onBack, userRole, userId }: StudentDetailsProps) {
  const [student, setStudent] = useState<Student | null>(null);
  const [transcript, setTranscript] = useState<StudentTranscript | null>(null);
  const [gpa, setGpa] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showGradeModal, setShowGradeModal] = useState(false);
  const [selectedEnrollmentId, setSelectedEnrollmentId] = useState<number | null>(null);
  const [gradeForm, setGradeForm] = useState<Grade>({ enrollment_id: 0, grade: undefined, semester: 0 });
  const [isEditingGrade, setIsEditingGrade] = useState(false);
  const [currentGradeId, setCurrentGradeId] = useState<number | null>(null);

  useEffect(() => {
    if (userRole === 'student' && studentId !== userId) {
      setError("Access denied. You can only view your own details.");
      setLoading(false);
      return;
    }

    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        const [studentRes, transcriptRes, gpaRes] = await Promise.all([
          getStudent(studentId),
          getStudentTranscript(studentId),
          calculateStudentGPA(studentId),
        ]);

        setStudent(studentRes.data);
        setTranscript(transcriptRes.data);
        setGpa(gpaRes.data.gpa);

      } catch (err: any) {
        setError(`Failed to load student data: ${err.response?.data?.error || err.message}`);
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();

  }, [studentId, userRole, userId]);

  const handleAddGradeClick = (enrollmentId: number) => {
    setSelectedEnrollmentId(enrollmentId);
    setGradeForm({ enrollment_id: enrollmentId, grade: undefined, semester: 0 });
    setIsEditingGrade(false);
    setCurrentGradeId(null);
    setShowGradeModal(true);
  };

  const handleEditGradeClick = (transcriptCourse: TranscriptCourse) => {
    if (transcriptCourse.grade_id === undefined || transcriptCourse.grade_id === null) {
      console.error("Attempted to edit a course without a grade ID");
      return;
    }
    setSelectedEnrollmentId(transcriptCourse.enrollment_id);
    setGradeForm({
      enrollment_id: transcriptCourse.enrollment_id,
      grade: transcriptCourse.grade,
      semester: transcriptCourse.semester ?? 0,
    });
    setIsEditingGrade(true);
    setCurrentGradeId(transcriptCourse.grade_id);
    setShowGradeModal(true);
  }

  const handleDeleteGrade = async (gradeId: number | undefined) => {
    if (gradeId === undefined) return;
    if (window.confirm('Are you sure you want to delete this grade?')) {
      setLoading(true);
      setError(null);
      try {
        await deleteGrade(gradeId);
        alert('Grade deleted successfully!');

        const [transcriptRes, gpaRes] = await Promise.all([
          getStudentTranscript(studentId),
          calculateStudentGPA(studentId),
        ]);
        setTranscript(transcriptRes.data);
        setGpa(gpaRes.data.gpa);

      } catch (err: any) {
        setError(`Failed to delete grade: ${err.response?.data?.error || err.message}`);
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
  }


  const handleGradeFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setGradeForm({
      ...gradeForm,
      [name]: name === 'grade' ? parseFloat(value) || undefined : name === 'semester' ? parseInt(value, 10) || 0 : value,
    });
  };

  const handleGradeFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedEnrollmentId === null) {
      setError("No enrollment selected for grading.");
      return;
    }

    setLoading(true);
    setError(null);


    if (gradeForm.semester === 0) {
      setError("Semester is required.");
      setLoading(false);
      return;
    }

    try {
      if (isEditingGrade && currentGradeId !== null) {
        await updateGrade(currentGradeId, gradeForm);
        alert('Grade updated successfully!');
      } else {
        await addGrade(gradeForm);
        alert('Grade added successfully!');
      }
      setShowGradeModal(false);
      const [transcriptRes, gpaRes] = await Promise.all([
        getStudentTranscript(studentId),
        calculateStudentGPA(studentId),
      ]);
      setTranscript(transcriptRes.data);
      setGpa(gpaRes.data.gpa);


    } catch (err: any) {
      setError(`Failed to save grade: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };


  if (loading && !student) {
    return <div className="text-center text-gray-600">Loading student details...</div>;
  }

  if (error && !student) {
    return <div className="text-center text-red-600">{error}</div>;
  }

  if (!student) {
    return <div className="text-center text-gray-600">Student not found.</div>;
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      {userRole === 'faculty' && (
        <button onClick={onBack} className="mb-6 bg-gray-300 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-400 transition duration-200 ease-in-out">
          Back to Students List
        </button>
      )}
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Student Details: {student.name}</h2>


      {error && <div className="text-center text-red-600 mb-4">{error}</div>}


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
                  {userRole === 'faculty' && (
                    <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Actions</th>
                  )}
                </tr>
              </thead>
              <tbody>
                {transcript.courses.map((course, index) => (
                  <tr key={course.enrollment_id} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.course_code}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.course_title}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.credits}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.semester ?? 'N/A'}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.grade?.toFixed(2) ?? 'N/A'}</td>
                    {userRole === 'faculty' && (
                      <td className="py-3 px-4 border-b text-sm text-gray-700">
                        {course.grade_id !== null && course.grade_id !== undefined ? (
                          <>
                            <button
                              onClick={() => handleEditGradeClick(course)}
                              className="mr-2 bg-yellow-600 text-white px-3 py-1 rounded-md hover:bg-yellow-700 transition duration-200 ease-in-out text-xs"
                            >
                              Edit Grade
                            </button>
                            <button
                              onClick={() => handleDeleteGrade(course.grade_id)}
                              className="bg-red-600 text-white px-3 py-1 rounded-md hover:bg-red-700 transition duration-200 ease-in-out text-xs"
                            >
                              Delete Grade
                            </button>
                          </>
                        ) : (
                            <button
                              onClick={() => handleAddGradeClick(course.enrollment_id)}
                              className="bg-green-600 text-white px-3 py-1 rounded-md hover:bg-green-700 transition duration-200 ease-in-out text-xs"
                            >
                              Add Grade
                            </button>
                          )}
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



      </div>


      {showGradeModal && userRole === 'faculty' && selectedEnrollmentId !== null && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white p-6 rounded-lg shadow-xl max-w-sm w-full">
            <h3 className="text-xl font-semibold text-gray-800 mb-4">{isEditingGrade ? 'Edit Grade' : 'Add Grade'}</h3>
            {error && <p className="text-red-600 text-xs italic mb-4">{error}</p>}
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

                />
              </div>
              <div className="mb-6">
                <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="semester">
                  Semester (e.g., 20231):
                </label>
                <input
                  type="number"
                  id="semester"
                  name="semester"
                  value={gradeForm.semester || ''}
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
                  onClick={() => setShowGradeModal(false)}
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

