import { readdir, readFile, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const root = new URL('../../resources/views/', import.meta.url);
const voidTags = new Set(['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr']);
const blockTags = 'article|aside|button|details|div|fieldset|footer|form|h[1-6]|header|label|li|main|nav|ol|option|p|section|select|summary|table|tbody|td|th|thead|tr|ul';

async function files(dir) {
  const entries = await readdir(dir, { withFileTypes: true });
  const paths = await Promise.all(entries.map(async entry => {
    const path = join(dir, entry.name);
    return entry.isDirectory() ? files(path) : entry.name.endsWith('.html') ? [path] : [];
  }));
  return paths.flat().sort();
}

function markup(line) {
  const result = [];
  let i = 0;
  while (i < line.length) {
    if (line.startsWith('{{', i)) {
      const end = line.indexOf('}}', i + 2);
      if (end < 0) throw new Error('Template action spans lines');
      result.push(line.slice(i, end + 2));
      i = end + 2;
    } else if (line[i] === '<') {
      let quote = '';
      let end = i + 1;
      for (; end < line.length; end++) {
        const ch = line[end];
        if (quote) {
          if (ch === quote) quote = '';
        } else if (ch === '"' || ch === "'") {
          quote = ch;
        } else if (ch === '>') {
          end++;
          break;
        }
      }
      result.push(line.slice(i, end));
      i = end;
    } else {
      i++;
    }
  }
  return result;
}

function format(source) {
  // Break only between adjacent elements. Inline text, punctuation and Go
  // actions stay byte-for-byte intact, including textarea and pre content.
  const expanded = source.replace(new RegExp(`>(?=<(?:${blockTags})\\b)`, 'gi'), '>\n');
  const output = [];
  const stack = [];
  for (const original of expanded.split(/\r?\n/)) {
    const line = original.trim();
    if (!line) continue;
    const parts = markup(line);
    const first = parts[0] ?? '';
    const action = /^{{\s*(.*?)\s*}}$/s.exec(first)?.[1] ?? '';
    const close = /^<\/([a-z][\w:-]*)\b/i.exec(first)?.[1]?.toLowerCase();
    let indent = stack.length;
    if (action === 'end' || /^else(?:\s|$)/.test(action) || close) indent = Math.max(0, indent - 1);
    output.push('  '.repeat(indent) + line);

    for (const part of parts) {
      const command = /^{{\s*(.*?)\s*}}$/s.exec(part)?.[1] ?? '';
      const closing = /^<\/([a-z][\w:-]*)\b/i.exec(part)?.[1]?.toLowerCase();
      const opening = /^<([a-z][\w:-]*)\b/i.exec(part)?.[1]?.toLowerCase();
      if (/^(?:define|block|if|range|with)(?:\s|$)/.test(command)) {
        stack.push('go');
      } else if (command === 'end') {
        const index = stack.lastIndexOf('go');
        if (index >= 0) stack.splice(index, 1);
      } else if (closing) {
        const index = stack.lastIndexOf(`html:${closing}`);
        if (index >= 0) stack.splice(index, 1);
      } else if (opening && !voidTags.has(opening) && !part.endsWith('/>')) {
        stack.push(`html:${opening}`);
      }
    }
  }
  return output.join('\n') + '\n';
}

const check = process.argv.includes('--check');
let changed = 0;
for (const path of await files(root.pathname)) {
  const before = await readFile(path, 'utf8');
  const after = format(before);
  if (before === after) continue;
  changed++;
  if (!check) await writeFile(path, after);
  process.stdout.write(`${check ? 'Needs formatting' : 'Formatted'}: ${path}\n`);
}
if (check && changed) process.exitCode = 1;
