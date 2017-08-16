pkg_name=foulkon
pkg_description="foulkon"
pkg_origin=chef
pkg_version="0.0.1"
pkg_maintainer="Chef Software Inc. <support@chef.io>"
pkg_license=('Apache-2.0')
pkg_bin_dirs=(bin)
pkg_build_deps=(
  # core/which # let's just ignore those errors. works fine without.
)
pkg_deps=(
  core/postgresql # for psql in hooks/init
  core/curl # for bootstrapping in hooks/init
)
pkg_scaffolding=afiune/scaffolding-go
scaffolding_go_base_path=github.com/Tecsisa
scaffolding_go_build_deps=(
  # github.com/Masterminds/glide # note: let's use the version foulkon uses
                                 # instead of master (this way)
)

do_prepare() {
  set -e

  build_line "mkdir -p $GOPATH/bin; export PATH=$GOPATH/bin:$PATH"
  mkdir -p $GOPATH/bin
  export PATH=$GOPATH/bin:$PATH
}

do_build() {
  pushd $scaffolding_go_pkg_path >/dev/null

  build_line "make deps generate"
  make deps generate

  # Note: We don't do 'make bin', because it's only these two we need
  #       (It's not worth installing env, and fixing up paths etc...)
  build_line "CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/{worker,proxy}"
  CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/worker
  CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/proxy
  popd
}

do_install() {
  build_line "copying worker and proxy binary"
  cp "${scaffolding_go_gopath:?}/bin/worker" $pkg_prefix/bin
  cp "${scaffolding_go_gopath:?}/bin/proxy" $pkg_prefix/bin
}
