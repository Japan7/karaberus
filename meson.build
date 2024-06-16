project('karaberus', 'c', license: 'AGPL-3.0-and-later', version: '0.1.0')

go = find_program('go', required: true)
ffmpeg = find_program('ffmpeg', required: false)

go_files = files(
    'go.mod',
    'go.sum',
)

karaberus_server_path = join_paths(meson.current_source_dir(), 'server')

subdir('karaberus_tools')
install_headers(karaberus_tools_h)

libkaraberus_tools = library(
    'karaberus_tools',
    karaberus_tools,
    karaberus_tools_h,
    dependencies: dakara_check,
    install: true,
)

subdir('server')
karaberus_inputs = go_files + karaberus_server_files

go_os = get_option('go_os')
go_arch = get_option('go_arch')

c_compiler = meson.get_compiler('c')
cc_cmd = ' '.join(c_compiler.cmd_array())

test_db_file = meson.current_build_dir() / 'karaberus.db'
karaberus_ui_dist_dir = meson.current_build_dir() / 'ui_dist'
pkg_config_path = join_paths(meson.current_build_dir(), 'meson-uninstalled')

# go build env
go_env = environment()
go_env.set('CGO_ENABLED', '1')
go_env.set('CGO_CFLAGS', '-Wno-unused-parameter')
go_env.set('CC', cc_cmd)

go_env.set('PKG_CONFIG_PATH', pkg_config_path)

if go_os != ''
    go_env.set('GOOS', go_os)
endif
if go_arch != ''
    go_env.set('GOARCH', go_arch)
endif

go_env.set('KARABERUS_DB_FILE', test_db_file)
go_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)

# go test env
go_test_env = environment()
go_test_env.set('CGO_ENABLED', '1')
go_test_env.set('CGO_CFLAGS', '-Wno-unused-parameter')
go_test_env.set('CC', cc_cmd)

go_test_env.set('PKG_CONFIG_PATH', pkg_config_path)

if go_os != ''
    go_test_env.set('GOOS', go_os)
endif
if go_arch != ''
    go_test_env.set('GOARCH', go_arch)
endif

go_test_env.set('KARABERUS_DB_FILE', test_db_file)
go_test_env.set('KARABERUS_DELETE_DB', '1')
go_test_env.set('KARABERUS_TEST_DIR', meson.current_build_dir())
go_test_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)

npm = find_program('npm')

ui_source = meson.current_source_dir() / 'ui'

npm_env = environment()
npm_env.set('NPM', npm.full_path())
npm_env.set('SOURCE', ui_source)

python = find_program('python3')

subdir('ui')

karaberus_ui = custom_target(
    'karaberus_ui',
    build_by_default: true,
    env: npm_env,
    input: ui_src,
    output: 'ui_dist',
    command: [
        python,
        meson.current_source_dir() / 'scripts' / 'npm_run.py',
        'build',
        meson.current_build_dir() / 'ui_dist',
    ],
    install: true,
    install_dir: 'share' / 'karaberus',
)

karaberus_ui_run = run_target(
    'ui_run',
    env: npm_env,
    command: [
        python,
        meson.current_source_dir() / 'scripts' / 'npm_run.py',
        'dev',
    ],
)

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
    depends: [libkaraberus_tools, karaberus_ui],
    install: true,
    install_dir: 'bin',
    install_mode: 'rwxr-xr-x',
)

karaberus_run = run_target(
    'run',
    command: [
        go,
        'run',
        '-C', meson.current_source_dir(),
        meson.current_source_dir(),
    ],
    depends: [libkaraberus_tools, karaberus_ui],
    env: go_env,
)

test_deps = [libkaraberus_tools]

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
    warning('ffmpeg executable was not found, S3 tests are disabled')
    go_test_env.set('SKIP_S3_TESTS', 'y')
endif

test(
    'karaberus_test',
    go,
    args: ['test', karaberus_server_path],
    depends: test_deps,
    env: go_test_env,
    timeout: 300,
)
