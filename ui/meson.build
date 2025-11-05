# npm wrapper env
npm_env = environment()
npm_env.set('NPM', npm.full_path())
npm_env.set('SOURCE', meson.current_source_dir())

ui_pkg = files('package-lock.json')

ui_src = files(
    'public' / 'vite.svg',
    'index.html',
    'package.json',
    'postcss.config.js',
    'vite.config.ts',
)

subdir('src')

npx = find_program('npx')

karaberus_ui_deps = []
openapi_ts = custom_target(
    'openapi_ts',
    command: [
        npx,
        'openapi-typescript',
        openapi_spec,
        '-o',
        '@CURRENT_SOURCE_DIR@/src/utils/karaberus.d.ts',
    ],
    output: 'karaberus.d.ts',
    capture: true,
)
karaberus_ui_deps += openapi_ts

npm_run = meson.global_source_root() / 'scripts' / 'npm_run.py'

karaberus_ui_deps += custom_target(
    'karaberus_ui_deps',
    build_by_default: true,
    env: npm_env,
    input: ui_pkg,
    output: '.npm_ci',
    command: [python, npm_run, 'ci', '@OUTPUT@'],
)

karaberus_ui = custom_target(
    'karaberus_ui',
    build_by_default: true,
    env: npm_env,
    input: [ui_src, openapi_ts],
    output: 'ui_dist',
    command: [python, npm_run, 'build', karaberus_ui_dist_dir],
    depends: karaberus_ui_deps,
    install: true,
    install_dir: 'share' / 'karaberus',
    install_tag: 'runtime',
)

run_target(
    'ui_run',
    env: npm_env,
    command: [python, npm_run, 'dev'],
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
    depends: karaberus_ui_deps,
)
