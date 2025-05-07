import React from 'react';

function Error() {
  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="text-center p-6 bg-red-500 text-white rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold">404</h1>
        <p className="mt-2 text-xl">Oops! The page you are looking for was not found.</p>
      </div>
    </div>
  );
}

export default Error;
