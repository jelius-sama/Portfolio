import { StrictMode } from 'react'
import '@/index.css'
import { BrowserRouter } from 'react-router-dom'
import { Router } from '@/index'
import { Toaster } from '@/components/ui/sonner'
import { ThemeProvider } from '@/contexts/theme'
import { createRoot } from 'react-dom/client'
import { ConfigProvider } from '@/contexts/config'

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

const Entry = () => {
  return (
    <StrictMode>
      <ConfigProvider>
        <BrowserRouter>
          <ThemeProvider defaultTheme="dark" storageKey="theme">
            <Router />
            <Toaster richColors={true} />
          </ThemeProvider>
        </BrowserRouter>
      </ConfigProvider>
    </StrictMode>
  )
}

const reactRoot = createRoot(rootEl);
reactRoot.render(<Entry />);
