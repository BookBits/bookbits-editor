import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import postcss from 'rollup-plugin-postcss';
import { terser } from 'rollup-plugin-terser';
import typescript from '@rollup/plugin-typescript';

const production = !process.env.ROLLUP_WATCH;

export default [{
	input: 'static/main.js',
	output: {
		file: './public/bundle.js',
		format: 'iife',
		name: 'mainBundle',
		exports: 'named'
	},
	plugins: [
		resolve(), 
		commonjs(),
		postcss(),
		production && terser() // minify, but only in production
	]
}, {
		input: 'static/login.ts',
		output: {
			file: './public/login-bundle.js',
			format: 'iife',
			name: 'loginBundle',
			exports: 'named'
		},
		plugins: [
			typescript(),
			resolve(),
			production && terser() // minify, but only in production
		]
	}];
