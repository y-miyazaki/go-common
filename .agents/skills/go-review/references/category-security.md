## Security (SEC)

**SEC-01 (SHOULD): Validate inputs; ban string-concat SQL**

Check: Are input validation, prepared statements, and sanitization implemented?
Why: Unvalidated input and SQL string concatenation enable SQL injection, XSS attacks, data tampering
Fix: Mandatory prepared statements, implement validation, implement sanitization

**SEC-02 (SHOULD): Escape/sanitize outputs for HTML/JSON/CRLF sinks**

Check: Are HTML escaping, JSON injection prevention, and CRLF injection prevention present?
Why: Missing escaping causes XSS vulnerabilities, response tampering, session hijacking
Fix: Use html/template, context-appropriate escaping for output

**SEC-03 (SHOULD): Authenticate endpoints; verify JWT; enforce RBAC**

Check: Are all endpoints authenticated, JWT signature verified, and RBAC implemented?
Why: Skipped authentication and insufficient verification enable unauthorized access, privilege escalation, data leakage
Fix: Mandatory authentication for all endpoints, JWT signature verification, RBAC implementation

**SEC-04 (SHOULD): Mask passwords/tokens in logs**

Check: Are sensitive information masking functions and password/token masking present?
Why: Logging passwords and tokens causes credential leakage, GDPR violations
Fix: Implement sensitive information masking functions, structured logging, log rotation

**SEC-05 (SHOULD): Least privilege; no production debug; explicit CORS**

Check: Are least privilege principle, production debug disabled, and explicit CORS settings present?
Why: Insecure defaults cause security breaches, increased attack success rate
Fix: Least privilege principle, disable production debug, explicit CORS settings
