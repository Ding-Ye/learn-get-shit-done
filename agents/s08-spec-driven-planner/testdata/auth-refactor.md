# Auth refactor — Specification

## Goal

Replace cookie-based session auth with stateless JWTs, while preserving
existing API surfaces and rotating secrets without downtime.

## Phase P01: Token issuance

Goal: Mint signed JWTs on successful login.

Tasks:
- add POST /token endpoint
- sign tokens with HS256 + 12h expiry
- thread JWT_SECRET through config

Verify: hitting POST /token with valid credentials returns a parseable JWT
whose `exp` is between 11h59m and 12h01m in the future.

## Phase P02: Token validation

Goal: Validate tokens in middleware and surface 401 on failure.

Tasks:
- write middleware/auth that calls jwt.Parse
- map invalid-signature to 401 with descriptive body
- add unit tests for expired + tampered tokens

Verify: integration test sends a tampered JWT and asserts the response is 401.
