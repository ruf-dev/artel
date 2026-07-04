import {fileURLToPath, URL} from 'node:url'

import {defineConfig, loadEnv} from 'vite'
import react from '@vitejs/plugin-react'

export default ({mode}: { mode: string }) => {
    process.env = {...process.env, ...loadEnv(mode, process.cwd())};

    const uc = defineConfig({
        base: '/',
        plugins: [react()],

        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url)),
            },
        },
    });

    uc.server = {
        port: 5175,
        host: true,
        allowedHosts: ['.loca.lt', 'localhost', '127.0.0.1', 'artel-alexskilled.redsock.ru'],
        hmr: {
            host: 'artel-alexskilled.redsock.ru',
            protocol: 'wss',
            clientPort: 443,
        }
    }

    return uc
}
