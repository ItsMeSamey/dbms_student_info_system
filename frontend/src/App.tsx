// src/App.tsx
import { useState } from 'react';
import StudentList from './components/Students/StudentList';
import CourseList from './components/Courses/CourseList';
import EnrollmentList from './components/Enrollments/EnrollmentList';
import StudentForm from './components/Students/StudentForm';
import CourseForm from './components/Courses/CourseForm';
import EnrollmentForm from './components/Enrollments/EnrollmentForm';
import StudentDetails from './components/Students/StudentDetails';

import './App.css';

type View = 'students' | 'courses' | 'enrollments' | 'add-student' | 'add-course' | 'add-enrollment' | 'student-details';

function App() {
  const [currentView, setCurrentView] = useState<View>('students');
  const [selectedStudentId, setSelectedStudentId] = useState<number | null>(null);

  const renderView = () => {
    switch (currentView) {
      case 'students':
        return <StudentList onViewDetails={(id) => { setSelectedStudentId(id); setCurrentView('student-details'); }} onAddStudent={() => setCurrentView('add-student')} />;
      case 'courses':
        return <CourseList onAddCourse={() => setCurrentView('add-course')} />;
      case 'enrollments':
        return <EnrollmentList onAddEnrollment={() => setCurrentView('add-enrollment')} />;
      case 'add-student':
        return <StudentForm onSuccess={() => setCurrentView('students')} />;
      case 'add-course':
        return <CourseForm onSuccess={() => setCurrentView('courses')} />;
      case 'add-enrollment':
        return <EnrollmentForm onSuccess={() => setCurrentView('enrollments')} />;
      case 'student-details':
        if (selectedStudentId === null) {
          setCurrentView('students'); // Redirect if no student selected
          return null;
        }
        return <StudentDetails studentId={selectedStudentId} onBack={() => setCurrentView('students')} />;
      default:
        return <StudentList onViewDetails={(id) => { setSelectedStudentId(id); setCurrentView('student-details'); }} onAddStudent={() => setCurrentView('add-student')} />;
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 font-san w-screen"> {/* Added font-sans */}
      <nav className="bg-gradient-to-r from-blue-600 to-blue-800 p-4 text-white shadow-lg w-screen"> {/* Added gradient and shadow */}
        <div className="container mx-auto flex flex-col md:flex-row justify-between items-center w-screen"> {/* Added flex-col for mobile */}
          <h1 className="text-xl font-bold mb-2 md:mb-0">Student Info System</h1> {/* Adjusted margin for mobile */}
          <div className="flex space-x-4"> {/* Added space-x for spacing */}
            <button onClick={() => setCurrentView('students')} className="hover:underline underline-offset-4">Students</button> {/* Added underline effect */}
            <button onClick={() => setCurrentView('courses')} className="hover:underline underline-offset-4">Courses</button>
            <button onClick={() => setCurrentView('enrollments')} className="hover:underline underline-offset-4">Enrollments</button>
          </div>
        </div>
      </nav>
      <div className="container mx-auto p-4 mt-4"> {/* Added margin-top */}
        {renderView()}
      </div>
    </div>
  );
}

export default App;
