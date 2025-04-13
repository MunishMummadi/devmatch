'use client';

import Link from 'next/link';
import {
    SignInButton,
    SignUpButton,
    SignedIn,
    SignedOut,
    UserButton,
} from '@clerk/nextjs';

export default function Header() {
    return (
        <header style={styles.header}>
            <nav style={styles.nav}>
                <div style={styles.logo}>
                    <Link href="/" style={{ color: 'white', textDecoration: 'none' }}>
                        DevMatch
                    </Link>
                </div>

                <ul style={styles.navLinks}>
                    <li><Link href="/messages" style={styles.link}>Messages</Link></li>
                </ul>

                <div style={styles.actions}>
                    <SignedOut>
                        <SignInButton mode="modal">
                            <button style={styles.button}>Login</button>
                        </SignInButton>
                        <SignUpButton mode="modal">
                            <button style={styles.button}>Signup</button>
                        </SignUpButton>
                    </SignedOut>

                    <SignedIn>
                        <UserButton />
                    </SignedIn>

                    <button style={{ ...styles.button, ...styles.upgradeBtn }}>
                        Upgrade to Pro
                    </button>
                </div>
            </nav>
        </header>
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
    link: {
        color: 'white',
        textDecoration: 'none',
        fontWeight: '500',
    },
    actions: {
        display: 'flex',
        gap: '1rem',
        alignItems: 'center',
    },
    button: {
        background: 'transparent',
        border: '1px solid #ffffff',
        color: '#ffffff',
        padding: '0.5rem 1rem',
        cursor: 'pointer',
        borderRadius: '6px',
    },
    upgradeBtn: {
        backgroundColor: '#ff9800',
        borderColor: '#ff9800',
    },
};
