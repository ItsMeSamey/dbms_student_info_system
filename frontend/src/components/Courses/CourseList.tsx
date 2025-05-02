// src/components/Courses/CourseList.tsx
import { useEffect, useState } from 'react';
import { Course } from '../../types/types';
import { getCourses, deleteCourse } from '../../api/api';

interface CourseListProps {
  onAddCourse: () => void;
}

function CourseList({ onAddCourse }: CourseListProps) {
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchCourses();
  }, []);

  const fetchCourses = async () => {
    try {
      const response = await getCourses();
      setCourses(response.data);
      setLoading(false);
    } catch (err) {
      setError('Failed to fetch courses');
      setLoading(false);
      console.error(err);
    }
  };

  const handleDelete = async (id: number | undefined) => {
    if (id === undefined) return;
    if (window.confirm('Are you sure you want to delete this course?')) {
      try {
        await deleteCourse(id);
        fetchCourses(); // Refresh the list
      } catch (err) {
        alert('Failed to delete course');
        console.error(err);
      }
    }
  };

  if (loading) {
    return <div className="text-center text-gray-600">Loading courses...</div>;
  }

  if (error) {
    return <div className="text-center text-red-600">{error}</div>;
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-xl">
      <h2 className="text-2xl font-bold text-gray-800 mb-4 border-b pb-2">Courses</h2>
      <button
        onClick={onAddCourse}
        className="mb-6 bg-green-600 text-white px-6 py-2 rounded-md hover:bg-green-700 transition duration-200 ease-in-out shadow-md"
      >
        Add New Course
      </button>
      {courses.length === 0 ? (
        <p className="text-gray-600">No courses found.</p>
      ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full bg-white border border-gray-200 rounded-md overflow-hidden">
              <thead className="bg-gray-200">
                <tr>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">ID</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Code</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Title</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Credits</th>
                  <th className="py-3 px-4 border-b text-left text-sm font-semibold text-gray-700">Actions</th>
                </tr>
              </thead>
              <tbody>
                {courses.map((course, index) => (
                  <tr key={course.id} className={`${index % 2 === 0 ? 'bg-gray-50' : 'bg-white'} hover:bg-gray-100 transition duration-150 ease-in-out`}>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.id}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.code}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.title}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">{course.credits}</td>
                    <td className="py-3 px-4 border-b text-sm text-gray-700">
                      <button
                        onClick={() => handleDelete(course.id)}
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

export default CourseList;

