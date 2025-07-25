import { StrictMode, useState, Fragment } from 'react'
import '@/index.css'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/contexts/theme'
import { createRoot } from 'react-dom/client'
import { ConfigProvider } from '@/contexts/config'
import { lazy, Suspense, useLayoutEffect, useEffect, type ReactNode } from 'react'
import { Outlet, Route, Routes, useLocation, useParams } from 'react-router-dom'
import { Header } from "@/components/layout/header"
import { useConfig } from "@/contexts/config"
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'
import Loading from "@/pages/loading"
import { ErrorBoundary } from "@/error-boundary"
import { LoadingBoundary } from "@/loading-boundary"

const queryClient = new QueryClient()

const Home = lazy(() => import("@/pages/home"))
const Analytics = lazy(() => import("@/pages/analytics"))
const ClientError = lazy(() => import("@/pages/error"))
const InternalServerError = lazy(() => import("@/pages/internal-server-error"))
const Blogs = lazy(() => import("@/pages/blogs"))
const Blog = lazy(() => import("@/pages/blog"))
const Links = lazy(() => import("@/pages/links"))
const NotFound = lazy(() => import("@/pages/not-found"))
const Toaster = lazy(() => import('@/components/ui/sonner'))

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
  const priorityPaths = ["/", "/links", "/blogs", "/blog/"]

  useLayoutEffect(() => {
    document.documentElement.scrollTo({
      top: 0,
      left: 0,
      behavior: "instant",
    });
  }, [pathname]);

  const isPriorityPath = priorityPaths.some(
    path => pathname.startsWith(path) || path === pathname
  );

  return (
    <Fragment>
      <ErrorBoundary fallback={<ClientError />}>
        <Header />
        <Suspense fallback={isPriorityPath ? null : <Loading />}>
          {isPriorityPath ? <LoadingBoundary /> : <Outlet />}
        </Suspense>
      </ErrorBoundary>
    </Fragment>
  )
}

const ServerErrorWrapper = ({ comp }: { comp: ReactNode }) => {
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

  return !isSSRLoaded ? <Loading /> : errorPath === pathname ? <InternalServerError /> : comp
}

export const Authenticate = ({ page }: { page: React.ReactNode }) => {
  const { token } = useParams<{ token: string }>()
  const [status, setStatus] = useState<"pending" | "success" | "expired" | "error">("pending")

  useEffect(() => {
    if (!token) return

    fetch(`/api/verify_auth?token=${token}`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (res.status === 200) {
          setStatus("success")
        } else if (res.status === 498) {
          setStatus("expired")
        } else {
          setStatus("error")
        }
      })
      .catch(() => {
        setStatus("error")
      })
  }, [token])

  if (status === "pending") return <Loading />
  if (status !== "success") return <NotFound />

  return page
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
                <Route path='/' element={<ServerErrorWrapper comp={<Home />} />} />
                <Route path='/links' element={<ServerErrorWrapper comp={<Links />} />} />
                <Route path='/:token' element={<ServerErrorWrapper comp={<Authenticate page={<Analytics />} />} />} />
                <Route path='/blogs' element={<ServerErrorWrapper comp={<Blogs />} />} />
                <Route path="/blog/:id" element={<ServerErrorWrapper comp={<Blog />} />} />
                <Route path='*' element={<ServerErrorWrapper comp={<NotFound />} />} />
              </Route>
            </Routes>

            <Suspense><Toaster richColors={true} /></Suspense>
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
