// app/page.js
import Link from 'next/link';

export default function HomePage() {
  return (
    <section style={styles.heroSection}>
      <h1>Who is Impact for?</h1>
      <p>
        Whether you’re a programmer, designer, freelancer or you need Impact for a
        whole team, our pricing just makes sense.
      </p>

      {/* Below could be your "feature" columns */}
      <div style={styles.featuresContainer}>
        <div style={styles.featureBox}>
          <h2>Marketing</h2>
          <p>Reveal best strategies from the market and your competitors.</p>
        </div>
        <div style={styles.featureBox}>
          <h2>Research</h2>
          <p>Understand your market, your competitors, and your customers.</p>
        </div>
        <div style={styles.featureBox}>
          <h2>Sales</h2>
          <p>Enhance performance throughout your sales funnel.</p>
        </div>
      </div>
    </section>
  );
}

const styles = {
  heroSection: {
    background: '#0052cc',
    color: '#fff',
    textAlign: 'center',
    padding: '4rem 1rem',
  },
  featuresContainer: {
    display: 'flex',
    flexDirection: 'column',
    gap: '2rem',
    marginTop: '2rem',
  },
  featureBox: {
    backgroundColor: '#fff',
    color: '#333',
    borderRadius: '8px',
    padding: '1.5rem',
    margin: '0 auto',
    maxWidth: '300px',
  },
};
