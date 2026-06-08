/* Shared bits across all variants — logo, mini icons, helpers.
   Exposes onto window so each variant <script> file can use them. */

function ArtelMark({ size = 40, animated = true }) {
  // A coral disc with a stylized "А" carved out — fits the Zpotify
  // visual DNA (one solid mark = brand) but distinctly Artel.
  return (
    <svg
      viewBox="0 0 100 100"
      style={{
        width: size, height: size,
        display: 'block', flexShrink: 0,
        animation: animated ? 'artelPulse 2.6s ease-in-out infinite' : 'none',
      }}
    >
      <defs>
        <mask id="artel-A-cut">
          <rect width="100" height="100" fill="white" />
          {/* The "А" cut-out — bold geometric */}
          <path
            d="M 50 22 L 28 78 L 38 78 L 43.5 64 L 56.5 64 L 62 78 L 72 78 Z M 46.5 56 L 50 38 L 53.5 56 Z"
            fill="black"
          />
        </mask>
      </defs>
      <circle cx="50" cy="50" r="50" fill="#FF4B3E" mask="url(#artel-A-cut)" />
    </svg>
  );
}

function ArtelWordmark({ size = 22, color = 'var(--text-primary)' }) {
  return (
    <span style={{
      fontFamily: 'Comfortaa, sans-serif',
      fontWeight: 700,
      fontSize: size,
      letterSpacing: '-0.02em',
      color,
      lineHeight: 1,
    }}>artel</span>
  );
}

function CoralDot({ size = 6 }) {
  return (
    <span style={{
      display: 'inline-block',
      width: size, height: size, borderRadius: '50%',
      background: '#FF4B3E',
      boxShadow: '0 0 8px rgba(255,75,62,0.55)',
      verticalAlign: 'middle',
    }}/>
  );
}

/* tiny iconography — outlines only, white-based. Match Zpotify's no-emoji rule. */
const Icon = {
  mail: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 7l9 6 9-6"/>
    </svg>
  ),
  box: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 7l9-4 9 4-9 4-9-4z"/><path d="M3 7v10l9 4 9-4V7"/><path d="M12 11v10"/>
    </svg>
  ),
  brain: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 4a3 3 0 0 0-3 3v1a3 3 0 0 0-2 5 3 3 0 0 0 2 5v1a3 3 0 0 0 6 0V4a2 2 0 0 0-3 0z"/>
      <path d="M15 4a3 3 0 0 1 3 3v1a3 3 0 0 1 2 5 3 3 0 0 1-2 5v1a3 3 0 0 1-6 0"/>
    </svg>
  ),
  flask: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 3h6"/><path d="M10 3v6L4 19a2 2 0 0 0 2 3h12a2 2 0 0 0 2-3l-6-10V3"/>
    </svg>
  ),
  tg: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="currentColor">
      <path d="M22 3.5L2.5 11l5 1.7L18 6 9.7 14l-.4 5.5 3.4-3.2 5.5 4c.9.5 1.5.2 1.7-.8L22 4.3c.2-.7-.2-1.1-.9-.8z"/>
    </svg>
  ),
  arrow: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||16} height={p.size||16} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
      <path d="M5 12h14M13 5l7 7-7 7"/>
    </svg>
  ),
  check: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||14} height={p.size||14} fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M5 12l4 4L19 6"/>
    </svg>
  ),
  obsidian: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 2L4 7l3 13h10l3-13z"/><path d="M9 7l3-5 3 5"/><path d="M7 20l5-9 5 9"/>
    </svg>
  ),
  chart: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||18} height={p.size||18} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 3v18h18"/><path d="M7 14l4-4 3 3 5-7"/>
    </svg>
  ),
  alert: (p={}) => (
    <svg viewBox="0 0 24 24" width={p.size||16} height={p.size||16} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 2L2 21h20L12 2z"/><path d="M12 9v5"/><circle cx="12" cy="17" r=".8" fill="currentColor"/>
    </svg>
  ),
};

/* Section wrapper — gives each variant artboard a controlled column */
function PageScaffold({ children, bg = 'var(--bg-base)', maxW = 1280 }) {
  return (
    <div style={{
      width: '100%', background: bg, color: 'var(--text-primary)',
      fontFamily: 'Comfortaa, sans-serif', overflow: 'hidden',
    }}>
      {children}
    </div>
  );
}

Object.assign(window, { ArtelMark, ArtelWordmark, CoralDot, Icon, PageScaffold });
