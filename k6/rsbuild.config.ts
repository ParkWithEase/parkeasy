import path from 'node:path';
import { defineConfig } from '@rsbuild/core';
import { pluginTypeCheck } from '@rsbuild/plugin-type-check';
import { globSync } from 'glob';

function entryFromGlob(pattern: string): Record<string, string> {
    return globSync(pattern, { dotRelative: true }).reduce(
        (acc, x) => {
            const key = path.basename(x, path.extname(x));
            acc[key] = x;
            return acc;
        },
        {} as Record<string, string>,
    );
}

export default defineConfig({
    plugins: [pluginTypeCheck()],
    source: {
        entry: entryFromGlob('./src/**/*test.ts'),
    },
    tools: {
        htmlPlugin: false,
        rspack: {
            target: ['node', 'es2022'],
        },
    },
    output: {
        target: 'node',
        filenameHash: false,
        minify: false,
        externals: [/^k6(\/.*)?/, /^https?:\/\/(\/.*)?/],
        sourceMap: {
            js: 'source-map',
        },
    },
});
