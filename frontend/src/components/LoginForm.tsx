import { useState } from 'react';
import { LoginRequest } from '../types/types';

interface LoginFormProps {
  onLogin: (credentials: LoginRequest) => void;
  loading: boolean;
  error: string | null;
}

function LoginForm({ onLogin, loading, error }: LoginFormProps) {
  const [credentials, setCredentials] = useState<LoginRequest>({
    id: 0,
    password: '',
    role: 'student',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setCredentials({
      ...credentials,
      [name]: name === 'id' ? parseInt(value, 10) || 0: value,
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (credentials.id === 0) {
      alert("Please enter a valid ID.");
      return;
    }
    onLogin(credentials);
  };

  return (
    <div className="bg-white p-8 rounded-lg shadow-xl max-w-sm mx-auto mt-20">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">Login</h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="id">
            ID: {/* Changed label from Account to ID */}
          </label>
          <input
            type="number" // Changed input type to number
            id="id"
            name="id"
            value={credentials.id === 0 ? '' : credentials.id}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
            min="1"
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="password">
            Password:
          </label>
          <input
            type="password"
            id="password"
            name="password"
            value={credentials.password}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="role">
            Role:
          </label>
          <select
            id="role"
            name="role"
            value={credentials.role}
            onChange={handleChange}
            className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          >
            <option value="student">Student</option>
            <option value="faculty">Faculty</option>
          </select>
        </div>
        {error && <p className="text-red-600 text-xs italic mb-4">{error}</p>}
        <div className="flex items-center justify-between">
          <button
            type="submit"
            className={`bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-md focus:outline-none focus:shadow-outline transition duration-200 ease-in-out ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            disabled={loading}
          >
            {loading ? 'Logging In...' : 'Login'}
          </button>
        </div>
      </form>
    </div>
  );
}

export default LoginForm;

