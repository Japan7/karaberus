mugen_client_file = files('mugen.go')

karaberus_server_files += mugen_client_file

mugen_test_inputs = mugen_client_file + files('mugen_test.go')
mugen_test_output = 'mugen_test'

mugen_test = custom_target(
    'mugen_test',
    env: go_env,
    input: mugen_test_inputs,
    output: mugen_test_output,
    command: [
        go,
        'test',
        '-C',
        meson.current_source_dir(),
        '-buildmode=pie',
        '-trimpath',
        '-ldflags',
        '-linkmode=external -s',
        '-o',
        meson.current_build_dir() / mugen_test_output,
        '-c',
        meson.current_source_dir(),
    ],
    depends: karaberus_generate,
)

if get_option('test')
    test('mugen_run_test', mugen_test, env: go_test_env, timeout: 300)
endif
