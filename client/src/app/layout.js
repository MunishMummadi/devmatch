// app/layout.js
import './globals.css';
import React from 'react';

// If you need a special font from Google, you can import it:
// import { Inter } from 'next/font/google';
// const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'DevMatch Dashboard',
  description: 'A sample Next.js 13 dashboard layout',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {/* NAVIGATION / HEADER */}
        <header style={styles.header}>
          <nav style={styles.nav}>
            <div style={styles.logo}>impact</div>
            <ul style={styles.navLinks}>
              <li>Messages</li>
              {/* <li>Front pages</li>
              <li>App pages</li>
              <li>Support</li> */}
            </ul>
            <div style={styles.actions}>
              {/* <button style={styles.button}>Login</button> */}
              <button style={styles.button}>Signup</button>
              <button style={{ ...styles.button, ...styles.upgradeBtn }}>
                Upgrade to Pro
              </button>
            </div>
          </nav>
        </header>

        {/* MAIN CONTENT: children come from each route's page.js */}
        <main>{children}</main>
      </body>
    </html>
  );
}

const styles = {
  header: {
    backgroundColor: '#0052cc',
    color: '#fff',
    padding: '1rem 2rem',
  },
  nav: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  logo: {
    fontWeight: 'bold',
    fontSize: '1.5rem',
    textTransform: 'uppercase',
  },
  navLinks: {
    listStyle: 'none',
    display: 'flex',
    gap: '1rem',
  },
  actions: {
    display: 'flex',
    gap: '1rem',
  },
  button: {
    background: 'transparent',
    border: '1px solid #ffffff',
    color: '#ffffff',
    padding: '0.5rem 1rem',
    cursor: 'pointer',
  },
  upgradeBtn: {
    backgroundColor: '#ff9800',
    borderColor: '#ff9800',
  },
};
