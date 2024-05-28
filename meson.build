project('karaberus', 'c', license: 'AGPL-3.0-and-later', version: '0.1.0')

go = find_program('go', required: true)

go_files = files(
    'go.mod',
    'go.sum',
)

karaberus_tools_path = join_paths(meson.current_source_dir(), 'karaberus_tools')
karaberus_server_path = join_paths(meson.current_source_dir(), 'server')

subdir('karaberus_tools')
subdir('server')

karaberus_inputs = go_files + karaberus_server_files

go_env = environment()
go_env.set('CGO_ENABLED', '1')

pkg_config_path = join_paths(meson.current_build_dir(), 'meson-uninstalled')
go_env.set('PKG_CONFIG_PATH', pkg_config_path)

go_os = get_option('go_os')
go_arch = get_option('go_arch')
if go_os != ''
    go_env.set('GOOS', go_os)
endif
if go_arch != ''
    go_env.set('GOARCH', go_arch)
endif

c_compiler = meson.get_compiler('c')
go_env.set('CC', ' '.join(c_compiler.cmd_array()))

test_db_file = meson.current_build_dir() / 'karaberus.db'
go_env.set('KARABERUS_DB_FILE', test_db_file)

install_headers(karaberus_tools_h)

libkaraberus_tools = library(
    'karaberus_tools',
    karaberus_tools,
    karaberus_tools_h,
    dependencies: dakara_check,
    install: true,
)

pkg = import('pkgconfig')
pkg.generate(libkaraberus_tools)

karaberus_output = join_paths(meson.current_build_dir(), 'karaberus')

karaberus = custom_target(
    'karaberus',
    build_by_default: true,
    env: go_env,
    input: karaberus_inputs,
    output: 'karaberus',
    command: [
        go,
        'build',
        '-C', meson.current_source_dir(),
        '-buildmode=pie',
        '-trimpath',
        '-ldflags',
        '-linkmode=external -s',
        '-o', karaberus_output,
        meson.current_source_dir(),
    ],
    depends: libkaraberus_tools,
    install: true,
    install_dir: 'bin',
    install_mode: 'rwxr-xr-x',
)

test(
    'karaberus_test',
    go,
    args: ['test', karaberus_server_path],
    depends: libkaraberus_tools,
    env: go_env,
)
