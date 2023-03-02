package handlers

var module = `
package sui.alert

import future.keywords.if
import future.keywords.in

default allow := false

allow if {
    input.packageId == "GET1"
    input.path == ["salary", input.subject.user]
}

allow if is_admin

is_admin if "admin" in input.subject.groups
`
