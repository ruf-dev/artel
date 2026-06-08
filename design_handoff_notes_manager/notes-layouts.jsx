/* ═══════════════════════════════════════════════════════════
   notes-layouts.jsx — Artel Notes Manager UI Components
   Exports: LayoutClassic, LayoutFocused, LayoutPanel,
            TypoComfortaa, TypoMono, TypoSerif
   ═══════════════════════════════════════════════════════════ */

const N = {
  sb: '#0d0d0d', ed: '#000', elev: '#141414',
  bDim: 'rgba(255,255,255,0.06)', bMid: 'rgba(255,255,255,0.11)',
  p1: 'rgba(255,255,255,0.92)', p2: 'rgba(255,255,255,0.55)',
  p3: 'rgba(255,255,255,0.32)', p4: 'rgba(255,255,255,0.18)',
  coral: '#FF4B3E', coralDim: 'rgba(255,75,62,0.08)', coralBorder: 'rgba(255,75,62,0.22)',
  font: 'Comfortaa, sans-serif', mono: 'JetBrains Mono, monospace',
  topH: 44, titleH: 46, sbW: 240, panelW: 200,
};

/* ─── ArtelMark tiny (for chips and topbar) ─────────────── */
let _mid = 0;
function AMark({ size = 20, id }) {
  const mid = React.useMemo(() => id || `am${_mid++}`, []);
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ flexShrink: 0, display: 'block' }}>
      <defs><mask id={mid}>
        <rect width="100" height="100" fill="white"/>
        <path d="M50 22L28 78h10l5.5-14h13L62 78h10ZM46.5 56L50 38l3.5 18Z" fill="black"/>
      </mask></defs>
      <circle cx="50" cy="50" r="50" fill="#FF4B3E" mask={`url(#${mid})`}/>
    </svg>
  );
}

/* ─── WikiChip — rendered wiki-link ─────────────────────── */
function WikiChip({ name }) {
  return (
    <span style={{
      display: 'inline-flex', alignItems: 'center', gap: 5,
      background: N.coralDim, border: `1px solid ${N.coralBorder}`,
      borderRadius: 5, padding: '1px 8px 1px 5px',
      fontFamily: N.font, fontSize: 12, color: N.p1,
      cursor: 'pointer', verticalAlign: 'middle', lineHeight: 1.6,
      whiteSpace: 'nowrap', userSelect: 'none',
    }}>
      <AMark size={11} />
      {name}
    </span>
  );
}

/* ─── RawLink — wiki-link in source/edit mode ────────────── */
function RawLink({ name }) {
  return (
    <span style={{
      background: 'rgba(255,75,62,0.11)', border: `1px solid rgba(255,75,62,0.25)`,
      borderRadius: 3, padding: '0 4px',
      color: 'rgba(255,160,130,0.9)', fontFamily: N.mono, fontSize: 11,
    }}>{`[[${name}]]`}</span>
  );
}

/* ─── ModeBar — Edit / Preview / Read switcher ───────────── */
function ModeBar({ active = 'preview' }) {
  const modes = [{ id: 'edit', label: 'Edit' }, { id: 'preview', label: 'Preview' }, { id: 'read', label: 'Read' }];
  return (
    <div style={{
      display: 'inline-flex', padding: 3, gap: 0,
      background: N.elev, border: `1px solid ${N.bDim}`, borderRadius: 999,
      fontFamily: N.mono, fontSize: 10, letterSpacing: '0.06em',
    }}>
      {modes.map(m => (
        <div key={m.id} style={{
          padding: '5px 13px', borderRadius: 999, whiteSpace: 'nowrap', cursor: 'pointer',
          background: m.id === active ? N.sb : 'transparent',
          color: m.id === active ? N.p1 : N.p3,
          boxShadow: m.id === active ? `inset 0 0 0 1px ${N.bMid}` : 'none',
        }}>{m.label}</div>
      ))}
    </div>
  );
}

/* ─── Sidebar micro-components ───────────────────────────── */
function SLabel({ children }) {
  return (
    <div style={{ padding: '12px 16px 5px', fontFamily: N.mono, fontSize: 9,
      letterSpacing: '0.12em', textTransform: 'uppercase', color: N.p4 }}>
      {children}
    </div>
  );
}

function TItem({ name, active, depth = 0, hasArrow, open, isFolder }) {
  const ArrowIcon = () => (
    <svg viewBox="0 0 8 8" width={8} height={8} style={{ flexShrink: 0, opacity: 0.35, transform: open ? 'rotate(90deg)' : 'none', transition: 'transform 0.15s' }}>
      <path d="M2 1.5l4 2.5-4 2.5z" fill="currentColor"/>
    </svg>
  );
  const FileIcon = () => (
    <svg viewBox="0 0 14 14" width={10} height={10} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" style={{ flexShrink: 0, opacity: 0.35 }}>
      <rect x="2" y="1" width="10" height="12" rx="1.5"/><path d="M5 4.5h4M5 7h4M5 9.5h2"/>
    </svg>
  );
  const FolderIcon = () => (
    <svg viewBox="0 0 14 14" width={10} height={10} fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" style={{ flexShrink: 0, opacity: 0.45 }}>
      <path d="M1 3.5h4l1.5 1.5H13v7H1z"/>
    </svg>
  );
  return (
    <div style={{
      display: 'flex', alignItems: 'center', gap: 5,
      padding: `4px 16px 4px ${16 + depth * 12}px`,
      fontSize: 12, fontFamily: N.font, cursor: 'pointer',
      color: active ? N.p1 : N.p2,
      background: active ? N.coralDim : 'transparent',
      borderLeft: active ? `2px solid ${N.coral}` : '2px solid transparent',
    }}>
      {hasArrow ? <ArrowIcon /> : <span style={{ width: 8, flexShrink: 0 }}/>}
      {isFolder ? <FolderIcon /> : <FileIcon />}
      <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1 }}>{name}</span>
    </div>
  );
}

function TagPill({ tag }) {
  return (
    <span style={{
      padding: '3px 8px', background: 'rgba(255,255,255,0.04)',
      border: `1px solid ${N.bDim}`, borderRadius: 4,
      fontSize: 10, fontFamily: N.mono, color: N.p3, cursor: 'pointer', letterSpacing: '0.06em',
    }}>{tag}</span>
  );
}

/* ─── Mini knowledge graph ───────────────────────────────── */
function MiniGraph() {
  const nodes = [
    { x: 85, y: 52, r: 5, a: true }, { x: 158, y: 36 }, { x: 36, y: 84 },
    { x: 142, y: 92 }, { x: 198, y: 70 }, { x: 102, y: 108 },
  ];
  const edges = [[0,1],[0,2],[0,3],[1,3],[1,4],[2,5],[3,5]];
  return (
    <svg width="100%" height="76" viewBox="0 0 236 120" style={{ display: 'block' }}>
      {edges.map(([a,b],i) => <line key={i} x1={nodes[a].x} y1={nodes[a].y} x2={nodes[b].x} y2={nodes[b].y} stroke="rgba(255,255,255,0.08)" strokeWidth="1.5"/>)}
      {nodes.map((n,i) => <circle key={i} cx={n.x} cy={n.y} r={n.r||3} fill={n.a ? N.coral : 'rgba(255,255,255,0.22)'} style={n.a ? { filter: `drop-shadow(0 0 4px ${N.coral})` } : {}}/>)}
    </svg>
  );
}

/* ─── Search bar ─────────────────────────────────────────── */
function SearchBar() {
  return (
    <div style={{
      margin: '8px 10px', display: 'flex', alignItems: 'center', gap: 8,
      background: 'rgba(255,255,255,0.035)', border: `1px solid ${N.bDim}`,
      borderRadius: 6, padding: '7px 10px',
    }}>
      <svg viewBox="0 0 16 16" width={12} height={12} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" style={{ color: N.p4, flexShrink: 0 }}>
        <circle cx="6.5" cy="6.5" r="4.5"/><path d="M10.5 10.5l3 3"/>
      </svg>
      <span style={{ fontSize: 11, fontFamily: N.mono, color: N.p4, letterSpacing: '0.04em' }}>Search notes…</span>
    </div>
  );
}

/* ─── Divider ────────────────────────────────────────────── */
function Div() { return <div style={{ margin: '5px 16px', height: 1, background: N.bDim, flexShrink: 0 }}/>; }

/* ═══════════════════════════════════════════════════════════
   FULL SIDEBAR
   ═══════════════════════════════════════════════════════════ */
function NoteSidebar({ slim }) {
  const w = slim ? 220 : N.sbW;
  return (
    <div style={{ width: w, background: N.sb, borderRight: `1px solid ${N.bDim}`, display: 'flex', flexDirection: 'column', flexShrink: 0, height: '100%', overflow: 'hidden' }}>
      <SearchBar />
      <div style={{ flex: 1, overflowY: 'auto', overflowX: 'hidden' }}>
        <SLabel>★ Favorites</SLabel>
        <TItem name="Project Brainstorm" active />
        <TItem name="Meeting Notes" />
        <Div />
        <SLabel>⏱ Recent</SLabel>
        <TItem name="API Design" />
        <TItem name="Research v2" />
        <TItem name="Quick Capture" />
        <Div />
        <SLabel>All Notes</SLabel>
        <TItem name="Projects" hasArrow open isFolder />
        <TItem name="Project Brainstorm" active depth={1} />
        <TItem name="Meeting Notes" depth={1} />
        <TItem name="Team Roadmap" depth={1} />
        <TItem name="Research" hasArrow isFolder />
        <TItem name="API Design" depth={1} />
        <TItem name="Competitor Analysis" depth={1} />
        <TItem name="Quick Capture" />
        <Div />
        <SLabel># Tags</SLabel>
        <div style={{ padding: '6px 14px 10px', display: 'flex', flexWrap: 'wrap', gap: 5 }}>
          {['#design','#api','#backlog','#research','#agent'].map(t => <TagPill key={t} tag={t}/>)}
        </div>
        {!slim && <>
          <Div />
          <SLabel>◉ Graph</SLabel>
          <div style={{ padding: '2px 10px 10px' }}><MiniGraph /></div>
        </>}
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   ICON SIDEBAR (Layout B)
   ═══════════════════════════════════════════════════════════ */
function IconSidebar() {
  const icons = [
    { id: 'search', svg: <svg viewBox="0 0 16 16" width={16} height={16} fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round"><circle cx="6.5" cy="6.5" r="4.5"/><path d="M10.5 10.5l3 3"/></svg> },
    { id: 'tree',   svg: <svg viewBox="0 0 16 16" width={16} height={16} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"><path d="M2 3h4M2 8h8M2 13h5"/><circle cx="13" cy="3" r="1.5"/><circle cx="13" cy="8" r="1.5"/><circle cx="13" cy="13" r="1.5"/><path d="M6 3h5M10 8h1M7 13h4"/></svg> },
    { id: 'star',   svg: <svg viewBox="0 0 16 16" width={15} height={15} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"><path d="M8 2l1.6 3.4L14 6.2l-3 2.8.7 4-3.7-2-3.7 2 .7-4-3-2.8 4.4-.8z"/></svg> },
    { id: 'tag',    svg: <svg viewBox="0 0 16 16" width={15} height={15} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"><path d="M9 2H4L2 4v5.5L9 16l6-6-6-8z"/><circle cx="5.5" cy="5.5" r="1"/></svg> },
    { id: 'graph',  svg: <svg viewBox="0 0 16 16" width={15} height={15} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"><circle cx="8" cy="8" r="1.5"/><circle cx="3" cy="3.5" r="1.5"/><circle cx="13" cy="3.5" r="1.5"/><circle cx="3" cy="12.5" r="1.5"/><circle cx="13" cy="12.5" r="1.5"/><path d="M8 8L3 3.5M8 8l5-4.5M8 8l-5 4.5M8 8l5 4.5"/></svg> },
  ];
  return (
    <div style={{ width: 44, background: N.sb, borderRight: `1px solid ${N.bDim}`, display: 'flex', flexDirection: 'column', alignItems: 'center', padding: '10px 0', gap: 4, flexShrink: 0 }}>
      <div style={{ marginBottom: 10 }}><AMark size={22} id="icon-sb-logo" /></div>
      {icons.map((it, i) => (
        <div key={it.id} style={{
          width: 32, height: 32, borderRadius: 6, display: 'flex', alignItems: 'center', justifyContent: 'center',
          color: i === 1 ? N.p1 : N.p3,
          background: i === 1 ? 'rgba(255,255,255,0.07)' : 'transparent', cursor: 'pointer',
        }}>{it.svg}</div>
      ))}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   NOTE TITLE BAR
   ═══════════════════════════════════════════════════════════ */
function TitleBar({ mode }) {
  return (
    <div style={{ height: N.titleH, background: N.ed, borderBottom: `1px solid ${N.bDim}`, display: 'flex', alignItems: 'center', padding: '0 32px', gap: 16, flexShrink: 0 }}>
      <div style={{ flex: 1 }}>
        <div style={{ fontSize: 17, fontWeight: 600, fontFamily: N.font, color: N.p1, letterSpacing: '-0.01em' }}>Project Brainstorm</div>
        <div style={{ fontSize: 10, fontFamily: N.mono, color: N.p4, marginTop: 1 }}>edited 2h ago · 480 words · 6 links</div>
      </div>
      <ModeBar active={mode} />
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   NOTE CONTENT — Preview / WYSIWYG (rendered markdown)
   ═══════════════════════════════════════════════════════════ */
function NoteRendered({ bodyFont, bodySize = 14, lineH = 1.72 }) {
  const bf = bodyFont || N.font;
  const h1 = { fontSize: 26, fontWeight: 700, fontFamily: N.font, color: N.p1, letterSpacing: '-0.02em', margin: '0 0 14px' };
  const h2 = { fontSize: 16, fontWeight: 600, fontFamily: N.font, color: N.p1, letterSpacing: '-0.01em', margin: '26px 0 10px' };
  const body = { fontSize: bodySize, fontFamily: bf, color: N.p2, lineHeight: lineH, margin: '0 0 14px' };
  const bq = { borderLeft: '2px solid rgba(255,75,62,0.4)', paddingLeft: 14, margin: '0 0 14px', fontSize: bodySize - 0.5, fontFamily: bf, color: N.p3, fontStyle: 'italic', lineHeight: lineH };
  const codeBlock = { background: '#0d0d0d', border: `1px solid ${N.bDim}`, borderRadius: 6, padding: '12px 16px', margin: '12px 0 16px', fontFamily: N.mono, fontSize: 11, color: N.p2, lineHeight: 1.65, whiteSpace: 'pre', overflow: 'hidden', display: 'block' };
  const tasks = [
    { done: true,  text: 'Design the sidebar layout' },
    { done: true,  text: 'Define wiki-link chip component' },
    { done: false, text: 'Build ', link: 'API Design', after: ' module' },
    { done: false, text: 'Review ', link: 'Competitor Analysis' },
    { done: false, text: 'Sync with ', link: 'Team Roadmap' },
  ];
  return (
    <div style={{ flex: 1, overflowY: 'auto', padding: '30px 40px' }}>
      <div style={h1}>Project Brainstorm</div>
      <div style={bq}>Central hub for all product ideas and team discussions.</div>
      <div style={h2}>Core Concept</div>
      <p style={body}>The artel notes system is built around connected knowledge. Every note links to others through wiki-links, creating a living graph of your team's thinking.</p>
      <p style={body}>Key integrations with <WikiChip name="Meeting Notes" /> and <WikiChip name="Research" /> feed into this document automatically via agent writes.</p>
      <div style={h2}>Open Tasks</div>
      <div style={{ marginBottom: 18 }}>
        {tasks.map((t, i) => (
          <div key={i} style={{ display: 'flex', alignItems: 'flex-start', gap: 9, marginBottom: 7, fontSize: bodySize, fontFamily: bf, color: t.done ? N.p4 : N.p2, lineHeight: lineH, textDecoration: t.done ? 'line-through' : 'none' }}>
            <div style={{ width: 14, height: 14, borderRadius: 3, flexShrink: 0, marginTop: 3, border: `1.5px solid ${t.done ? N.coral : N.bMid}`, background: t.done ? N.coral : 'transparent', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              {t.done && <svg viewBox="0 0 10 10" width={8} height={8}><path d="M2 5l2.5 2.5L8 3" stroke="white" strokeWidth="1.6" strokeLinecap="round" fill="none"/></svg>}
            </div>
            <span>{t.text}{t.link && <WikiChip name={t.link} />}{t.after}</span>
          </div>
        ))}
      </div>
      <div style={h2}>Tech Stack</div>
      <code style={codeBlock}>{`interface Note {\n  id: string\n  title: string\n  links: WikiLink[]\n  tags: string[]\n}`}</code>
      <div style={h2}>See Also</div>
      <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
        {['Meeting Notes','Research','API Design','Backlog','Team Roadmap'].map(n => <WikiChip key={n} name={n} />)}
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   NOTE CONTENT — Raw markdown (source edit)
   ═══════════════════════════════════════════════════════════ */
function NoteRaw() {
  const L = ({ ch, h1, h2, dim, muted, children }) => (
    <div style={{ fontFamily: N.mono, lineHeight: 1.75, color: h1 ? N.p1 : h2 ? N.p1 : dim ? N.p4 : muted ? N.p3 : N.p2, fontWeight: h1 ? 700 : h2 ? 500 : 400, fontSize: h1 ? 15 : h2 ? 12 : 12 }}>{children}</div>
  );
  return (
    <div style={{ flex: 1, overflowY: 'auto', padding: '28px 32px' }}>
      <L h1># Project Brainstorm</L>
      <L dim>&nbsp;</L>
      <L muted>&gt; Central hub for all product ideas and team discussions.</L>
      <L dim>&nbsp;</L>
      <L h2>## Core Concept</L>
      <L dim>&nbsp;</L>
      <div style={{ fontFamily: N.mono, fontSize: 12, lineHeight: 1.75, color: N.p2 }}>
        The artel notes system is built around connected knowledge.
        Every note links through <RawLink name="Meeting Notes" /> and <RawLink name="Research" />.
      </div>
      <L dim>&nbsp;</L>
      <L h2>## Open Tasks</L>
      <L dim>&nbsp;</L>
      <L muted>- [x] Design the sidebar layout</L>
      <L muted>- [x] Define wiki-link chip component</L>
      <div style={{ fontFamily: N.mono, fontSize: 12, lineHeight: 1.75, color: N.p2 }}>- [ ] Build <RawLink name="API Design" /> module</div>
      <div style={{ fontFamily: N.mono, fontSize: 12, lineHeight: 1.75, color: N.p2 }}>- [ ] Review <RawLink name="Competitor Analysis" /></div>
      <div style={{ fontFamily: N.mono, fontSize: 12, lineHeight: 1.75, color: N.p2 }}>- [ ] Sync with <RawLink name="Team Roadmap" /></div>
      <L dim>&nbsp;</L>
      <L h2>## Tech Stack</L>
      <L dim>&nbsp;</L>
      {['```ts','interface Note {','  id: string','  title: string','  links: WikiLink[]','  tags: string[]','}','```'].map((l,i) => (
        <div key={i} style={{ fontFamily: N.mono, fontSize: 11, lineHeight: 1.6, color: i===0||i===7 ? N.p4 : N.p3, background: i>0&&i<7 ? 'rgba(255,255,255,0.02)' : 'transparent', paddingLeft: i>0&&i<7 ? 16 : 0 }}>{l}</div>
      ))}
      <L dim>&nbsp;</L>
      <L h2>## See Also</L>
      <L dim>&nbsp;</L>
      <div style={{ fontFamily: N.mono, fontSize: 12, lineHeight: 1.75, color: N.p2, display: 'flex', flexWrap: 'wrap', gap: 8, alignItems: 'center' }}>
        {['Meeting Notes','Research','API Design','Backlog'].map((n,i) => [
          <RawLink key={n} name={n} />,
          i < 3 && <span key={`d${i}`} style={{ color: N.p4 }}>·</span>
        ])}
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   RIGHT PANEL (Layout C)
   ═══════════════════════════════════════════════════════════ */
function RightPanel() {
  return (
    <div style={{ width: N.panelW, background: N.sb, borderLeft: `1px solid ${N.bDim}`, flexShrink: 0, overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
      <SLabel>Outline</SLabel>
      {[['Core Concept',0],['Open Tasks',0],['Tech Stack',0],['See Also',0]].map(([n,d],i) => (
        <div key={i} style={{ padding: `4px 16px 4px ${16+d*10}px`, fontSize: 11, fontFamily: N.font, color: N.p3, cursor: 'pointer' }}>— {n}</div>
      ))}
      <Div />
      <SLabel>Backlinks</SLabel>
      {[['Meeting Notes','mentions this note'],['Research v2','links to project'],['Team Roadmap','references tasks']].map(([n,c],i) => (
        <div key={i} style={{ padding: '6px 16px', cursor: 'pointer' }}>
          <div style={{ fontSize: 11, fontFamily: N.font, color: N.p2 }}>{n}</div>
          <div style={{ fontSize: 9, fontFamily: N.mono, color: N.p4, marginTop: 1 }}>{c}</div>
        </div>
      ))}
      <Div />
      <SLabel>Properties</SLabel>
      <div style={{ padding: '0 16px 10px' }}>
        {[['created','2026-05-12'],['modified','2h ago'],['words','480'],['links','6']].map(([k,v]) => (
          <div key={k} style={{ display: 'flex', justifyContent: 'space-between', padding: '3px 0', fontSize: 10, fontFamily: N.mono }}>
            <span style={{ color: N.p4 }}>{k}</span>
            <span style={{ color: N.p3 }}>{v}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   TOP BAR shared
   ═══════════════════════════════════════════════════════════ */
function TopBar({ breadcrumb, modeBar, extra, id }) {
  return (
    <div style={{ height: N.topH, background: '#0a0a0a', borderBottom: `1px solid ${N.bDim}`, display: 'flex', alignItems: 'center', padding: '0 16px', gap: 10, flexShrink: 0 }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        <AMark size={22} id={id} />
        <span style={{ fontWeight: 700, fontSize: 14, letterSpacing: '-0.02em', fontFamily: N.font, color: N.p1 }}>artel</span>
        <span style={{ fontFamily: N.mono, fontSize: 10, color: N.p4, paddingLeft: 8, borderLeft: `1px solid ${N.bDim}`, letterSpacing: '0.06em' }}>notes</span>
      </div>
      {breadcrumb && <span style={{ fontFamily: N.mono, fontSize: 11, color: N.p4, marginLeft: 8 }}>{breadcrumb}</span>}
      <div style={{ flex: 1 }} />
      {modeBar}
      {extra}
      <div style={{ display: 'flex', alignItems: 'center', gap: 5, padding: '5px 10px', background: 'rgba(255,255,255,0.04)', border: `1px solid ${N.bDim}`, borderRadius: 6, fontFamily: N.mono, fontSize: 11, color: N.p3, cursor: 'pointer' }}>
        <svg viewBox="0 0 12 12" width={10} height={10} fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round"><path d="M6 1v10M1 6h10"/></svg> New
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   LAYOUT A — Classic (Live Preview / WYSIWYG)
   ═══════════════════════════════════════════════════════════ */
function LayoutClassic() {
  return (
    <div style={{ width: 1280, height: 800, background: N.ed, display: 'flex', flexDirection: 'column', fontFamily: N.font, color: N.p1, overflow: 'hidden' }}>
      <TopBar id="tb-a" />
      <div style={{ flex: 1, display: 'flex', overflow: 'hidden' }}>
        <NoteSidebar />
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
          <TitleBar mode="preview" />
          <NoteRendered />
        </div>
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   LAYOUT B — Focused (Raw / Source Edit, icon sidebar)
   ═══════════════════════════════════════════════════════════ */
function LayoutFocused() {
  return (
    <div style={{ width: 1280, height: 800, background: N.ed, display: 'flex', fontFamily: N.font, color: N.p1, overflow: 'hidden' }}>
      <IconSidebar />
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* Breadcrumb bar */}
        <div style={{ height: N.topH, background: '#0a0a0a', borderBottom: `1px solid ${N.bDim}`, display: 'flex', alignItems: 'center', padding: '0 20px', gap: 6, flexShrink: 0 }}>
          {['artel','Projects','Project Brainstorm'].map((seg, i, arr) => [
            <span key={seg} style={{ fontFamily: N.mono, fontSize: 11, color: i === arr.length-1 ? N.p1 : N.p4 }}>{seg}</span>,
            i < arr.length-1 && <span key={`s${i}`} style={{ fontFamily: N.mono, fontSize: 11, color: N.p4 }}>›</span>
          ])}
          <div style={{ flex: 1 }} />
          <ModeBar active="edit" />
        </div>
        {/* Editor with line numbers */}
        <div style={{ flex: 1, display: 'flex', overflow: 'hidden' }}>
          <div style={{ width: 38, background: '#080808', borderRight: `1px solid ${N.bDim}`, padding: '28px 0', flexShrink: 0 }}>
            {Array.from({ length: 26 }, (_, i) => (
              <div key={i} style={{ height: 21, display: 'flex', alignItems: 'center', justifyContent: 'flex-end', paddingRight: 10, fontFamily: N.mono, fontSize: 10, color: N.p4 }}>{i + 1}</div>
            ))}
          </div>
          <NoteRaw />
        </div>
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   LAYOUT C — Panel (Read mode + right panel)
   ═══════════════════════════════════════════════════════════ */
function LayoutPanel() {
  return (
    <div style={{ width: 1280, height: 800, background: N.ed, display: 'flex', flexDirection: 'column', fontFamily: N.font, color: N.p1, overflow: 'hidden' }}>
      <TopBar id="tb-c" breadcrumb="Projects  ›  Project Brainstorm" modeBar={<ModeBar active="read" />} />
      <div style={{ flex: 1, display: 'flex', overflow: 'hidden' }}>
        <NoteSidebar slim />
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
          <NoteRendered />
        </div>
        <RightPanel />
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════════
   TYPOGRAPHY VARIANTS
   ═══════════════════════════════════════════════════════════ */
function TypoSample({ meta, bodyFont, bodySize = 13, lineH = 1.72 }) {
  const bf = bodyFont;
  return (
    <div style={{ width: 460, height: 340, background: N.ed, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
      <div style={{ height: 34, background: '#0a0a0a', borderBottom: `1px solid ${N.bDim}`, display: 'flex', alignItems: 'center', padding: '0 18px', flexShrink: 0 }}>
        <span style={{ fontFamily: N.mono, fontSize: 9, color: N.p4, letterSpacing: '0.1em', textTransform: 'uppercase' }}>{meta}</span>
      </div>
      <div style={{ flex: 1, padding: '20px 24px', overflow: 'hidden' }}>
        <div style={{ fontSize: 20, fontWeight: 700, fontFamily: N.font, color: N.p1, letterSpacing: '-0.02em', marginBottom: 10 }}>Project Brainstorm</div>
        <div style={{ borderLeft: '2px solid rgba(255,75,62,0.38)', paddingLeft: 12, marginBottom: 14, fontSize: bodySize - 0.5, fontFamily: bf, color: N.p3, fontStyle: 'italic', lineHeight: lineH }}>
          Central hub for all product ideas.
        </div>
        <div style={{ fontSize: 14, fontWeight: 600, fontFamily: N.font, color: N.p1, marginBottom: 8 }}>Core Concept</div>
        <div style={{ fontSize: bodySize, fontFamily: bf, color: N.p2, lineHeight: lineH, marginBottom: 14 }}>
          The notes system is built around connected knowledge. Every note links through <WikiChip name="Meeting Notes" /> and <WikiChip name="Research" />.
        </div>
        <div style={{ fontSize: 14, fontWeight: 600, fontFamily: N.font, color: N.p1, marginBottom: 8 }}>Open Tasks</div>
        <div style={{ fontSize: bodySize, fontFamily: bf, color: N.p3, lineHeight: lineH }}>
          — Design sidebar layout ✓<br />— Build <WikiChip name="API Design" /> module
        </div>
      </div>
    </div>
  );
}

function TypoComfortaa() { return <TypoSample meta="Comfortaa · brand font · warm rounded sans" bodyFont="Comfortaa, sans-serif" bodySize={13} lineH={1.7} />; }
function TypoMono()      { return <TypoSample meta="JetBrains Mono · code-native · focused" bodyFont="JetBrains Mono, monospace" bodySize={12} lineH={1.78} />; }
function TypoSerif()     { return <TypoSample meta="Georgia · editorial serif · warm + readable" bodyFont="Georgia, 'Times New Roman', serif" bodySize={14} lineH={1.8} />; }

Object.assign(window, { LayoutClassic, LayoutFocused, LayoutPanel, TypoComfortaa, TypoMono, TypoSerif });
