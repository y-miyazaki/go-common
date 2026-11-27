# Troubleshooting

- If evidence is partial, mark affected checks as deferred with explicit reason.
- If PR context is incomplete, list missing evidence and mark impacted checks as deferred.
- If PR contains only generated or vendor files, set `Status: skipped` in `## Checks Summary` with reason `no reviewable Go source`.
