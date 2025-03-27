# Install Node deps on change of package.json
local_resource(
  'yarn',
  cmd='yarn install',
  deps=['package.json'],
)

# Rebuild frontend if source files change
local_resource(
  'frontend',
  cmd='node ./ci/build.mjs',
  deps=['src'],
  resource_deps=['yarn'],
)

# Rebuild and run Go webserver on code changes
local_resource(
  'server',
  cmd='go build .',
  deps=[
    'main.go',
    'pkg',
  ],
  ignore=['doc-render', 'src'],
  serve_cmd='./doc-render --listen=:15642',
  readiness_probe=probe(
    http_get=http_get_action(15642, path='/api/healthz'),
    initial_delay_secs=1,
  ),
  resource_deps=['frontend'],
)
