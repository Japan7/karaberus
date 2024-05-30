project('karaberus', 'c', license: 'AGPL-3.0-and-later', version: '0.1.0')

go = find_program('go', required: true)
ffmpeg = find_program('ffmpeg', required: false)

go_files = files(
    'go.mod',
    'go.sum',
)

karaberus_server_path = join_paths(meson.current_source_dir(), 'server')

subdir('server')

karaberus_inputs = go_files + karaberus_server_files

go_env = environment()
go_env.set('CGO_ENABLED', '1')
go_env.set('CGO_CFLAGS', '-Wno-unused-parameter')

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
go_env.set('KARABERUS_DELETE_DB', '1')
go_env.set('KARABERUS_TEST_DIR', meson.current_build_dir())

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
    install: true,
    install_dir: 'bin',
    install_mode: 'rwxr-xr-x',
)

test_deps = []

if ffmpeg.found()
    create_test_mkv = custom_target(
        'create_test_mkv',
        output: 'karaberus_test.mkv',
        build_by_default: false,
        command: [
            ffmpeg,
            '-f', 'lavfi',
            '-i', 'testsrc=d=30',
            '-f', 'lavfi',
            '-i', 'sine',
            '-map', '0',
            '-map', '1',
            '-t', '30',
            '-c:v', 'libx264',
            '-preset', 'ultrafast',
            '-c:a', 'libopus',
            'karaberus_test.mkv',
        ],
    )

    test_deps += [create_test_mkv]
else
    message('ffmpeg executable was not found, S3 tests are disabled')
    go_env.set('SKIP_S3_TESTS', 'y')
endif

test(
    'karaberus_test',
    go,
    args: ['test', karaberus_server_path],
    depends: test_deps,
    env: go_env,
    timeout: 300,
)
