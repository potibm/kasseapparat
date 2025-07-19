const stored = localStorage.getItem("flowbite-theme-mode");

if (stored === "dark") {
  document.documentElement.classList.add("dark");
} else if (stored === "light") {
  document.documentElement.classList.remove("dark");
} else {
  // Optional fallback: prefers-color-scheme
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  if (prefersDark) {
    document.documentElement.classList.add("dark");
  }
}
