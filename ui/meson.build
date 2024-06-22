ui_pkg = files('package-lock.json')

ui_src = files(
  'index.html',
  'package.json',
  'vite.config.ts',
  'public' / 'vite.svg',
)

subdir('src')

openapi_spec = custom_target(
  'openapi_spec',
  command: [go, 'run', meson.global_source_root(), 'openapi'],
  output: 'openapi.yaml',
  capture: true,
  depend_files: karaberus_inputs,
)

npx = find_program('npx')

custom_target(
  'openapi_ts',
  command: [
    npx,
    'openapi-typescript',
    openapi_spec,
    '-o', '@CURRENT_SOURCE_DIR@/src/utils/karaberus.d.ts',
  ],
  output: 'karaberus.d.ts',
)
