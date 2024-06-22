ui_pkg = files('package-lock.json')

ui_src = files(
    'index.html',
    'package.json',
    'vite.config.ts',
    'public' / 'vite.svg',
)

subdir('src')

npx = find_program('npx')

openapi_ts = custom_target(
    'openapi_ts',
    command: [
        npx,
        'openapi-typescript',
        openapi_spec,
        '-o', '@CURRENT_SOURCE_DIR@/src/utils/karaberus.d.ts',
    ],
    output: 'karaberus.d.ts',
    capture: true,
)

npm_run = meson.global_source_root() / 'scripts' / 'npm_run.py'

karaberus_ui_deps = custom_target(
    'karaberus_ui_deps',
    build_by_default: true,
    env: npm_env,
    input: ui_pkg,
    output: '.npm_ci',
    command: [
        python,
        npm_run,
        'ci',
        '@OUTPUT@',
    ],
)

karaberus_ui = custom_target(
    'karaberus_ui',
    build_by_default: true,
    env: npm_env,
    input: [ui_src, openapi_ts],
    output: 'ui_dist',
    command: [
        python,
        npm_run,
        'build',
        karaberus_ui_dist_dir,
    ],
    depends: karaberus_ui_deps,
    install: true,
    install_dir: 'share' / 'karaberus',
)

karaberus_ui_run = run_target(
    'ui_run',
    env: npm_env,
    command: [
        python,
        npm_run,
        'dev',
    ],
    depends: karaberus_ui_deps,
)

test(
    'karaberus_ui_lint',
    python,
    args: [npm_run, 'lint'],
    env: npm_env,
    depends: karaberus_ui_deps,
)

test(
    'karaberus_ui_check',
    python,
    args: [npm_run, 'check'],
    env: npm_env,
)
