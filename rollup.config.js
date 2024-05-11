import commonjs from '@rollup/plugin-commonjs';
import resolve from '@rollup/plugin-node-resolve';
import typescript from '@rollup/plugin-typescript';
import postcss from 'rollup-plugin-postcss';
import { terser } from 'rollup-plugin-terser';

const production = !process.env.ROLLUP_WATCH;

export default [{
	input: './static/main.ts',
	output: {
		file: './public/bundle-main.js',
		format: 'iife',
		name: 'mainBundle',
		exports: 'named'
	},
	plugins: [
		resolve(),
		typescript(),
		production && terser() // minify, but only in production
	]
}, {
		input: './static/app.ts',
		output: {
			file: './public/bundle-app.js',
			format: 'iife',
			name: 'appBundle',
			exports: 'named',
		},
		plugins: [
			resolve(),
			typescript(),
			production && terser()
		]
}, {
		input: './static/packages.js',
		output: {
			name: 'packagesBundle',
			exports: 'named',
			file: './public/bundle-packages.js',
			format: 'iife',
		},
		plugins: [
			resolve(),
			commonjs(),
			postcss(),
			production && terser() // minify, but only in production
		]
	}];
