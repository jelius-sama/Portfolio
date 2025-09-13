import { StrictMode, useState, Fragment } from 'react'
import '@/index.css'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/contexts/theme'
import { createRoot } from 'react-dom/client'
import { ConfigProvider } from '@/contexts/config'
import { lazy, Suspense, useLayoutEffect, useEffect, type ReactNode } from 'react'
import { Outlet, Route, Routes, useLocation } from 'react-router-dom'
import { Header } from "@/components/layout/header"
import { useConfig } from "@/contexts/config"
import { QueryClientProvider, QueryClient } from '@tanstack/react-query'
import Loading from "@/pages/loading"
import { ErrorBoundary } from "@/error-boundary"
import { LoadingBoundary } from "@/loading-boundary"

const queryClient = new QueryClient()

const Home = lazy(() => import("@/pages/home"))
const Achievements = lazy(() => import("@/pages/achievements"))
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
  const priorityPaths = ["/", "/achievements", "/links", "/blogs", "/blog/*"]

  useLayoutEffect(() => {
    document.documentElement.scrollTo({
      top: 0,
      left: 0,
      behavior: "instant",
    });
  }, [pathname]);

  const isPriorityPath = priorityPaths.some((path) => {
    if (path.endsWith("/*")) {
      // INFO: Remove the '*' but keep the '/'
      const basePath = path.slice(0, -1);
      return pathname.startsWith(basePath);
    } else {
      return pathname === path;
    }
  });

  useEffect(() => {
    if (isPriorityPath) return;

    const timeout = setTimeout(() => {
      import('@/analytics')
        .then((module) => {
          module.sendAnalytics();
        })
        .catch((err) => {
          console.error("Error loading analytics function:", err);
        });
    }, 3000);

    return () => clearTimeout(timeout);
  }, [pathname]);

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
  const [status, setStatus] = useState<"pending" | "success" | "error">("pending")

  useEffect(() => {
    fetch(`/api/verify_auth`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (res.status === 200) {
          setStatus("success")
        } else if (res.status === 498) {
          // INFO: Expired token
          setStatus("error")
        } else {
          setStatus("error")
        }
      })
      .catch(() => {
        setStatus("error")
      })
  }, [])

  if (status === "pending") return <Loading />
  if (status === "error") return <NotFound />

  return page
}

const LDJson = () => {
  return <script
    type="application/ld+json"
    data-aria-hidden="true"
    aria-hidden="true"
    dangerouslySetInnerHTML={{
      __html: JSON.stringify({
        "@context": "https://schema.org",
        "@type": "Person",
        name: "Jelius",
        url: "https://jelius.dev",
        sameAs: [
          "https://www.linkedin.com/in/jelius-basumatary-485044339/",
          "https://github.com/jelius-sama"
        ]
      }),
    }}
  />
}

const reactRoot = createRoot(rootEl);
reactRoot.render(
  <StrictMode>
    <LDJson />
    <ConfigProvider>
      <BrowserRouter>
        <QueryClientProvider client={queryClient}>
          <ThemeProvider defaultTheme="dark" storageKey="theme">
            <Routes>
              <Route path='/' element={<App />}>
                <Route path='/' element={<ServerErrorWrapper comp={<Home />} />} />
                <Route path='/achievements' element={<ServerErrorWrapper comp={<Achievements />} />} />
                <Route path='/links' element={<ServerErrorWrapper comp={<Links />} />} />
                <Route path='/analytics' element={<ServerErrorWrapper comp={<Authenticate page={<Analytics />} />} />} />
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
