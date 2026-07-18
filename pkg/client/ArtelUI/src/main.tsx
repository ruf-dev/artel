import '@/index.css'
import '@/processes/AuthFetchInterceptor.ts'

import {createRoot} from 'react-dom/client'
import {BrowserRouter} from "react-router-dom"
import {QueryClient, QueryClientProvider} from "@tanstack/react-query"

import Router from "@/app/routing/Router.tsx"

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
    <QueryClientProvider client={queryClient}>
        <BrowserRouter>
            <link href="https://fonts.googleapis.com/css2?family=Comfortaa:wght@400;500;600;700&display=swap" rel="stylesheet"/>
            <Router/>
        </BrowserRouter>
    </QueryClientProvider>
)
