const CACHE_NAME = "jelius-pwa-cache-v1";

const urlsToCache = [
  "/",
  "/assets/compressed/jelius.webp",
  "/assets/manifest.json",
  "/assets/icons/apple-icon-180.png",
  "/assets/icons/apple-splash-2048-2732.jpg",
  "/assets/icons/apple-splash-2732-2048.jpg",
  "/assets/icons/apple-splash-1668-2388.jpg",
  "/assets/icons/apple-splash-2388-1668.jpg",
  "/assets/icons/apple-splash-1536-2048.jpg",
  "/assets/icons/apple-splash-2048-1536.jpg",
  "/assets/icons/apple-splash-1488-2266.jpg",
  "/assets/icons/apple-splash-2266-1488.jpg",
  "/assets/icons/apple-splash-1640-2360.jpg",
  "/assets/icons/apple-splash-2360-1640.jpg",
  "/assets/icons/apple-splash-1668-2224.jpg",
  "/assets/icons/apple-splash-2224-1668.jpg",
  "/assets/icons/apple-splash-1620-2160.jpg",
  "/assets/icons/apple-splash-2160-1620.jpg",
  "/assets/icons/apple-splash-1290-2796.jpg",
  "/assets/icons/apple-splash-2796-1290.jpg",
  "/assets/icons/apple-splash-1179-2556.jpg",
  "/assets/icons/apple-splash-2556-1179.jpg",
  "/assets/icons/apple-splash-1284-2778.jpg",
  "/assets/icons/apple-splash-2778-1284.jpg",
  "/assets/icons/apple-splash-1170-2532.jpg",
  "/assets/icons/apple-splash-2532-1170.jpg",
  "/assets/icons/apple-splash-1125-2436.jpg",
  "/assets/icons/apple-splash-2436-1125.jpg",
  "/assets/icons/apple-splash-1242-2688.jpg",
  "/assets/icons/apple-splash-2688-1242.jpg",
  "/assets/icons/apple-splash-828-1792.jpg",
  "/assets/icons/apple-splash-1792-828.jpg",
  "/assets/icons/apple-splash-1242-2208.jpg",
  "/assets/icons/apple-splash-2208-1242.jpg",
  "/assets/icons/apple-splash-750-1334.jpg",
  "/assets/icons/apple-splash-1334-750.jpg",
  "/assets/icons/apple-splash-640-1136.jpg",
  "/assets/icons/apple-splash-1136-640.jpg",
];

// Install event — pre-cache important files
self.addEventListener("install", (event) => {
  self.skipWaiting(); // Force activation after install
  event.waitUntil(
    caches
      .open(CACHE_NAME)
      .then((cache) => cache.addAll(urlsToCache))
      .catch((err) => console.error("Cache install failed", err)),
  );
});

// Activate event — clean up old caches
self.addEventListener("activate", (event) => {
  clients.claim(); // Take control of pages ASAP
  event.waitUntil(
    caches.keys().then((cacheNames) =>
      Promise.all(
        cacheNames.map((name) => {
          if (name !== CACHE_NAME) {
            return caches.delete(name);
          }
        }),
      ),
    ),
  );
});

// Fetch event — try cache first, fall back to network
self.addEventListener("fetch", (event) => {
  if (event.request.method !== "GET") return;

  event.respondWith(
    caches.match(event.request).then((cached) => {
      return (
        cached ||
        fetch(event.request)
          .then((response) => {
            const cloned = response.clone();
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(event.request, cloned);
            });
            return response;
          })
          .catch(() => {
            // Optional: return a fallback offline page here
          })
      );
    }),
  );
});
