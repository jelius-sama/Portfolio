import { lazy, Suspense } from 'react'
import { Outlet, Route, Routes } from 'react-router-dom'

const Home = lazy(() => import("@/pages/home"))
const NotFound = lazy(() => import("@/pages/not-found"))

export const Router = () => {

  return (
    <Routes>
      <Route path='/' element={<App />}>
        <Route path='/' element={<Home />} />
        <Route path='*' element={<NotFound />} />
      </Route>
    </Routes>
  )
}

function App() {

  return (
    <Suspense>
      <Outlet />
    </Suspense>
  )
}

