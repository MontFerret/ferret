// Package module defines the extension contracts used to customize engine
// bootstrap and lifecycle behavior.
//
// Modules register host-scoped services, such as runtime library entries and
// encoding codecs, and can attach hooks that run during engine, plan, and
// session stages.
package module
