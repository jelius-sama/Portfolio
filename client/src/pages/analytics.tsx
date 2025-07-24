import { TerminalWindow } from "@/components/ui/terminal-window"
import { Select } from "@/components/ui/select"
import { BarChart3, Globe, Monitor, Clock, Users, Eye, TrendingUp, Activity } from "lucide-react"
import type { AnalyticsRecord, AnalyticsSummary } from "@/types/analytics"
import { useEffect, useMemo, useState, Fragment } from "react"
import { useConfig } from "@/contexts/config"
import { StaticMetadata } from "@/contexts/metadata"
import { useQuery } from "@tanstack/react-query"
import { toast } from "sonner"

// Utility functions for safe data handling
const safeString = (value: string | undefined | null, fallback = "Unknown"): string => {
  return value && value.trim() !== "" ? value : fallback
}

const safeNumber = (value: number | undefined | null, fallback = 0): number => {
  return typeof value === "number" && !isNaN(value) ? value : fallback
}

function safeArray<T>(value: T[] | undefined | null): T[] {
  return Array.isArray(value) ? value : [];
}

const extractPagePath = (url: string | undefined): string => {
  if (!url) return "/"
  try {
    const urlObj = new URL(url)
    return urlObj.pathname || "/"
  } catch {
    return url.startsWith("/") ? url : "/"
  }
}

const extractDomain = (url: string | undefined): string => {
  if (!url) return "Direct"
  try {
    const urlObj = new URL(url)
    return urlObj.hostname || "Direct"
  } catch {
    return "Direct"
  }
}

const getBrowserName = (userAgent: string | undefined): string => {
  if (!userAgent) return "Unknown"

  if (userAgent.includes("Firefox")) return "Firefox"
  if (userAgent.includes("Safari")) return "Safari"
  if (userAgent.includes("Chrome")) return "Chrome"
  if (userAgent.includes("Edge")) return "Edge"
  if (userAgent.includes("Opera")) return "Opera"

  return "Other"
}

const getLocationString = (record: AnalyticsRecord): string => {
  const parts = [record.city, record.region, record.country].filter(Boolean)
  return parts.length > 0 ? parts.join(", ") : "Unknown Location"
}

// TODO: Implement authentication
export default function Analytics() {
  const [timeRange, setTimeRange] = useState<"24h" | "7d" | "30d" | "90d">("7d")
  const [eventFilter, setEventFilter] = useState<"all" | "page_view" | "click" | "error">("all")
  const [data, setData] = useState<AnalyticsRecord[]>([])
  const { app: { portfolio: me } } = useConfig()

  const timeRangeOptions = [
    { value: "24h", label: "Last 24 Hours" },
    { value: "7d", label: "Last 7 Days" },
    { value: "30d", label: "Last 30 Days" },
    { value: "90d", label: "Last 90 Days" },
  ]

  const eventFilterOptions = [
    { value: "all", label: "All Events" },
    { value: "page_view", label: "Page Views" },
    { value: "click", label: "Clicks" },
    { value: "error", label: "Errors" },
  ]

  const { isPending, error: neofetchError, data: neofetchData } = useQuery({
    queryKey: ['neofetch'],
    queryFn: async () => {
      const res = await fetch('/api/neofetch')
      if (!res.ok) {
        const errMsg = "Failed to fetch netfetch command!"
        toast.error(errMsg)
        throw new Error(errMsg)
      }
      return res.json() as Promise<Array<{ label: string, value: string }>>
    }
  })

  const neofetch: undefined | Array<{ label: string, value: string }> = isPending ? undefined : neofetchError ? me.neofetch : neofetchData && neofetchData

  useEffect(() => {
    const controller = new AbortController()
    const { signal } = controller

    const fetchData = async () => {
      try {
        const params = new URLSearchParams()
        params.set("time_range", timeRange)
        params.set("event_type", eventFilter)

        const res = await fetch(`/api/analytics?${params.toString()}`, { signal })
        const json = await res.json()
        setData(json)
      } catch (err) {
        if (err instanceof Error && err.name !== "AbortError") {
          console.error("Failed to fetch analytics:", err)
        }
      }
    }

    fetchData()
    return () => controller.abort()
  }, [timeRange, eventFilter])

  const filteredData = useMemo(() => {
    return safeArray(data)
  }, [eventFilter, data])

  const summary: AnalyticsSummary = useMemo(() => {
    const pageViews = filteredData.filter((r) => r.event_type === "page_view")
    const uniqueSessions = new Set(filteredData.map((r) => r.session_id)).size

    // Safe calculations with fallbacks
    const validLoadTimes = pageViews
      .map(r => safeNumber(r.page_load_time_ms))
      .filter(time => time > 0)

    const validTimesOnPage = pageViews
      .map(r => safeNumber(r.time_on_page_sec))
      .filter(time => time > 0)

    // Calculate top pages safely
    const pageCount = new Map<string, number>()
    filteredData.forEach(record => {
      const page = extractPagePath(record.page_url)
      pageCount.set(page, (pageCount.get(page) || 0) + 1)
    })

    // Calculate referrer domains safely
    const referrerCount = new Map<string, number>()
    filteredData.forEach(record => {
      const domain = extractDomain(record.referrer_url)
      referrerCount.set(domain, (referrerCount.get(domain) || 0) + 1)
    })

    // Calculate device types safely
    const deviceCount = new Map<string, number>()
    filteredData.forEach(record => {
      const device = safeString(record.device_type, "unknown")
      deviceCount.set(device, (deviceCount.get(device) || 0) + 1)
    })

    // Calculate browsers safely
    const browserCount = new Map<string, number>()
    filteredData.forEach(record => {
      const browser = getBrowserName(record.user_agent)
      browserCount.set(browser, (browserCount.get(browser) || 0) + 1)
    })

    return {
      totalPageViews: pageViews.length,
      uniqueVisitors: uniqueSessions,
      avgTimeOnPage: validTimesOnPage.length > 0
        ? validTimesOnPage.reduce((acc, time) => acc + time, 0) / validTimesOnPage.length
        : 0,
      avgPageLoadTime: validLoadTimes.length > 0
        ? validLoadTimes.reduce((acc, time) => acc + time, 0) / validLoadTimes.length
        : 0,
      topPages: Array.from(pageCount.entries())
        .map(([page, views]) => ({ page, views }))
        .sort((a, b) => b.views - a.views)
        .slice(0, 5),
      topCountries: Array.from(referrerCount.entries())
        .map(([country, visitors]) => ({ country, visitors }))
        .sort((a, b) => b.visitors - a.visitors)
        .slice(0, 5),
      deviceTypes: Array.from(deviceCount.entries())
        .map(([type, count]) => ({ type, count }))
        .sort((a, b) => b.count - a.count),
      browsers: Array.from(browserCount.entries())
        .map(([browser, count]) => ({ browser, count }))
        .sort((a, b) => b.count - a.count),
      recentEvents: filteredData.slice(0, 10),
    }
  }, [filteredData])

  const formatTimestamp = (timestamp: string | undefined): string => {
    if (!timestamp) return "Unknown time"
    try {
      return new Date(timestamp).toLocaleString()
    } catch {
      return "Invalid date"
    }
  }

  const getEventTypeColor = (eventType: string | undefined): string => {
    switch (eventType) {
      case "page_view":
        return "text-green-400"
      case "click":
        return "text-blue-400"
      case "error":
        return "text-red-400"
      default:
        return "text-gray-400"
    }
  }

  const getEventTypeIcon = (eventType: string | undefined): string => {
    switch (eventType) {
      case "page_view":
        return "👁"
      case "click":
        return "👆"
      case "error":
        return "❌"
      default:
        return "📊"
    }
  }

  const formatDuration = (seconds: number): string => {
    if (seconds === 0) return "N/A"
    if (seconds < 60) return `${Math.round(seconds)}s`
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = Math.round(seconds % 60)
    return `${minutes}m ${remainingSeconds}s`
  }

  const formatLoadTime = (ms: number): string => {
    if (ms === 0) return "N/A"
    if (ms < 1000) return `${Math.round(ms)}ms`
    return `${(ms / 1000).toFixed(1)}s`
  }

  return (
    <Fragment>
      <StaticMetadata />

      <section className="max-w-6xl mx-auto pt-20 pb-12 px-6">
        {/* Page Title */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-orange-400 to-orange-600 flex items-center justify-center">
            <BarChart3 className="text-black" size={24} />
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Analytics Dashboard</h1>
          <p className="text-gray-400 font-mono">{me.first_name}'s portfolio performance and visitor insights</p>
        </div>

        {/* Filters */}
        <div className="flex flex-wrap gap-4 justify-center mb-8">
          <Select
            value={timeRange}
            onChange={(value) => setTimeRange(value as typeof timeRange)}
            options={timeRangeOptions}
            className="w-48"
          />
          <Select
            value={eventFilter}
            onChange={(value) => setEventFilter(value as typeof eventFilter)}
            options={eventFilterOptions}
            className="w-48"
          />
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <TerminalWindow title="page-views" className="h-full">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-green-500/20 rounded-full flex items-center justify-center">
                <Eye className="text-green-400" size={24} />
              </div>
              <div>
                <div className="text-2xl font-bold text-white">{summary.totalPageViews}</div>
                <div className="text-gray-400 text-sm font-mono">Total Page Views</div>
              </div>
            </div>
          </TerminalWindow>

          <TerminalWindow title="unique-visitors" className="h-full">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-blue-500/20 rounded-full flex items-center justify-center">
                <Users className="text-blue-400" size={24} />
              </div>
              <div>
                <div className="text-2xl font-bold text-white">{summary.uniqueVisitors}</div>
                <div className="text-gray-400 text-sm font-mono">Unique Visitors</div>
              </div>
            </div>
          </TerminalWindow>

          <TerminalWindow title="avg-time-on-page" className="h-full">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-orange-500/20 rounded-full flex items-center justify-center">
                <Clock className="text-orange-400" size={24} />
              </div>
              <div>
                <div className="text-2xl font-bold text-white">{formatDuration(summary.avgTimeOnPage)}</div>
                <div className="text-gray-400 text-sm font-mono">Avg. Time on Page</div>
              </div>
            </div>
          </TerminalWindow>

          <TerminalWindow title="avg-load-time" className="h-full">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-purple-500/20 rounded-full flex items-center justify-center">
                <TrendingUp className="text-purple-400" size={24} />
              </div>
              <div>
                <div className="text-2xl font-bold text-white">{formatLoadTime(summary.avgPageLoadTime)}</div>
                <div className="text-gray-400 text-sm font-mono">Avg. Load Time</div>
              </div>
            </div>
          </TerminalWindow>
        </div>

        {/* Charts and Data */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          {/* Top Pages */}
          <TerminalWindow title="top-pages">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ cat /var/log/popular-pages.log</div>
              <div className="space-y-3">
                {summary.topPages.length > 0 ? (
                  summary.topPages.map((page, index) => (
                    <div key={`${page.page}-${index}`} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <span className="text-orange-400">{index + 1}.</span>
                        <span className="text-gray-300 truncate">{page.page}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="w-20 bg-gray-800 rounded-full h-2">
                          <div
                            className="bg-orange-400 h-2 rounded-full"
                            style={{
                              width: `${Math.max(20, (page.views / Math.max(...summary.topPages.map((p) => p.views))) * 100)}%`,
                            }}
                          ></div>
                        </div>
                        <span className="text-gray-400 text-xs w-8">{page.views}</span>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-gray-500 text-center py-4">No page data available</div>
                )}
              </div>
            </div>
          </TerminalWindow>

          {/* Top Referrers */}
          <TerminalWindow title="traffic-sources">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ referrer-stats --top-sources</div>
              <div className="space-y-3">
                {summary.topCountries.length > 0 ? (
                  summary.topCountries.map((source, index) => (
                    <div key={`${source.country}-${index}`} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <Globe className="text-orange-400" size={16} />
                        <span className="text-gray-300 truncate">{source.country}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="w-16 bg-gray-800 rounded-full h-2">
                          <div
                            className="bg-green-400 h-2 rounded-full"
                            style={{
                              width: `${Math.max(20, (source.visitors / Math.max(...summary.topCountries.map((c) => c.visitors))) * 100)}%`,
                            }}
                          ></div>
                        </div>
                        <span className="text-gray-400 text-xs w-6">{source.visitors}</span>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-gray-500 text-center py-4">No referrer data available</div>
                )}
              </div>
            </div>
          </TerminalWindow>

          {/* Device Types */}
          <TerminalWindow title="device-breakdown">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ device-analyzer --summary</div>
              <div className="space-y-3">
                {summary.deviceTypes.length > 0 ? (
                  summary.deviceTypes.map((device, index) => (
                    <div key={`${device.type}-${index}`} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <Monitor className="text-orange-400" size={16} />
                        <span className="text-gray-300 capitalize">{device.type}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="w-20 bg-gray-800 rounded-full h-2">
                          <div
                            className="bg-blue-400 h-2 rounded-full"
                            style={{
                              width: `${Math.max(20, (device.count / Math.max(...summary.deviceTypes.map((d) => d.count))) * 100)}%`,
                            }}
                          ></div>
                        </div>
                        <span className="text-gray-400 text-xs w-8">{device.count}</span>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-gray-500 text-center py-4">No device data available</div>
                )}
              </div>
            </div>
          </TerminalWindow>

          {/* Browser Stats */}
          <TerminalWindow title="browser-stats">
            <div className="font-mono text-sm">
              <div className="text-orange-400 mb-4">$ browser-usage --report</div>
              <div className="space-y-3">
                {summary.browsers.length > 0 ? (
                  summary.browsers.map((browser, index) => (
                    <div key={`${browser.browser}-${index}`} className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <Activity className="text-orange-400" size={16} />
                        <span className="text-gray-300">{browser.browser}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="w-20 bg-gray-800 rounded-full h-2">
                          <div
                            className="bg-purple-400 h-2 rounded-full"
                            style={{
                              width: `${Math.max(20, (browser.count / Math.max(...summary.browsers.map((b) => b.count))) * 100)}%`,
                            }}
                          ></div>
                        </div>
                        <span className="text-gray-400 text-xs w-8">{browser.count}</span>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-gray-500 text-center py-4">No browser data available</div>
                )}
              </div>
            </div>
          </TerminalWindow>
        </div>

        {/* Recent Events */}
        <TerminalWindow title="recent-events" className="mb-8">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ tail -f /var/log/analytics.log</div>
            <div className="space-y-2 max-h-96 overflow-y-auto">
              {summary.recentEvents.length > 0 ? (
                summary.recentEvents.map((event) => (
                  <div
                    key={event.id}
                    className="flex items-start gap-3 p-3 bg-gray-800/30 rounded border border-gray-700/50"
                  >
                    <span className="text-lg">{getEventTypeIcon(event.event_type)}</span>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className={`font-bold ${getEventTypeColor(event.event_type)}`}>
                          {safeString(event.event_type, "UNKNOWN").toUpperCase()}
                        </span>
                        <span className="text-gray-500 text-xs">{formatTimestamp(event.event_timestamp)}</span>
                      </div>
                      <div className="text-gray-300 text-sm mb-1">
                        <span className="text-orange-400">URL:</span> {extractPagePath(event.page_url)}
                      </div>
                      <div className="flex flex-wrap gap-4 text-xs text-gray-400">
                        <span>
                          <span className="text-orange-400">Location:</span> {getLocationString(event)}
                        </span>
                        {event.device_type && (
                          <span>
                            <span className="text-orange-400">Device:</span> {event.device_type}
                          </span>
                        )}
                        <span>
                          <span className="text-orange-400">Browser:</span> {getBrowserName(event.user_agent)}
                        </span>
                        {event.page_load_time_ms && event.page_load_time_ms > 0 && (
                          <span>
                            <span className="text-orange-400">Load Time:</span> {formatLoadTime(event.page_load_time_ms)}
                          </span>
                        )}
                      </div>
                      {event.error_message && (
                        <div className="text-red-400 text-sm mt-1">
                          <span className="text-orange-400">Error:</span> {event.error_message}
                        </div>
                      )}
                      {event.element_id && (
                        <div className="text-blue-400 text-sm mt-1">
                          <span className="text-orange-400">Element:</span> {event.element_id}
                        </div>
                      )}
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-gray-500 text-center py-8">No recent events available</div>
              )}
            </div>
          </div>
        </TerminalWindow>

        {/* System Status */}
        <TerminalWindow title="system-status">
          <div className="font-mono text-sm">
            <div className="text-orange-400 mb-4">$ systemctl status analytics.service</div>
            <div className="space-y-1 text-gray-300">
              <div>● analytics.service - Portfolio Analytics Service</div>
              <div className="ml-2">Loaded: loaded (/etc/systemd/system/analytics.service; enabled)</div>
              <div className="ml-2">
                Active: <span className="text-green-400">active (running)</span>
              </div>
              <div className="ml-2">Events processed: {filteredData.length}</div>
              <div className="ml-2">
                Last event: {formatTimestamp(filteredData[filteredData.length - 1]?.event_timestamp)}
              </div>
              <div className="ml-2">Data collection: Client-side tracking enabled</div>
              <div className="ml-2">Missing data: Geographic info, OS details, browser versions</div>
            </div>
            <div className="text-orange-400 mt-4 mb-2 sm:block hidden">$ neofetch</div>
            <div className="text-orange-400 mt-4 mb-2 sm:hidden block">$ neofetch | less</div>

            {/* Neofetch Output */}
            <div className="flex gap-6 mb-6">
              {/* Left side - QR Code */}
              <div className="flex-shrink-0">
                <div className="w-64 h-64 bg-gray-800 border border-orange-500/30 rounded flex items-center justify-center">
                  {me.links["jelius.dev/links"]?.qr_code_link ? (
                    <img
                      src={me.links["jelius.dev/links"]?.qr_code_link}
                      alt="QR Code to links page"
                      className="w-62 h-62 rounded object-cover"
                    />
                  ) : (
                    <span className="text-gray-500">QR Code not available</span>
                  )}

                </div>
              </div>

              {/* Right side - System Info */}
              {neofetch && (
                <div className="sm:block flex-1 text-gray-300 space-y-1 hidden">
                  {neofetch.map((stat, index) => (
                    <div key={index}>
                      <span className="text-orange-400">{stat.label}:</span> {stat.value}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </TerminalWindow>

        {/* Footer */}
        <div className="text-center text-gray-500 text-sm font-mono mt-12">
          <div>Analytics Dashboard • Real-time monitoring</div>
        </div>
      </section>
    </Fragment>
  )
}
