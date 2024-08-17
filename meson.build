project(
    'karaberus',
    'c',
    license: 'AGPL-3.0-and-later',
    version: '0.1.0',
    default_options: [
        'default_library=shared',
        'ffmpeg:default_library=shared',
        'libass:default_library=shared',
        'fribidi:bin=false',
        'c_std=gnu17',
        'warning_level=2',
    ],
)

go = find_program('go', required: true)
ffmpeg = find_program('ffmpeg', required: get_option('s3_tests').enabled())
python = find_program('python3', required: true)
npm = find_program('npm', required: true)
meson_cmd = find_program('meson', required: true)

c_compiler = meson.get_compiler('c')

# go build env
go_env = environment()
if get_option('no_native_deps')
    go_env.set('CGO_ENABLED', '0')
else
    go_env.set('CGO_ENABLED', '1')
    go_env.set('CGO_CFLAGS', '-Wno-unused-parameter')
    go_env.set('CC', ' '.join(c_compiler.cmd_array()))
    pkg_config_path = meson.current_build_dir() / 'meson-uninstalled'
    go_env.set('PKG_CONFIG_PATH', pkg_config_path)
endif

go_os = get_option('go_os')
go_arch = get_option('go_arch')
if go_os != ''
    go_env.set('GOOS', go_os)
endif
if go_arch != ''
    go_env.set('GOARCH', go_arch)
endif

db_file = meson.current_build_dir() / 'karaberus.db'
karaberus_ui_dist_dir = meson.current_build_dir() / 'ui' / 'ui_dist'
go_env.set('KARABERUS_DB_FILE', db_file)
go_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)

meson.add_devenv(go_env)

# go test build env
go_test_build_env = environment()
go_test_build_env.set('CGO_ENABLED', '0')

# go test env
go_test_env = environment()
go_test_env.set('KARABERUS_DB_FILE', 'file::memory:?cache=shared')
go_test_env.set('KARABERUS_TEST_DIR_GENERATED', meson.current_build_dir())
go_test_env.set('KARABERUS_TEST_DIR', meson.current_source_dir() / 'tests')
go_test_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)

ui_source = meson.current_source_dir() / 'ui'

# npm wrapper env
npm_env = environment()
npm_env.set('NPM', npm.full_path())
npm_env.set('SOURCE', ui_source)


go_files = files('go.mod', 'go.sum', 'main.go', 'main_cgo.go')
go_files += files(
    'karaberus_tools' / 'cbinds.go',
    'karaberus_tools' / 'nocbinds.go',
    'karaberus_tools' / 'model.go',
)

go_tools_modfile = files('tools' / 'go.mod')
go_tools_opts = ['-modfile', go_tools_modfile]

karaberus_tools_deps = []

if not get_option('no_native_deps')
    dakara_check = dependency(
        'dakara_check',
        version: '>4.0.4',
        required: true,
        default_options: ['programs=false'],
    )
    karaberus_tools_deps += dakara_check
    karaberus_tools_deps += dependency('libavutil', required: true)
    karaberus_tools_deps += dependency('libavformat', required: true)

    go_files += files(
        'karaberus_tools' / 'karaberus_tools.c',
        'karaberus_tools' / 'karaberus_tools.h',
    )
endif

subdir('server')

karaberus_inputs = go_files + karaberus_server_files
karaberus_test_inputs = karaberus_inputs + karaberus_server_tests

if get_option('no_native_deps')
    karaberus = custom_target(
        'karaberus',
        build_by_default: true,
        env: go_env,
        input: karaberus_inputs,
        output: 'karaberus',
        command: [go, 'build', '-o', '@OUTPUT@', meson.current_source_dir()],
        install: true,
        install_dir: 'bin',
        install_mode: 'rwxr-xr-x',
    )
else
    karaberus_archive = custom_target(
        'karaberus.a',
        env: go_env,
        command: [
            go,
            'build',
            '-buildmode=c-archive',
            '-o',
            '@OUTPUT@',
            meson.current_source_dir(),
        ],
        output: 'karaberus.a',
    )

    # no-op but helps define the header for the executable below
    karaberus_header = custom_target(
        'karaberus.h',
        env: go_env,
        input: karaberus_archive,
        command: [go, 'version'],
        output: 'karaberus.h',
    )

    karaberus_main_c = files('karaberus' / 'karaberus.c')

    karaberus = executable(
        'karaberus',
        karaberus_main_c,
        link_with: karaberus_archive,
        dependencies: karaberus_tools_deps,
        sources: karaberus_header,
        install: true,
    )
endif

karaberus_server_path = meson.current_source_dir() / 'server'

karaberus_test = custom_target(
    'karaberus_test',
    env: go_test_build_env,
    input: karaberus_test_inputs,
    output: 'karaberus_test',
    command: [go, 'test', '-o', '@OUTPUT@', '-c', karaberus_server_path],
)

if not meson.is_cross_build()
    openapi_spec = custom_target(
        'openapi_spec',
        command: [karaberus, 'openapi'],
        output: 'openapi.yaml',
        capture: true,
        depend_files: karaberus_inputs,
    )
else
    if meson.can_run_host_binaries()
        openapi_spec = custom_target(
            'openapi_spec',
            command: [karaberus, 'openapi'],
            output: 'openapi.yaml',
            capture: true,
            depend_files: karaberus_inputs,
        )
    else
        openapi_spec = custom_target(
            'openapi_spec',
            command: [go, 'run', meson.current_source_dir(), 'openapi'],
            output: 'openapi.yaml',
            capture: true,
            depend_files: karaberus_inputs,
            env: {'CGO_ENABLED': '0'},
        )
    endif
endif

subdir('ui')

run_target('run', command: [karaberus], depends: [karaberus_ui], env: go_env)

test_deps = []

if get_option('s3_tests').enabled()
    test_deps += custom_target(
        'karaberus_test.mkv',
        output: 'karaberus_test.mkv',
        build_by_default: false,
        command: [
            ffmpeg,
            '-f',
            'lavfi',
            '-i',
            'testsrc=d=30',
            '-f',
            'lavfi',
            '-i',
            'sine',
            '-map',
            '0',
            '-map',
            '1',
            '-t',
            '30',
            '-c:v',
            'libx264',
            '-preset',
            'ultrafast',
            '-c:a',
            'libopus',
            '-y',
            '@OUTPUT@',
        ],
    )

    test_deps += custom_target(
        'karaberus_test.opus',
        output: 'karaberus_test.opus',
        build_by_default: false,
        command: [
            ffmpeg,
            '-f',
            'lavfi',
            '-i',
            'sine',
            '-map',
            '0',
            '-t',
            '30',
            '-c:a',
            'libopus',
            '-y',
            '@OUTPUT@',
        ],
    )
else
    go_test_env.set('SKIP_S3_TESTS', 'y')
endif

if get_option('test')
    test(
        'karaberus_run_test',
        karaberus_test,
        depends: test_deps,
        env: go_test_env,
        timeout: 300,
    )

    if get_option('staticcheck')
        staticcheck = custom_target(
            'staticcheck',
            command: [
                go,
                'install',
                go_tools_opts,
                'honnef.co/go/tools/cmd/staticcheck',
            ],
            output: 'staticcheck',
            env: {'GOBIN': meson.current_build_dir()},
        )
        test(
            'staticcheck',
            staticcheck,
            args: ['./...'],
            depends: karaberus,
            env: go_env,
            workdir: meson.current_source_dir(),
            timeout: 300,
        )
    endif

    if get_option('errcheck')
        errcheck = custom_target(
            'errcheck',
            command: [
                go,
                'install',
                go_tools_opts,
                'github.com/kisielk/errcheck',
            ],
            output: 'errcheck',
            env: {'GOBIN': meson.current_build_dir()},
        )
        test(
            'errcheck',
            errcheck,
            args: ['./...'],
            depends: karaberus,
            env: go_env,
            workdir: meson.current_source_dir(),
            timeout: 300,
        )
    endif
endif

test(
    'meson_format_check',
    meson_cmd,
    args: ['format', '-r', '-q', meson.current_source_dir()],
)
