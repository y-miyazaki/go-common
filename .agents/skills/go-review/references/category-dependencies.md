# Dependencies (DEP)

**DEP-01 (SHOULD): List direct deps in go.mod with pinned versions**

Check: Are direct dependencies explicitly in go.mod, versions pinned, and regularly updated?
Why: Depending on indirect dependencies and unpinned versions cause unstable builds, unexpected behavior
Fix: Explicitly list direct dependencies in go.mod, pin versions, regular updates
