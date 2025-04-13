"use client"

import './globals.css';
import { ClerkProvider } from '@clerk/nextjs';

// export const metadata = {
//   title: 'DevMatch',
//   description: 'Developer matchmaking dashboard',
// };

export default function RootLayout({ children }) {
  return (
    <ClerkProvider>
      <html lang="en">
        <body>
          {children}
        </body>
      </html>
    </ClerkProvider>
  );
}
