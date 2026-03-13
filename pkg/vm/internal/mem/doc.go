// Package mem provides VM storage helpers and the narrow ownership layer used
// to clean up direct register-held closers.
//
// Ownership policy:
//   - Produced values marked by the VM as frame-owned, for example via
//     writeProducedRegister, become owned by the active frame.
//   - Borrowed values do not become owned unless they are explicitly
//     transferred.
//   - Moving a value into a container or other external sink forfeits frame
//     ownership.
//   - Returns and tail-calls transfer ownership of surviving direct register
//     values.
//   - Register, window, and scratch storage only scrub slots to runtime.None;
//     they never close values directly.
//   - Automatic cleanup is intentionally limited to direct register-held
//     io.Closer values tracked by OwnedResources. It does not manage deep
//     container contents or act as a universal runtime lifetime manager.
package mem
