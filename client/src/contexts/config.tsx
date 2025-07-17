import { createContext, useContext } from "react"
import AppConfigJSON from "~/client.config.json";
import StaticRouteJSON from "~/static.route.json";
import { type StaticRoute } from "@/types/static.route";

type ConfigState = {
  app: typeof AppConfigJSON;
  staticRoute: StaticRoute[];
}

const ConfigProviderContext = createContext<ConfigState | undefined>(undefined)

export function ConfigProvider({ children }: { children: React.ReactNode }) {

  return (
    <ConfigProviderContext.Provider value={{ app: AppConfigJSON, staticRoute: StaticRouteJSON as StaticRoute[] }}>
      {children}
    </ConfigProviderContext.Provider>
  )
}

export const useConfig = () => {
  const context = useContext(ConfigProviderContext)

  if (context === undefined)
    throw new Error("useConfig must be used within a ConfigProvider")

  return context
}
