import { copyFileSync, mkdirSync } from 'node:fs';
import { join } from 'node:path';

const assets = [
  ['node_modules/htmx.org/dist/htmx.min.js', 'js/htmx.min.js'],
  ['resources/js/app.js', 'js/app.js'],
  ['node_modules/htmx.org/LICENSE', 'licenses/htmx.txt'],
  ['node_modules/tailwindcss/LICENSE', 'licenses/tailwindcss.txt'],
  ['resources/icons.svg', 'icons.svg'],
];

for (const [source, output] of assets) {
  const destination = join('public/assets', output);
  mkdirSync(join(destination, '..'), { recursive: true });
  copyFileSync(source, destination);
}
