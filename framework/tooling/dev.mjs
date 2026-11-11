import { spawn } from 'node:child_process';
import { copyFileSync, existsSync, watch } from 'node:fs';
import { join } from 'node:path';

if (existsSync('.env')) process.loadEnvFile('.env');

const children = new Set();
const watchers = [];
let stopping = false;

function stop(code = 0) {
  if (stopping) return;
  stopping = true;
  for (const watcher of watchers) watcher.close();
  for (const child of children) child.kill('SIGTERM');
  process.exitCode = code;
}

function start(command, args) {
  const child = spawn(command, args, {
    stdio: 'inherit',
    env: { ...process.env, APP_ENV: 'development' },
  });
  children.add(child);
  child.on('error', error => {
    console.error(`Unable to start ${command}: ${error.message}`);
    stop(1);
  });
  child.on('exit', code => {
    children.delete(child);
    if (!stopping) stop(code || 1);
  });
}

// CSS is compiled by Tailwind; copy other edited assets for live reload.
for (const directory of ['js']) {
  watchers.push(watch(`resources/${directory}`, (event, filename) => {
    if (!filename) return;
    const source = join('resources', directory, filename);
    if (existsSync(source)) copyFileSync(source, join('public/assets', directory, filename));
  }));
}
watchers.push(watch('resources/icons.svg', () => {
  copyFileSync('resources/icons.svg', 'public/assets/icons.svg');
}));

start('npm', ['run', 'watch']);
start(process.env.AIR_BIN || 'air', ['-c', 'framework/tooling/air.toml']);
process.on('SIGINT', () => stop());
process.on('SIGTERM', () => stop());
