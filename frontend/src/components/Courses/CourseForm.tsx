import { useState } from 'react';
import { Course } from '../../types/types';
import { createCourse } from '../../api/api';

interface CourseFormProps {
  onSuccess: () => void;
}

function CourseForm({ onSuccess }: CourseFormProps) {
  const [course, setCourse] = useState<Course>({
    code: '',
    title: '',
    credits: 0,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCourse({
      ...course,
      [name]: name === 'credits' ? parseFloat(value) || 0 : value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);


    if (course.code === '' || course.title === '' || course.credits <= 0) {
      setError("Course code, title, and positive credits are required.");
      setLoading(false);
      return;
    }


    try {
      await createCourse(course);
      alert('Course created successfully!');
      onSuccess();
    } catch (err: any) {
      setError(`Failed to create course: ${err.response?.data?.error || err.message}`);
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white p-8 rounded-lg shadow-xl max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">Add New Course</h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="code">
            Course Code:
          </label>
          <input
            type="text"
            id="code"
            name="code"
            value={course.code}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="title">
            Title:
          </label>
          <input
            type="text"
            id="title"
            name="title"
            value={course.title}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="credits">
            Credits:
          </label>
          <input
            type="number"
            id="credits"
            name="credits"
            value={course.credits}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
            min="1"
            step="0.1"
          />
        </div>
        {error && <p className="text-red-600 text-xs italic mb-4">{error}</p>}
        <div className="flex items-center justify-between">
          <button
            type="submit"
            className={`bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline transition duration-200 ease-in-out ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            disabled={loading}
          >
            {loading ? 'Saving...' : 'Add Course'}
          </button>
          <button
            type="button"
            onClick={onSuccess}
            className="inline-block align-baseline font-bold text-sm text-gray-600 hover:text-gray-800"
            disabled={loading}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default CourseForm;

