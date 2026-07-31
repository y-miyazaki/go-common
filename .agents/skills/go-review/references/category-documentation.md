# Documentation (DOC)

**DOC-01 (SHOULD): Package has a doc comment stating purpose**

Check: Are package doc comments, package purpose, and usage documented?
Why: Missing package doc comments make API understanding difficult, increase misuse, delay onboarding
Fix: Add package doc comments, document purpose, responsibility, usage examples

**DOC-02 (MUST): Public APIs have godoc covering args/returns/errors**

Check: Are all public APIs documented with godoc, arguments, return values, and error conditions specified?
Why: Missing or insufficient public function comments make API usage unclear, cause misuse
Fix: Document all public APIs with godoc, specify arguments, return values, error conditions

**DOC-03 (SHOULD): Comments are consistently English**

Check: Are comments unified in English, grammar-checked, and concise?
Why: Mixed languages and grammar errors reduce readability, make internationalization difficult
Fix: Unify in English, check grammar, write concisely and clearly
