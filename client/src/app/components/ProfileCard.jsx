'use client';

import { useState, useEffect } from 'react';
import { FaGithub, FaLinkedin } from 'react-icons/fa';

// This component now *only* renders the card.
export default function WaveProfileCard() {
    const [profileData, setProfileData] = useState({
        fullName: 'CHRISTIAN SØGAARD MOEN',
        handle: '@christianmoen',
        summary: 'Brand and communication strategy, graphic design, and illustration.',
        interests: 'Photography, Portraits, Art Direction',
        // Replace with your actual image URL or ensure it's valid
        image: 'https://images.unsplash.com/photo-1499714608240-22fc6ad53fb2',
        socials: {
            github: 'https://github.com/yourusername',
            linkedin: 'https://linkedin.com/in/yourusername',
        },
    });

    // useEffect remains the same
    useEffect(() => {
        // Fetch logic can go here if needed
    }, []);

    return (
        // Card Container: Correct classes applied
        <div className="relative w-full max-w-sm bg-white rounded-2xl overflow-hidden shadow-xl">

            {/* Wave / Gradient Top */}
            <div className="relative h-40 w-full overflow-hidden">
                {/* SVG definition remains the same */}
                <svg
                    className="absolute top-0 left-0 w-full h-full"
                    viewBox="0 0 500 150"
                    preserveAspectRatio="none"
                >
                    <defs>
                        <linearGradient id="waveGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                            <stop offset="0%" stopColor="#f953c6" />
                            <stop offset="100%" stopColor="#b91d73" />
                        </linearGradient>
                    </defs>
                    <path
                        d="M0.00,49.98 C149.99,150.00 303.56,-49.98 500.00,49.98 L500.00,0.00 L0.00,0.00 Z"
                        fill="url(#waveGradient)"
                    />
                </svg>
            </div>

            {/* Profile Image (overlapping the wave) */}
            {/* Added z-10 to ensure it's above the content flow */}
            <div className="absolute top-20 left-1/2 transform -translate-x-1/2 z-10">
                <div className="w-24 h-24 rounded-full overflow-hidden border-4 border-white shadow-md">
                    {profileData.image ? (
                        <img
                            src={profileData.image}
                            alt="Profile"
                            className="w-full h-full object-cover"
                            // Add error handling for the image if needed
                            onError={(e) => { e.target.onerror = null; e.target.src = "https://placehold.co/96x96/E2E8F0/9CA3AF?text=N/A"; }} // Example placeholder
                        />
                    ) : (
                        <div className="flex items-center justify-center w-full h-full bg-gray-200 text-gray-400 text-2xl">
                            +
                        </div>
                    )}
                </div>
            </div>

            {/* Card Content - MOVED INSIDE the main card container div */}
            <div className="pt-16 pb-8 px-6 flex flex-col items-center text-center">
                {/* Adjust pt-16 if needed based on image size/position */}

                {/* Full Name & Handle */}
                <h1 className="text-lg font-bold text-gray-800 mt-2"> {/* Added small top margin */}
                    {profileData.fullName}
                </h1>
                <span className="text-sm text-gray-500 mb-3">
                    {profileData.handle}
                </span>

                {/* Social Links */}
                <div className="flex gap-4 mt-3">
                    <a
                        href={profileData.socials.linkedin}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-gray-500 hover:text-gray-700 transition-colors duration-200" // Added transition
                        aria-label="LinkedIn Profile"
                    >
                        <FaLinkedin size={20} />
                    </a>
                    <a
                        href={profileData.socials.github}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-gray-500 hover:text-gray-700 transition-colors duration-200" // Added transition
                        aria-label="GitHub Profile"
                    >
                        <FaGithub size={20} />
                    </a>
                </div>

                {/* Summary Section */}
                <div className="mt-4 w-full bg-gray-100 rounded-lg p-4 text-sm text-gray-800 leading-relaxed shadow">
                    <h3 className="font-medium mb-1">Summary</h3>
                    <p>{profileData.summary}</p>
                </div>

                {/* Interests Section */}
                <div className="mt-4 w-full bg-gray-100 rounded-lg p-4 text-sm text-gray-800 leading-relaxed shadow">
                    <h3 className="font-medium mb-1">Interests</h3>
                    <p>{profileData.interests}</p>
                </div>
            </div>
            {/* The main card container div now correctly closes AFTER the content */}
        </div>
    );
}
