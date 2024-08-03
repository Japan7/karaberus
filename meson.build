project(
    'karaberus',
    'c',
    license: 'AGPL-3.0-and-later',
    version: '0.1.0',
    default_options: [
        'ffmpegaacsucks:prefer_static=true',
        'ffmpegaacsucks:default_library=static',
        'c_std=gnu17',
        'warning_level=2',
    ],
)

go = find_program('go', required: true)
ffmpeg = find_program('ffmpeg', required: false)
python = find_program('python3', required: true)
npm = find_program('npm', required: true)
qemu_user = find_program('qemu-' + host_machine.cpu_family(), required: false)
meson_cmd = find_program('meson', required: true)

go_files = files('go.mod', 'go.sum', 'main.go')
go_files += files(
    'karaberus_tools' / 'cbinds.go',
    'karaberus_tools' / 'nocbinds.go',
    'karaberus_tools' / 'model.go',
)

karaberus_server_path = join_paths(meson.current_source_dir(), 'server')

go_os = get_option('go_os')
go_arch = get_option('go_arch')

c_compiler = meson.get_compiler('c')
cc_cmd = ' '.join(c_compiler.cmd_array())

db_file = meson.current_build_dir() / 'karaberus.db'
karaberus_ui_dist_dir = meson.current_build_dir() / 'ui' / 'ui_dist'
pkg_config_path = meson.current_build_dir() / 'meson-uninstalled'

if not get_option('no_native_deps')
    karaberus_tools_deps = []

    dakara_check = dependency(
        'dakara_check',
        version: '>4.0.4',
        required: true,
        default_options: ['default_library=static', 'prefer_static=true'],
    )
    karaberus_tools_deps += dakara_check
    karaberus_tools_deps += dependency('libavutil', required: true)
    karaberus_tools_deps += dependency('libavformat', required: true)

    # force dakara_check deps to be linked whole too
    if dakara_check.type_name() == 'internal'
        dakara_check_sub = subproject('dakara_check')
        karaberus_tools_deps += dakara_check_sub.get_variable('deps')
    endif

    karaberus_tools_deps_fr = []

    foreach d : karaberus_tools_deps
        if d.type_name() == 'internal'
            karaberus_tools_deps_fr += d.as_link_whole()
        else
            karaberus_tools_deps_fr += d
        endif
    endforeach

    karaberus_tools_c = files('karaberus_tools' / 'karaberus_tools.c')
    karaberus_tools_h = files('karaberus_tools' / 'karaberus_tools.h')

    inc_karaberus_tools = include_directories('karaberus_tools')

    libkaraberus_tools = static_library(
        'karaberus_tools',
        karaberus_tools_c,
        karaberus_tools_h,
        dependencies: karaberus_tools_deps_fr,
        include_directories: inc_karaberus_tools,
        install: false,
    )

    pkg = import('pkgconfig')
    pkg.generate(libkaraberus_tools)
endif

# go build env
go_env = environment()
if get_option('no_native_deps')
    go_env.set('CGO_ENABLED', '0')
else
    go_env.set('CGO_ENABLED', '1')
endif
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

meson.add_devenv(go_env)

# go test env
go_test_env = environment()
go_test_env.set('KARABERUS_DB_FILE', 'file::memory:?cache=shared')
go_test_env.set('KARABERUS_TEST_DIR_GENERATED', meson.current_build_dir())
go_test_env.set('KARABERUS_TEST_DIR', meson.current_source_dir() / 'tests')
go_test_env.set('KARABERUS_UI_DIST_DIR', karaberus_ui_dist_dir)

if get_option('no_native_deps')
    go_build_flags = ['-buildmode=pie', '-trimpath']
else
    go_build_flags = [
        '-buildmode=pie',
        '-trimpath',
        '-ldflags',
        '-linkmode=external -s',
    ]
endif

ui_source = meson.current_source_dir() / 'ui'

npm_env = environment()
npm_env.set('NPM', npm.full_path())
npm_env.set('SOURCE', ui_source)

subdir('server')

karaberus_inputs = go_files + karaberus_server_files
karaberus_test_inputs = karaberus_inputs + karaberus_server_tests
karaberus_output = meson.current_build_dir() / 'karaberus'

karaberus_deps = []
if not get_option('no_native_deps')
    karaberus_deps += libkaraberus_tools
endif

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
        go_build_flags,
        '-o',
        karaberus_output,
        meson.current_source_dir(),
    ],
    depends: karaberus_deps,
    install: true,
    install_dir: 'bin',
    install_mode: 'rwxr-xr-x',
)

karaberus_test_output = meson.current_build_dir() / 'karaberus_test'

karaberus_test = custom_target(
    'karaberus_test',
    env: go_env,
    input: karaberus_test_inputs,
    output: 'karaberus_test',
    command: [
        go,
        'test',
        '-C',
        meson.current_source_dir(),
        go_build_flags,
        '-o',
        karaberus_test_output,
        '-c',
        karaberus_server_path,
    ],
    depends: [karaberus_deps, karaberus],
)

if go_os == '' and go_arch == ''
    openapi_spec = custom_target(
        'openapi_spec',
        command: [karaberus, 'openapi'],
        output: 'openapi.yaml',
        capture: true,
        depend_files: karaberus_inputs,
    )
else
    if qemu_user.found() and get_option('sysroot') != ''
        openapi_spec = custom_target(
            'openapi_spec',
            command: [
                qemu_user,
                '-L',
                get_option('sysroot'),
                karaberus,
                'openapi',
            ],
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

run_target(
    'run',
    command: [karaberus],
    depends: [karaberus_deps, karaberus_ui],
    env: go_env,
)

test_deps = [karaberus_deps, karaberus_test]

if ffmpeg.found()
    mkv_test_file = 'karaberus_test.mkv'
    test_deps += custom_target(
        'create_test_mkv',
        output: mkv_test_file,
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
            mkv_test_file,
        ],
    )

    opus_test_file = 'karaberus_test.opus'
    test_deps += custom_target(
        'create_test_opus',
        output: opus_test_file,
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
            opus_test_file,
        ],
    )
else
    warning('ffmpeg executable was not found, S3 tests are disabled')
    go_test_env.set('SKIP_S3_TESTS', 'y')
    summary('S3 tests', 'Disabled')
endif

errcheck = custom_target(
    'errcheck',
    command: [go, 'install', 'github.com/kisielk/errcheck@latest'],
    output: 'errcheck',
    env: {'GOBIN': meson.current_build_dir()},
)

if get_option('test')
    test(
        'karaberus_run_test',
        karaberus_test,
        depends: test_deps,
        env: go_test_env,
        timeout: 300,
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

test(
    'meson_format_check',
    meson_cmd,
    args: ['format', '-r', '-q', meson.current_source_dir()],
)
