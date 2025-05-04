import { useState, useEffect } from 'react';
import StudentList from './components/Students/StudentList';
import CourseList from './components/Courses/CourseList';
import EnrollmentList from './components/Enrollments/EnrollmentList';
import StudentForm from './components/Students/StudentForm';
import CourseForm from './components/Courses/CourseForm';
import EnrollmentForm from './components/Enrollments/EnrollmentForm';
import StudentDetails from './components/Students/StudentDetails';
import LoginForm from './components/LoginForm';
import { login, setAuthToken } from './api/api';
import { LoginRequest, AuthResponse } from './types/types';


import './App.css';

type View = 'login' | 'students' | 'courses' | 'enrollments' | 'add-student' | 'add-course' | 'add-enrollment' | 'student-details';

function App() {
  const [currentView, setCurrentView] = useState<View>('login');
  const [selectedStudentId, setSelectedStudentId] = useState<number | null>(null);
  const [user, setUser] = useState<{ id: number | null; role: 'student' | 'faculty' | null }>({ id: null, role: null });
  const [loginLoading, setLoginLoading] = useState(false);
  const [loginError, setLoginError] = useState<string | null>(null);

  // Check for token in localStorage on initial load
  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    const storedUserId = localStorage.getItem('userId');
    const storedUserRole = localStorage.getItem('userRole') as 'student' | 'faculty' | null;

    if (storedToken && storedUserId && storedUserRole) {
      setAuthToken(storedToken);
      setUser({ id: parseInt(storedUserId, 10), role: storedUserRole });
      // Redirect to appropriate view after auto-login
      if (storedUserRole === 'faculty') {
        setCurrentView('students'); // Faculty default view
      } else {
        setCurrentView('student-details'); // Student default view (their own details)
        setSelectedStudentId(parseInt(storedUserId, 10));
      }
    }
  }, []);


  const handleLogin = async (credentials: LoginRequest) => {
    setLoginLoading(true);
    setLoginError(null);
    try {
      const response = await login(credentials);
      const { token, id, role } = response.data;

      localStorage.setItem('token', token);
      localStorage.setItem('userId', id.toString());
      localStorage.setItem('userRole', role);

      setAuthToken(token);
      setUser({ id, role });

      if (role === 'faculty') {
        setCurrentView('students');
      } else { // student
        setSelectedStudentId(id);
        setCurrentView('student-details');
      }

    } catch (err: any) {
      setLoginError(`Login failed: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoginLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    localStorage.removeItem('userRole');
    setAuthToken(null);
    setUser({ id: null, role: null });
    setSelectedStudentId(null);
    setCurrentView('login');
  };


  const renderView = () => {
    if (!user.id || !user.role) {
      return <LoginForm onLogin={handleLogin} loading={loginLoading} error={loginError} />;
    }

    // Pass user and role to components that need it
    switch (currentView) {
      case 'students':
        return <StudentList userRole={user.role} userId={user.id} onViewDetails={(id) => { setSelectedStudentId(id); setCurrentView('student-details'); }} onAddStudent={() => setCurrentView('add-student')} />;
      case 'courses':
        return <CourseList userRole={user.role} userId={user.id} onAddCourse={() => setCurrentView('add-course')} />;
      case 'enrollments':
        return <EnrollmentList userRole={user.role} userId={user.id} onAddEnrollment={() => setCurrentView('add-enrollment')} />;
      case 'add-student':
        return <StudentForm onSuccess={() => setCurrentView('students')} />;
      case 'add-course':
        return <CourseForm onSuccess={() => setCurrentView('courses')} />;
      case 'add-enrollment':
        return <EnrollmentForm onSuccess={() => setCurrentView('enrollments')} />;
      case 'student-details':
        if (selectedStudentId === null) {
          setCurrentView('students');
          return null;
        }
        // Ensure student can only view their own details
        if (user.role === 'student' && selectedStudentId !== user.id) {
           // This case should ideally be prevented by navigation logic,
           // but handle defensively. Redirect to their own details or a forbidden page.
           setSelectedStudentId(user.id); // Redirect to their own details
           alert("You can only view your own details.");
           return <StudentDetails studentId={user.id} userRole={user.role} userId={user.id} onBack={() => setCurrentView('student-details')} />; // Pass user and role
        }
        return <StudentDetails studentId={selectedStudentId} userRole={user.role} userId={user.id} onBack={() => user.role === 'faculty' ? setCurrentView('students') : setCurrentView('student-details')} />; // Pass user and role
      default:
        return <StudentList userRole={user.role} userId={user.id} onViewDetails={(id) => { setSelectedStudentId(id); setCurrentView('student-details'); }} onAddStudent={() => setCurrentView('add-student')} />;
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 font-sans w-screen">
      <nav className="bg-gradient-to-r from-blue-600 to-blue-800 p-4 text-white shadow-lg w-screen">
        <div className="container mx-auto flex flex-col md:flex-row justify-between items-center w-screen">
          <h1 className="text-xl font-bold mb-2 md:mb-0">Student Info System</h1>
          {user.id && user.role && ( // Show navigation only when logged in
            <div className="flex space-x-4">
              {user.role === 'faculty' && ( // Faculty navigation
                <>
                  <button onClick={() => setCurrentView('students')} className="hover:underline underline-offset-4">Students</button>
                  <button onClick={() => setCurrentView('courses')} className="hover:underline underline-offset-4">Courses</button>
                  <button onClick={() => setCurrentView('enrollments')} className="hover:underline underline-offset-4">Enrollments</button>
                </>
              )}
              {user.role === 'student' && ( // Student navigation
                 <button onClick={() => { setSelectedStudentId(user.id); setCurrentView('student-details'); }} className="hover:underline underline-offset-4">My Details</button>
              )}
              <button onClick={handleLogout} className="hover:underline underline-offset-4">Logout</button>
            </div>
          )}
        </div>
      </nav>
      <div className="container mx-auto p-4 mt-4">
        {renderView()}
      </div>
    </div>
  );
}

export default App;

