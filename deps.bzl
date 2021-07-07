"""External dependencies.

Macros defined in this file should be called near the end of the WORKSPACE
file, once the workspace dependencies have been set up.

The go_dependencies() macro is automatically generated by running:
$ bazel run //:gazelle -- update-repos -from_file=go.mod
"""

load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    """Go dependencies.

    DO NOT EDIT: This macro is auto-generated by running:
    $ bazel run //:gazelle -- update-repos -from_file=go.mod
    """
    go_repository(
        name = "com_github_andybalholm_cascadia",
        importpath = "github.com/andybalholm/cascadia",
        sum = "h1:BuuO6sSfQNFRu1LppgbD25Hr2vLYW25JvxHs5zzsLTo=",
        version = "v1.1.0",
    )

    go_repository(
        name = "com_github_antchfx_htmlquery",
        importpath = "github.com/antchfx/htmlquery",
        sum = "h1:sP3NFDneHx2stfNXCKbhHFo8XgNjCACnU/4AO5gWz6M=",
        version = "v1.2.3",
    )

    go_repository(
        name = "com_github_antchfx_xmlquery",
        importpath = "github.com/antchfx/xmlquery",
        sum = "h1:kaEVzH1mNo/2AJZrhZjAaAUTy2Nn2zxGfYYU8jWfXOo=",
        version = "v1.3.6",
    )

    go_repository(
        name = "com_github_antchfx_xpath",
        importpath = "github.com/antchfx/xpath",
        sum = "h1:cJ0pOvEdN/WvYXxvRrzQH9x5QWKpzHacYO8qzCcDYAg=",
        version = "v1.1.10",
    )

    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:ZDRjVQ15GmhC3fiQ8ni8+OwkZQO4DARzQgrnXU1Liz8=",
        version = "v1.1.0",
    )

    go_repository(
        name = "com_github_gobwas_glob",
        importpath = "github.com/gobwas/glob",
        sum = "h1:A4xDbljILXROh+kObIiy5kIaPYD8e96x1tgBhUI5J+Y=",
        version = "v0.2.3",
    )

    go_repository(
        name = "com_github_gocolly_colly",
        importpath = "github.com/gocolly/colly",
        sum = "h1:qRz9YAn8FIH0qzgNUw+HT9UN7wm1oF9OBAilwEWpyrI=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:1r7pUrabqp18hOBcwBwiTsbnFeTZHV9eER/QT5JVZxY=",
        version = "v0.0.0-20200121045136-8c9f03a8e57e",
    )

    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:YF8+flBXS5eO826T4nzqPrxfhQThhXl0YzfuUPu4SBg=",
        version = "v1.3.1",
    )

    go_repository(
        name = "com_github_kennygrant_sanitize",
        importpath = "github.com/kennygrant/sanitize",
        sum = "h1:gN25/otpP5vAsO2djbMhF/LQX6R7+O1TB4yv8NzpJ3o=",
        version = "v1.2.4",
    )

    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_puerkitobio_goquery",
        importpath = "github.com/PuerkitoBio/goquery",
        sum = "h1:O5SP3b9JWqMSVMG69zMfj577zwkSNpxrFf7ybS74eiw=",
        version = "v1.7.0",
    )

    go_repository(
        name = "com_github_saintfish_chardet",
        importpath = "github.com/saintfish/chardet",
        sum = "h1:NugYot0LIVPxTvN8n+Kvkn6TrbMyxQiuvKdEwFdR9vI=",
        version = "v0.0.0-20120816061221-3af4cd4741ca",
    )

    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
        version = "v0.1.0",
    )

    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:TivCn/peBQ7UY8ooIcPgZFpTNSz0Q2U6UrFlUfqbe0Q=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_github_temoto_robotstxt",
        importpath = "github.com/temoto/robotstxt",
        sum = "h1:W2pOjSJ6SWvldyEuiFXNxz3xZ8aiWX5LbfDiOFd7Fxg=",
        version = "v1.1.2",
    )

    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
        version = "v1.6.7",
    )

    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:psW17arqaxU48Z5kZ0CQnkZWQJsqcURM6tKiBApRjXI=",
        version = "v0.0.0-20200622213623-75b288015ac9",
    )

    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:XpT3nA5TvE525Ne3hInMh6+GETgn27Zfm9dxsThnX2Q=",
        version = "v0.0.0-20210614182718-04defd469f4e",
    )

    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:b3NXsE2LusjYGGjL5bxEVZZORm/YEFFrWFjR8eFrw/c=",
        version = "v0.0.0-20210423082822-04245dca01da",
    )

    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:v+OssWQX+hTHEmOBgwxdZxK4zHq3yOs8F9J7mk0PY8E=",
        version = "v0.0.0-20201126162022-7de9c90e9dd1",
    )

    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:aRYxNxv6iGQlyVaZmk6ZgYEDa+Jg18DxebPSrd6bg1M=",
        version = "v0.3.6",
    )

    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:FDhOuMEY4JVRztM/gsbk+IKUQ8kj74bxZrgw87eMMVc=",
        version = "v0.0.0-20180917221912-90fa682c2a6e",
    )
