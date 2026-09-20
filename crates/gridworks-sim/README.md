# `gridworks-sim`

This crate is the sole canonical home of GRIDWORKS deterministic simulation mathematics. WP-002 provides versioned state/command contracts, explicit external-time advancement, a pinned seeded RNG, SHA-256 replay digests, typed errors and a generic proof domain.

The kernel has no network, database, Godot, filesystem, process-environment or device-clock dependency. Authoritative timestamps and rules/content versions are inputs. Facility/component mechanics, GDExtension integration and server validation remain separately authorized work-package boundaries.
