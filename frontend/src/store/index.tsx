// BACKEND_URL points at the labbaika Go backend (any of its cashbook-svc /
// payslip-svc / etc subcommands — they all listen on :8080).
//
// REACT_APP_BACKEND_URL is set at BUILD time, not runtime: CRA inlines
// REACT_APP_* into the bundle, so changing this after `npm run build` has no
// effect.
//   - local dev: unset, falls back to the CURRENT page's own hostname on
//                :8080. Not a hardcoded "localhost" — that would resolve to
//                whatever device is viewing the page. Opening the dev server
//                from another device on the LAN (CRA prints an "On Your
//                Network" URL for exactly this) means "localhost" means THAT
//                device, not this Mac, so the backend call would 404/refuse
//                with nothing listening there. window.location.hostname is
//                whatever the browser actually used to reach this page
//                (localhost, 127.0.0.1, or a LAN IP), so it's always correct.
//   - deployed:  "" (empty), so requests go same-origin and the reverse proxy
//                routes /api and /receipts to the backend. Keeps the browser
//                off cross-origin, so no CORS and no mixed-content on HTTPS.
export const BACKEND_URL =
  process.env.REACT_APP_BACKEND_URL ?? `${window.location.protocol}//${window.location.hostname}:8080`;
