# npm wrapper env
npm_env = environment()
npm_env.set('NPM', npm.full_path())
npm_env.set('SOURCE', meson.current_source_dir())
npm_env.set('VITE_REPOSITORY_URL', repo_url)
npm_env.set('VITE_REVISION', revision)

ui_pkg = files('package-lock.json')

ui_src = files(
    'public' / 'vite.svg',
    'index.html',
    'package.json',
    'postcss.config.js',
    'tailwind.config.js',
    'vite.config.ts',
)

subdir('src')

npx = find_program('npx')

ui_deps_target = []
openapi_ts = custom_target(
    'openapi_ts',
    command: [
        npx,
        'openapi-typescript',
        openapi_spec,
        '-o',
        '@OUTPUT@',
    ],
    output: 'karaberus.d.ts',
)
symlink = meson.global_source_root() / 'scripts' / 'symlink.py'
karaberus_d_ts = custom_target(
    'symlink_karaberus_d_ts',
    command: [
        symlink,
        openapi_ts,
        '@CURRENT_SOURCE_DIR@/src/utils/karaberus.d.ts',
    ],
    output: 'symlink_karaberus_d_ts',
    build_always_stale: true,
    capture: true,
)
ui_deps_target += karaberus_d_ts

npm_run = meson.global_source_root() / 'scripts' / 'npm_run.py'

ui_deps_target += custom_target(
    'ui_deps',
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
    depends: ui_deps_target,
    install: true,
    install_dir: 'share' / 'karaberus',
    install_tag: 'runtime',
)

run_target(
    'ui_run',
    env: npm_env,
    command: [python, npm_run, 'dev'],
    depends: ui_deps_target,
)

test(
    'karaberus_ui_lint',
    python,
    args: [npm_run, 'lint'],
    env: npm_env,
    depends: ui_deps_target,
    timeout: 300,
)

test(
    'karaberus_ui_check',
    python,
    args: [npm_run, 'check'],
    env: npm_env,
    depends: ui_deps_target,
)
