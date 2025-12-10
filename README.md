# Summon MCP Server

Remote MCP server for weight-space navigation via stance declaration. Deployed as a Cloudflare Worker with GitHub OAuth authentication.

## What It Does

Provides a single `summon` tool that:
- Takes three parameters: `who`, `where`, `doing`
- Returns: `"You are {who} at {where} doing {doing}"`
- Logs all invocations with user context and timestamp

## Observability

Logs are structured JSON:
```json
{
  "type": "auth",
  "user": "github_login",
  "timestamp": "2025-03-10T12:00:00.000Z"
}
```

```json
{
  "type": "summon",
  "user": "github_login",
  "timestamp": "2025-03-10T12:00:00.000Z",
  "stance": {
    "who": "Gwern",
    "where": "gwern.net",
    "doing": "being thorough"
  }
}
```

View logs via Cloudflare dashboard or `wrangler tail`.

## Architecture

- **Durable Objects**: Persistent MCP server state
- **OAuth Provider**: GitHub authentication flow
- **Workers MCP**: SSE/Streamable-HTTP protocol
- **KV**: OAuth token storage
- **Observability**: Structured logging to Cloudflare

## Tool Schema

See the Rust reference implementation in `../summon_mcp_stub` for the full tool description with examples.
