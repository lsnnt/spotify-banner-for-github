# spotify-banner-for-github

> A self-hosted Go server that renders your 20 recently played Spotify tracks as an SVG — perfect for embedding in your GitHub profile README.

![banner preview](https://spotify.fustin.top/?temp=foo)

---

## How it works

Most Spotify-for-README tools require you to register an OAuth app and manage a client secret. This one doesn't.

Instead, it uses the `sp_dc` session cookie that the Spotify web player already sets in your browser, and drives the full PKCE authorization flow on your behalf — no developer app registration needed. The resulting bearer token is cached in-memory with mutex-based thread safety and automatically refreshed on any non-200 response from the API.

```
sp_dc cookie → PKCE auth flow → bearer token (cached) → /v1/me/player/recently-played → SVG
```

The server also sets `Cache-Control: no-cache, no-store, must-revalidate` on every response, which is necessary to bypass [GitHub's Camo image proxy](https://github.blog/2014-01-28-proxying-user-images/) — otherwise your banner would be stale for hours.

---

## Setup

**Prerequisites:** Go 1.21+

```bash
git clone https://github.com/lsnnt/spotify-banner-for-github
cd spotify-banner-for-github
```

**1. Get your `sp_dc` cookie**

Open [open.spotify.com](https://open.spotify.com) in your browser, open DevTools → Application → Cookies, and copy the value of `sp_dc`.

**2. Create a `.env` file**

```env
SPDC="your_sp_dc_cookie_value_here"
```

**3. Build and run**

```bash
make        # or: go build . && ./spotify-banner-for-github
```

The server starts on `:8080`. Visit `http://localhost:8080/` to see your SVG.

**4. Add to your GitHub README**

Deploy to any publicly reachable server (a VPS, fly.io, etc.), then embed:

```markdown
![Spotify recently played](https://your-server.example.com/)
```

> **Note on GitHub Camo caching:** Even with correct cache headers, GitHub's image proxy can take a few minutes to reflect updates. This is a GitHub-side limitation and not specific to this project.

---

## Self-hosting

The binary is a single statically-linked executable with no external dependencies beyond the `.env` file. A minimal systemd unit or Docker setup works well.

```bash
# Example: run with systemd
ExecStart=/opt/spotify-banner-for-github/spotify-banner-for-github
EnvironmentFile=/opt/spotify-banner-for-github/.env
```

---

## A note on the `sp_dc` cookie

The `sp_dc` cookie is a long-lived session credential tied to your Spotify account. Treat it like a password:

- Never commit your `.env` file
- Rotate it by logging out and back in to open.spotify.com
- The cookie does expire; when it does, the server will log `invalid sp_dc cookie or expired` — just update your `.env` and restart

This approach is unofficial and not supported by Spotify. Use it for personal, non-commercial purposes.

---

## Contributing

The SVG design is intentionally minimal right now — contributions to make it look better are very welcome. Areas that could use work:

- Album art fetching and rendering
- Dark/light theme variants
- Artist names alongside track names
- Configurable number of tracks (currently hardcoded to 20)

Open an issue or PR if you want to take something on.

---

## License

GPL-3.0 — see [LICENSE](LICENSE).