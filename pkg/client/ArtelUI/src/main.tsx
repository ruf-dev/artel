import '@/index.css'
import '@/processes/AuthFetchInterceptor.ts'

import {createRoot} from 'react-dom/client'
import {BrowserRouter} from "react-router-dom"
import {QueryClient} from "@tanstack/react-query"
import {PersistQueryClientProvider} from "@tanstack/react-query-persist-client"
import {createSyncStoragePersister} from "@tanstack/query-sync-storage-persister"

import Router from "@/app/routing/Router.tsx"
import {
    queryCacheBuster,
    queryCacheMaxAge,
    queryCacheStorageKey,
    shouldPersistQuery,
} from "@/app/queryPersist.ts"

const queryClient = new QueryClient()

const persister = createSyncStoragePersister({
    storage: window.localStorage,
    key: queryCacheStorageKey(),
})

const persistOptions = {
    persister,
    maxAge: queryCacheMaxAge,
    buster: queryCacheBuster,
    dehydrateOptions: {shouldDehydrateQuery: shouldPersistQuery},
}

createRoot(document.getElementById('root')!).render(
    <PersistQueryClientProvider client={queryClient} persistOptions={persistOptions}>
        <BrowserRouter>
            <link href="https://fonts.googleapis.com/css2?family=Comfortaa:wght@400;500;600;700&display=swap" rel="stylesheet"/>
            <Router/>
        </BrowserRouter>
    </PersistQueryClientProvider>
)
