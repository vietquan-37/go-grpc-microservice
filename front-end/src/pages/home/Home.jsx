import React, { useEffect, useState } from 'react';
import axios from 'axios';

function Home() {
  const [user, setUser] = useState(null);
  const user_id = localStorage.getItem('id'); // Get user ID from localStorage

  useEffect(() => {
    const fetchUser = async () => {
      try {
        if (!user_id) return;
        const res = await axios.get(`${process.env.VITE_BACK_END}/v1/user/${user_id}`);
        setUser(res.data);
      } catch (error) {
        console.error('Error fetching user:', error);
      }
    };

    fetchUser();
  }, [user_id]);

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-white shadow-lg rounded-2xl p-6">
        <h2 className="text-2xl font-bold mb-4 text-center text-blue-600">User Profile</h2>
        {user ? (
          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="font-medium text-gray-700">User ID:</span>
              <span className="text-gray-900">{user.user_id}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium text-gray-700">Name:</span>
              <span className="text-gray-900">{user.user_name}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium text-gray-700">Phone:</span>
              <span className="text-gray-900">{user.phone_number}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium text-gray-700">Role:</span>
              <span className="text-gray-900 capitalize">{user.role}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium text-gray-700">Created At:</span>
              <span className="text-gray-900">
                {new Date(user.create_at).toLocaleString()}
              </span>
            </div>
          </div>
        ) : (
          <p className="text-center text-gray-500">
            {user_id ? 'Loading user data...' : 'User ID not found in localStorage'}
          </p>
        )}
      </div>
    </div>
  );
}

export default Home;
