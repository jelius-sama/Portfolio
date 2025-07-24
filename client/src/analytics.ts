const generateSessionID = () => {
  const id = crypto.randomUUID()
  localStorage.setItem("session_id", id)
  return id
}

export function sendAnalytics() {
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

  fetch("/api/analytics", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  })
    .catch((err) => {
      console.error("Failed to send analytics:", err);
    });
}
