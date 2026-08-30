// Ambient types for the Telegram Mini App runtime injected by
// https://telegram.org/js/telegram-web-app.js (loaded in index.html).
// Only the surface the mini-app auto-login actually touches is declared.

interface TelegramWebApp {
    initData: string
    ready(): void
}

interface Telegram {
    WebApp: TelegramWebApp
}

interface Window {
    Telegram?: Telegram
}
