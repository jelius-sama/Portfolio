import { StrictMode, useState } from 'react'
import '@/index.css'
import { BrowserRouter } from 'react-router-dom'
import { Toaster } from '@/components/ui/sonner'
import { ThemeProvider } from '@/contexts/theme'
import { createRoot } from 'react-dom/client'
import { ConfigProvider } from '@/contexts/config'
import { lazy, Suspense, useLayoutEffect, useEffect, type ReactNode } from 'react'
import { Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { Header } from "@/components/layout/header"
import { useConfig } from "@/contexts/config"
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'

const queryClient = new QueryClient()

const Home = lazy(() => import("@/pages/home"))
const InternalServerError = lazy(() => import("@/pages/internal-server-error"))
const Blogs = lazy(() => import("@/pages/blogs"))
const Blog = lazy(() => import("@/pages/blog"))
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

const generateSessionID = () => {
  const id = crypto.randomUUID()
  localStorage.setItem("session_id", id)
  return id
}

const sendAnalytics = async () => {
  const url = new URL(window.location.href)
  const params = url.searchParams

  const payload = {
    session_id: localStorage.getItem("session_id") || generateSessionID(),
    event_type: "page_view",
    event_timestamp: new Date().toISOString(),
    page_url: window.location.href,
    referrer_url: document.referrer || undefined,
    ip_address: undefined,
    country: undefined,
    region: undefined,
    city: undefined,
    user_agent: navigator.userAgent,
    device_type: /Mobi|Android/.test(navigator.userAgent) ? "mobile" : "desktop",
    browser_name: navigator.userAgent.split(" ")[0],
    browser_version: undefined,
    os_name: undefined,
    os_version: undefined,
    screen_width: window.screen.width,
    screen_height: window.screen.height,
    viewport_width: window.innerWidth,
    viewport_height: window.innerHeight,
    language: navigator.language,
    utm_source: params.get("utm_source") || undefined,
    utm_medium: params.get("utm_medium") || undefined,
    utm_campaign: params.get("utm_campaign") || undefined,
    utm_term: params.get("utm_term") || undefined,
    utm_content: params.get("utm_content") || undefined,
    page_load_time_ms: performance.timing.loadEventEnd - performance.timing.navigationStart,
    time_on_page_sec: undefined,
    scroll_depth_pct: undefined,
    element_id: undefined,
    error_message: undefined,
    metadata: {}
  }

  try {
    await fetch("/api/analytics", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
  } catch (err) {
    console.error("Failed to send analytics:", err)
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

  useEffect(() => {
    sendAnalytics()
  }, [pathname])

  return (
    <Suspense>
      <Header />
      <Outlet />
    </Suspense>
  )
}

const ErrorWrapper = ({ comp }: { comp: ReactNode }) => {
  const [errorPath, setErrorPath] = useState<string | null>(null)
  const { pathname } = useLocation()
  const { setSSRData } = useConfig()
  // INFO: The following state is to avoid race condition
  const [isSSRLoaded, setIsSSRLoaded] = useState(false)

  useLayoutEffect(() => {
    const script = document.getElementById('__SERVER_DATA__')
    if (script && script.textContent) {
      try {
        const data = JSON.parse(script.textContent)
        setSSRData(data)

        // Identify if it's error data (you can refine this signature check)
        if ("status" in data && data.status === 500) {
          setErrorPath(pathname)
        }

        script.remove()
      } catch (err) {
        console.error('Failed to parse SSR data:', err)
      }
    }
    setIsSSRLoaded(true)
  }, [])

  // Clear error when navigating to a different path
  useEffect(() => {
    if (errorPath && pathname !== errorPath) {
      setErrorPath(null)
    }
  }, [pathname, errorPath])

  return !isSSRLoaded ? null : errorPath === pathname ? <InternalServerError /> : comp
}

const reactRoot = createRoot(rootEl);
reactRoot.render(
  <StrictMode>
    <ConfigProvider>
      <BrowserRouter>
        <QueryClientProvider client={queryClient}>
          <ThemeProvider defaultTheme="dark" storageKey="theme">
            <Routes>
              <Route path='/' element={<App />}>
                <Route path='/' element={<ErrorWrapper comp={<Home />} />} />
                <Route path='/links' element={<ErrorWrapper comp={<Links />} />} />
                <Route path='/blogs' element={<ErrorWrapper comp={<Blogs />} />} />
                <Route path="/blog/:id" element={<ErrorWrapper comp={<Blog />} />} />
                <Route path='*' element={<ErrorWrapper comp={<NotFound />} />} />
              </Route>
            </Routes>
            <Toaster richColors={true} />
          </ThemeProvider>
        </QueryClientProvider>
      </BrowserRouter>
    </ConfigProvider>
  </StrictMode>
);

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
