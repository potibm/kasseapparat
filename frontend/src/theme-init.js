try {
  if (typeof window !== "undefined" && typeof localStorage !== "undefined") {
    const stored = localStorage.getItem("flowbite-theme-mode");

    if (stored === "dark") {
      document.documentElement.classList.add("dark");
    } else if (stored === "light") {
      document.documentElement.classList.remove("dark");
    } else {
      // fallback: use system preference, if available
      if (typeof window.matchMedia === "function") {
        const prefersDark = window.matchMedia(
          "(prefers-color-scheme: dark)",
        ).matches;
        if (prefersDark) {
          document.documentElement.classList.add("dark");
        }
      }
    }
  }
} catch (e) {
  console.warn("[theme-init] Theme detection failed:", e);
}
