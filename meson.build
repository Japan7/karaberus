project(
    'karaberus',
    'c',
    license: 'AGPL-3.0-and-later',
    version: '0.1.0',
    default_options: [
        'ffmpegaacsucks:prefer_static=true',
        'ffmpegaacsucks:default_library=static',
    ],
)

go = find_program('go', required: true)
ffmpeg = find_program('ffmpeg', required: false)
python = find_program('python3')
npm = find_program('npm')

go_files = files('go.mod', 'go.sum', 'main.go', 'karaberus_tools' / 'cbinds.go')

karaberus_server_path = join_paths(meson.current_source_dir(), 'server')

subdir('server')

go_os = get_option('go_os')
go_arch = get_option('go_arch')

c_compiler = meson.get_compiler('c')
cc_cmd = ' '.join(c_compiler.cmd_array())

db_file = meson.current_build_dir() / 'karaberus.db'
test_db_file = meson.current_build_dir() / 'karaberus_test.db'
karaberus_ui_dist_dir = meson.current_build_dir() / 'ui' / 'ui_dist'
pkg_config_path = meson.current_build_dir() / 'meson-uninstalled'

dakara_check = dependency(
    'dakara_check',
    version: '>=4.0.1',
    required: true,
    default_options: ['default_library=static', 'prefer_static=true'],
)

dakara_check_objects = []

if dakara_check.type_name() == 'internal'
    dakara_check_sub = subproject('dakara_check')
    libdakara_check = dakara_check_sub.get_variable('libdakara_check')
    dakara_check_objects = libdakara_check.extract_all_objects(recursive: true)
endif

karaberus_tools_c = files('karaberus_tools' / 'karaberus_tools.c')
karaberus_tools_h = files('karaberus_tools' / 'karaberus_tools.h')

inc_karaberus_tools = include_directories('karaberus_tools')

libkaraberus_tools = static_library(
    'karaberus_tools',
    karaberus_tools_c,
    karaberus_tools_h,
    dependencies: dakara_check,
    objects: dakara_check_objects,
    include_directories: inc_karaberus_tools,
    install: true,
)

pkg = import('pkgconfig')
pkg.generate(libkaraberus_tools)

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

go_env.set('KARABERUS_DB_FILE', db_file)
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
go_test_env.set('KARABERUS_DB_DELETE', '1')
go_test_env.set('KARABERUS_GENERATED_TEST_DIR', meson.current_build_dir())
go_test_env.set('KARABERUS_TEST_DIR', meson.current_source_dir() / 'tests')
go_test_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)

ui_source = meson.current_source_dir() / 'ui'

npm_env = environment()
npm_env.set('NPM', npm.full_path())
npm_env.set('SOURCE', ui_source)

karaberus_inputs = go_files + karaberus_server_files
karaberus_output = meson.current_build_dir() / 'karaberus'

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

karaberus_test_output = meson.current_build_dir() / 'karaberus_test'

karaberus_test = custom_target(
    'karaberus_test',
    env: go_env,
    input: karaberus_inputs,
    output: 'karaberus_test',
    command: [
        go,
        'test',
        '-C', meson.current_source_dir(),
        '-buildmode=pie',
        '-trimpath',
        '-ldflags',
        '-linkmode=external -s',
        '-o', karaberus_test_output,
        '-c', karaberus_server_path,
    ],
    depends: libkaraberus_tools,
)

openapi_spec = custom_target(
    'openapi_spec',
    command: [karaberus, 'openapi'],
    output: 'openapi.yaml',
    capture: true,
    depend_files: karaberus_inputs,
    depends: libkaraberus_tools,
    env: go_env,
)

subdir('ui')

run_target(
    'run',
    command: [karaberus],
    depends: [libkaraberus_tools, karaberus_ui],
    env: go_env,
)

run_target(
    'create-token',
    command: [
        go,
        'run',
        '-C', meson.current_source_dir(),
        meson.current_source_dir(),
        'create-token',
    ],
    depends: libkaraberus_tools,
    env: go_env,
)

test_deps = [libkaraberus_tools]

if ffmpeg.found()
    mkv_test_file = 'karaberus_test.mkv'
    test_deps += custom_target(
        'create_test_mkv',
        output: mkv_test_file,
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
            '-y', mkv_test_file,
        ],
    )

    opus_test_file = 'karaberus_test.opus'
    test_deps += custom_target(
        'create_test_opus',
        output: opus_test_file,
        build_by_default: false,
        command: [
            ffmpeg,
            '-f', 'lavfi',
            '-i', 'sine',
            '-map', '0',
            '-t', '30',
            '-c:a', 'libopus',
            '-y', opus_test_file,
        ],
    )
else
    warning('ffmpeg executable was not found, S3 tests are disabled')
    go_test_env.set('SKIP_S3_TESTS', 'y')
endif

test(
    'karaberus_run_test',
    karaberus_test,
    depends: test_deps,
    env: go_test_env,
    timeout: 300,
)
