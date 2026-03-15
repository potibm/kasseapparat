/**
 * Initializes the theme (dark/light mode) based on localStorage or system preferences.
 * This should run as early as possible to prevent "Flash of Unstyled Content" (FOUC).
 */
export const initTheme = (): void => {
  try {
    if (
      typeof globalThis !== "undefined" &&
      typeof localStorage !== "undefined"
    ) {
      const stored = localStorage.getItem("flowbite-theme-mode");

      if (stored === "dark") {
        document.documentElement.classList.add("dark");
      } else if (stored === "light") {
        document.documentElement.classList.remove("dark");
      } else if (typeof globalThis.matchMedia === "function") {
        // fallback: use system preference, if available
        const prefersDark = globalThis.matchMedia(
          "(prefers-color-scheme: dark)",
        ).matches;
        if (prefersDark) {
          document.documentElement.classList.add("dark");
        }
      }
    }
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : "Unknown error";
    console.warn("[theme-init] Theme detection failed:", message);
  }
};

// Auto-execute if this is imported as a side-effect
initTheme();
