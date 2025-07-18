import { StrictMode } from 'react'
import '@/index.css'
import { BrowserRouter } from 'react-router-dom'
import { Toaster } from '@/components/ui/sonner'
import { ThemeProvider } from '@/contexts/theme'
import { createRoot } from 'react-dom/client'
import { ConfigProvider } from '@/contexts/config'
import { lazy, Suspense, useLayoutEffect } from 'react'
import { Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { Header } from "@/components/layout/header"
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'

const queryClient = new QueryClient()

const Home = lazy(() => import("@/pages/home"))
const Development = lazy(() => import("@/pages/development"))
const Links = lazy(() => import("@/pages/links"))
const NotFound = lazy(() => import("@/pages/not-found"))

let rootEl = document.getElementById('root') as HTMLDivElement | null;

if (!rootEl) {
  if (process.env.NODE_ENV === "development") {
    throw new Error("Root element not found!")
  } else {
    const div = document.createElement('div');
    div.id = "root"
    document.body.appendChild(div);
    rootEl = div
  }
}

const App = () => {
  const { pathname } = useLocation();

  useLayoutEffect(() => {
    document.documentElement.scrollTo({
      top: 0,
      left: 0,
      behavior: "instant",
    });
  }, [pathname]);

  return (
    <Suspense>
      <Header />
      <Outlet />
    </Suspense>
  )
}

const Entry = () => {

  return (
    <StrictMode>
      <ConfigProvider>
        <BrowserRouter>
          <QueryClientProvider client={queryClient}>
            <ThemeProvider defaultTheme="dark" storageKey="theme">
              <Routes>
                <Route path='/' element={<App />}>
                  <Route path='/' element={<Home />} />
                  <Route path='/links' element={<Links />} />
                  <Route path='/blogs' element={<Development />} />
                  <Route path='*' element={<NotFound />} />
                </Route>
              </Routes>
              <Toaster richColors={true} />
            </ThemeProvider>
          </QueryClientProvider>
        </BrowserRouter>
      </ConfigProvider>
    </StrictMode>
  )
}

const reactRoot = createRoot(rootEl);
reactRoot.render(<Entry />);

// Registering service worker
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/assets/sw.js')
      .then((registration) => {
        console.log('Service Worker registered with scope:', registration.scope);
      })
      .catch((err) => {
        console.error('Service Worker registration failed:', err);
      });
  });
}
