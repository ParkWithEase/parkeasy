import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';
import process from 'node:process';
import IstanbulPlugin from 'vite-plugin-istanbul';

export default defineConfig({
    build: {
        sourcemap: process.env.VITE_COVERAGE == 'true'
    },
    plugins: [
        sveltekit(),
        ...(process.env.VITE_COVERAGE == 'true'
            ? [
                  IstanbulPlugin({
                      include: 'src/*',
                      exclude: ['node_modules', 'test/'],
                      extension: ['.js', '.ts', '.svelte'],
                      requireEnv: false,
                      forceBuildInstrument: true
                  })
              ]
            : [])
    ],
    test: {
        include: ['src/**/*.{test,spec}.{js,ts}'],
        coverage: {
            provider: 'istanbul',
            reporter: ['text', 'clover']
        }
    }
});
