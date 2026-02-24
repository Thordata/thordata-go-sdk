# Examples

All examples load `.env` from the repository root. Start by copying `.env.example` to `.env` and filling in your credentials:

```bash
cp .env.example .env
# Edit .env with your real credentials
```

Then run any example:

```bash
go run ./examples/proxy_basic
```

---

## Available Examples

### `proxy_basic` - Proxy Network

Demonstrates high-performance proxy requests with connection pooling.

**Requirements:**
- `THORDATA_PROXY_HOST`
- `THORDATA_PROXY_PORT`
- `THORDATA_RESIDENTIAL_USERNAME` / `THORDATA_RESIDENTIAL_PASSWORD` (or other proxy credentials)

**Run:**
```bash
go run ./examples/proxy_basic
```

---

### `serp_basic` - SERP Search

Demonstrates Google/Bing/Yandex search API.

**Requirements:**
- `THORDATA_SCRAPER_TOKEN`

**Run:**
```bash
go run ./examples/serp_basic
```

---

### `universal_basic` - Universal Scrape (Web Unlocker)

Demonstrates web scraping with JS rendering and smart waiting.

**Requirements:**
- `THORDATA_SCRAPER_TOKEN`

**Run:**
```bash
go run ./examples/universal_basic
```

---

### `tasks_basic` - Web Scraper Tasks

Demonstrates complete task lifecycle: create → wait → download.

**Requirements:**
- `THORDATA_SCRAPER_TOKEN`
- `THORDATA_PUBLIC_TOKEN`
- `THORDATA_PUBLIC_KEY`
- `THORDATA_TASK_SPIDER_ID`
- `THORDATA_TASK_SPIDER_NAME`
- `THORDATA_TASK_PARAMETERS_JSON`

**Optional:**
- `THORDATA_TASK_FILE_NAME`

**Run:**
```bash
go run ./examples/tasks_basic
```

---

### `locations_basic` - Locations API

Demonstrates querying available countries, states, cities, and ASNs.

**Requirements:**
- `THORDATA_PUBLIC_TOKEN`
- `THORDATA_PUBLIC_KEY`

**Run:**
```bash
go run ./examples/locations_basic
```

---

### `comprehensive` - Comprehensive Demo

Demonstrates all major SDK features in a single example. Automatically skips features if required credentials are missing.

**Requirements:**
- Various credentials (see code for details)

**Run:**
```bash
go run ./examples/comprehensive
```

---

### `browser_connection` - Browser Connection URL

Demonstrates how to get the WebSocket URL for connecting to Thordata Scraping Browser.

**Requirements:**
- `THORDATA_BROWSER_USERNAME`
- `THORDATA_BROWSER_PASSWORD`

**Run:**
```bash
go run ./examples/browser_connection
```

---

## Notes

- All examples use the shared `examples/internal/example` helper for env loading and client initialization.
- Examples automatically skip if required environment variables are missing.
- For integration testing, set `THORDATA_INTEGRATION=true` in `.env`.
