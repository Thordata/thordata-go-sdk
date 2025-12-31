# Examples

All examples load `.env` from repo root. Start by copying `.env.example` to `.env`.

```bash
cp .env.example .env
```

Then run any example:

```bash
go run ./examples/serp_basic
```

## SERP
Requires:
- THORDATA_SCRAPER_TOKEN

Run:
```bash
go run ./examples/serp_basic
go run ./examples/serp_google_news
```

## Universal Scraper
Requires:
- THORDATA_SCRAPER_TOKEN

Run:
```bash
go run ./examples/universal_basic
```

## Locations API
Requires:
- THORDATA_SCRAPER_TOKEN
- THORDATA_PUBLIC_TOKEN
- THORDATA_PUBLIC_KEY

Run:
```bash
go run ./examples/locations_basic
```

## Web Scraper Tasks (create -> wait -> download)
Requires:
- THORDATA_SCRAPER_TOKEN
- THORDATA_PUBLIC_TOKEN
- THORDATA_PUBLIC_KEY
- THORDATA_SPIDER_ID
- THORDATA_SPIDER_NAME
- THORDATA_TASK_PARAMETERS_JSON

Optional:
- THORDATA_TASK_FILE_NAME

Run:
```bash
go run ./examples/tasks_basic
```

## Proxy Network
Requires:
- THORDATA_PROXY_HOST
- THORDATA_PROXY_PORT

Run:
```bash
go run ./examples/proxy_basic
```