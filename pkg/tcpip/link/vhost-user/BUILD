load("//tools/go_stateify:defs.bzl", "go_library")

package(licenses = ["notice"])

go_library(
    name = "vhost-user",
    srcs = ["vhostuser.go"],
    importpath = "gvisor.googlesource.com/gvisor/pkg/tcpip/link/vhost-user",
    visibility = ["//:sandbox"],
    deps = [
        "//pkg/tcpip",
        "//pkg/tcpip/buffer",
        "//pkg/tcpip/stack",
    ],
)
