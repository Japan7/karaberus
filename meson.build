project(
    'karaberus',
    'c',
    license: 'AGPL-3.0-and-later',
    version: '0.1.0',
    default_options: [
        'buildtype=debugoptimized',
        'default_library=shared',
        'ffmpeg:default_library=shared',
        'libass:default_library=shared',
        'ffmpeg:zlib=enabled',
        'fribidi:bin=false',
        'c_std=gnu17',
        'warning_level=2',
    ],
)

go = find_program('go', required: true)
python = find_program('python3', required: true)
npm = find_program('npm', required: true)
meson_cmd = find_program('meson', required: true)
if not get_option('no_native_deps')
    ffmpeg = find_program('ffmpeg', required: get_option('s3_tests').enabled())
endif

c_compiler = meson.get_compiler('c')

# go build env
go_env = environment()
go_test_env = environment()
if get_option('no_native_deps')
    go_env.set('CGO_ENABLED', '0')
    go_test_env.set('NO_NATIVE_DEPS', '1')
else
    go_env.set('CGO_ENABLED', '1')
    if get_option('debug')
        go_env.set('CGO_CFLAGS', '-Wno-unused-parameter -Og')
    else
        go_env.set('CGO_CFLAGS', '-Wno-unused-parameter')
    endif
    go_env.set('CC', ' '.join(c_compiler.cmd_array()))

    pkg_config = find_program('pkg-config')
    go_env.set('PKG_CONFIG', pkg_config.full_path())

    pkg_config_path = get_option('pkg_config_path')
    pkg_config_path += meson.current_build_dir() / 'meson-uninstalled'
    go_env.set('PKG_CONFIG_PATH', ':'.join(pkg_config_path))
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

if get_option('builtin_oidc_env')
    go_env.set('KARABERUS_OIDC_ISSUER', 'http://localhost:9998')
    go_env.set('KARABERUS_OIDC_KEY_ID', 'karaberus')
    go_env.set('KARABERUS_OIDC_CLIENT_ID', 'web')
    go_env.set('KARABERUS_OIDC_CLIENT_SECRET', 'secret')
    go_env.set('KARABERUS_OIDC_GROUPS_CLAIM', 'groups')
    go_env.set('KARABERUS_OIDC_SCOPES', 'openid profile email groups')
    go_env.set('KARABERUS_OIDC_ADMIN_GROUP', 'admin')
    go_env.set('KARABERUS_OIDC_JWT_SIGN_KEY', 'sign_key')
endif

if get_option('builtin_s3_env')
    go_env.set('KARABERUS_S3_ENDPOINT', '127.0.0.1:9000')
endif

meson.add_devenv(go_env)

# go test build env
go_test_build_env = environment()
go_test_build_env.set('CGO_ENABLED', '0')

# go test env
go_test_env.set('KARABERUS_DB_FILE', 'file::memory:?cache=shared')
go_test_env.set('KARABERUS_TEST_DIR_GENERATED', meson.current_build_dir())
go_test_env.set('KARABERUS_TEST_DIR', meson.current_source_dir() / 'tests')
go_test_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)
go_test_env.set(
    'KARABERUS_S3_TEST_DB_FILE',
    meson.current_build_dir() / 's3_tests.db',
)

# go test binaries env
gobin_env = environment()
gobin_env.set('GOBIN', meson.current_build_dir())

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

go_tools_modfile = meson.current_source_dir() / 'tools' / 'go.mod'
go_tools_opts = [
    '-C',
    meson.current_source_dir() / 'tools',
    '-modfile',
    go_tools_modfile,
]

go_tools_files = files('tools' / 'go.sum') + go_tools_modfile

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
    karaberus_tools_deps += dependency(
        'appleframeworks',
        modules: ['foundation', 'security'],
        required: host_machine.system() == 'darwin',
    )

    go_files += files(
        'karaberus_tools' / 'karaberus_tools.c',
        'karaberus_tools' / 'karaberus_tools.h',
    )
endif

subdir('server')

karaberus_inputs = go_files + karaberus_server_files
karaberus_test_inputs = karaberus_inputs + karaberus_server_tests

if get_option('debug')
    go_build_options = ['-gcflags=all=-N -l']
else
    go_build_options = []
endif

if get_option('no_native_deps')
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
            '-C',
            meson.current_source_dir(),
            go_build_options,
            '-o',
            karaberus_output,
            meson.current_source_dir(),
        ],
        install: true,
        install_dir: 'bin',
        install_mode: 'rwxr-xr-x',
    )
    karaberus_go_build = karaberus
else
    karaberus_output = meson.current_build_dir() / 'karaberus.a'
    karaberus_archive = custom_target(
        'karaberus.a',
        env: go_env,
        command: [
            go,
            'build',
            '-C',
            meson.current_source_dir(),
            go_build_options,
            '-buildmode=c-archive',
            '-o',
            karaberus_output,
            meson.current_source_dir(),
        ],
        input: karaberus_inputs,
        output: 'karaberus.a',
    )
    karaberus_go_build = karaberus_archive

    # no-op but helps define the header for the executable below
    karaberus_header = custom_target(
        'karaberus.h',
        env: go_env,
        input: karaberus_archive,
        command: [go, 'version'],
        output: 'karaberus.h',
    )

    karaberus_main_c = files('karaberus' / 'karaberus.c')
    karaberus_a = meson.current_build_dir() / 'karaberus.a'

    karaberus = executable(
        'karaberus',
        karaberus_main_c,
        link_with: karaberus_archive,
        dependencies: karaberus_tools_deps,
        sources: karaberus_header,
        install: true,
    )
endif

go_test_env.set('KARABERUS_BIN', karaberus.full_path())

zitadel_karaberus = custom_target(
    'zitadel_karaberus',
    command: [
        go,
        'build',
        go_tools_opts,
        '-o',
        meson.current_build_dir() / 'zitadel-karaberus',
        'github.com/odrling/zitadel-karaberus',
    ],
    depend_files: go_tools_files,
    output: 'zitadel-karaberus',
)
go_test_env.set('OIDC_SERVER_EXE', zitadel_karaberus.full_path())

gofakes3 = custom_target(
    'gofakes3',
    command: [
        go,
        'build',
        go_tools_opts,
        '-o',
        meson.current_build_dir() / 'gofakes3',
        'github.com/johannesboyne/gofakes3/cmd/gofakes3',
    ],
    depend_files: go_tools_files,
    output: 'gofakes3',
)

go_test_env.set('GOFAKES3_EXE', gofakes3.full_path())

karaberus_test_output = meson.current_build_dir() / 'karaberus_test'
karaberus_server_path = meson.current_source_dir() / 'server'

karaberus_test = custom_target(
    'karaberus_test',
    env: go_test_build_env,
    input: karaberus_test_inputs,
    output: 'karaberus_test',
    command: [
        go,
        'test',
        '-C',
        meson.current_source_dir(),
        '-o',
        karaberus_test_output,
        '-c',
        karaberus_server_path,
    ],
    # so it can't build while karaberus is being built
    depends: karaberus_go_build,
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
            command: [
                go,
                'run',
                '-C',
                meson.current_source_dir(),
                '.',
                'openapi',
            ],
            output: 'openapi.yaml',
            capture: true,
            depend_files: karaberus_inputs,
            env: {'CGO_ENABLED': '0'},
        )
    endif
endif

subdir('ui')

run_target('run', command: [karaberus], depends: [karaberus_ui], env: go_env)
run_target('oidc', command: [zitadel_karaberus])
run_target(
    's3',
    command: [
        gofakes3,
        '-backend',
        'fs',
        '-initialbucket',
        'karaberus',
        '-fs.path',
        meson.current_build_dir() / 's3',
        '-fs.create',
    ],
)

test_deps = [zitadel_karaberus, gofakes3]

if get_option('s3_tests').enabled()

    if get_option('no_native_deps')
        warning(
            'S3 tests will run but checks won\'t be tested because of -Dno_native_deps=true.',
        )
        test_deps += custom_target(
            'karaberus_test.mkv',
            output: 'karaberus_test.mkv',
            build_by_default: false,
            command: [python, files('tests' / 'make_dummy_file.py')],
            capture: true,
        )

        test_deps += custom_target(
            'karaberus_test.opus',
            output: 'karaberus_test.opus',
            build_by_default: false,
            command: [python, files('tests' / 'make_dummy_file.py')],
            capture: true,
        )
    else
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
    endif

    test(
        's3_tests',
        python,
        args: [files('tests' / 'tests.py')],
        env: go_test_env,
    )
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
            depend_files: go_tools_files,
            output: 'staticcheck',
            env: gobin_env,
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
            depend_files: go_tools_files,
            output: 'errcheck',
            env: gobin_env,
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
