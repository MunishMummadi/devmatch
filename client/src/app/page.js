// Inside HomePage.js
'use client';

import Header from './components/Header';
import ProfileCard from './components/ProfileCard';

export default function HomePage() {
  return (
    <>
      <Header />
      {/* The main container centers the wrapper div */}
      <main className="bg-gray-100 py-12 px-4 flex justify-center items-start min-h-screen">
        {/* This wrapper div constrains the width of its child */}
        <div className="w-full max-w-sm"> {/* Adjust max-w-sm as needed */}
          <ProfileCard />
        </div>
      </main>
    </>
  );
}