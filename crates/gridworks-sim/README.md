# `gridworks-sim`

This crate is the sole future home of canonical deterministic simulation mathematics. WP-001 provides only the versioned command/snapshot boundary and deterministic smoke tests.

The boundary deliberately has no network, database, Godot, device-clock or service dependency. Authoritative timestamps and rules/content versions are inputs. Gameplay state transitions, seeded randomness, replay compatibility and GDExtension integration belong to their separately authorized work packages.
