# Root module security and logger qualification

The root github.com/Laisky/gorm module selects golang.org/x/crypto v0.57.0 and therefore requires Go 1.26.0; the patched toolchain is pinned to Go 1.26.9. The nested v2 module is separate and unchanged.

A broken String method on a named string previously replaced SQL log values with fmt's panic diagnostic. The logger now preserves the underlying string on that failure while retaining successful Stringer behavior and log context. Tests assert SQL content, source, elapsed time and row count, with ordinary Stringer/string/number/nil controls; none of the database assertions are skipped or weakened. Dynamic formatted assertions now use an explicit %s format for newer vet compatibility.

Qualification: the strengthened panicking-Stringer case failed before the source fix; full race suites pass for the declared sqlite dialect and an isolated PostgreSQL 17 fixture, including PostgreSQL read-only transaction and composite-key behavior. Module verification, vet and build pass. Source govulncheck finds zero called or imported vulnerabilities and one unused required-module advisory. MySQL/MSSQL and the nested v2 module were not exercised. No shared database or deployment was used; the synthetic PostgreSQL fixture was removed.
