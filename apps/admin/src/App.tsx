import { FOUNDATION_AREAS, PRESENTATION_LAYERS } from "./contracts";
import "./styles.css";

export function App() {
  return (
    <main className="shell">
      <header className="header">
        <div>
          <p className="eyebrow">GRIDWORKS</p>
          <h1>Administrative foundation</h1>
        </div>
        <span className="status">WP-001 foundation</span>
      </header>
      <section className="intro" aria-labelledby="intro-title">
        <p className="eyebrow">Control plane</p>
        <h2 id="intro-title">Boundaries before operations</h2>
        <p>
          This application reserves the administrative, moderation and content-publication surface.
          Authoritative domain state remains server-owned and presentation remains non-authoritative.
        </p>
      </section>
      <section className="grid" aria-label="Foundation areas">
        <article className="card">
          <h2>Domain ownership</h2>
          <ul>
            {FOUNDATION_AREAS.map((area) => (
              <li key={area.id}>
                <strong>{area.id}</strong>
                <span>{area.owner}</span>
              </li>
            ))}
          </ul>
        </article>
        <article className="card">
          <h2>Presentation layers</h2>
          <ul>
            {PRESENTATION_LAYERS.map((layer) => (
              <li key={layer}>
                <strong>{layer}</strong>
                <span>server domain state</span>
              </li>
            ))}
          </ul>
        </article>
      </section>
    </main>
  );
}
