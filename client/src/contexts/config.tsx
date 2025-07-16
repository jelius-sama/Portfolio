import { type AppConfig } from "@/types/config";
import { type StaticRoute } from "@/types/static.route";
import { createContext, useContext } from "react"
import AppConfigJSON from "~/client.config.json";
import StaticRouteJSON from "~/static.route.json";

type ConfigState = {
  app: AppConfig;
  staticRoute: StaticRoute[];
}

const ConfigProviderContext = createContext<ConfigState | undefined>(undefined)

export function ConfigProvider({ children }: { children: React.ReactNode }) {
  const appConfig: ConfigState["app"] = AppConfigJSON;
  const staticRoute: ConfigState["staticRoute"] = StaticRouteJSON as StaticRoute[];

  return (
    <ConfigProviderContext.Provider value={{ app: appConfig, staticRoute: staticRoute }}>
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
