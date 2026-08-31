# Release policy

The first release is created only from a successful merged main commit. The
workflow enables the repository immutable-releases setting with the GitHub
REST API, verifies the setting, creates one annotated tag, publishes exact
source/contract/lock/version/manifest/checksum assets, and audits the REST
release object, tag object, target commit, asset IDs, sizes, and SHA-256
digests. Existing tags and releases are never reused, deleted, or rewritten.

If an early release is found to be non-durable, it remains preserved as a
historical release and the next unused patch version is issued; the earlier
tag and release are not repaired in place.
