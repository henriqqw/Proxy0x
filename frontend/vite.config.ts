import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

// Production only: the dev server needs inline scripts and websockets for HMR.
// No remote origins: fonts are system fonts and flags are bundled.
const csp: Plugin = {
  name: 'proxy0x-csp',
  apply: 'build',
  transformIndexHtml: html => html.replace(
    '<head>',
    `<head>\n    <meta http-equiv="Content-Security-Policy" content="default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'none'; form-action 'none'" />`,
  ),
}

export default defineConfig({
  plugins: [react(), csp],
  resolve: {
    alias: {
      'wailsjs': resolve(__dirname, './wailsjs'),
    },
  },
  build: {
    outDir: 'dist',
  },
})
